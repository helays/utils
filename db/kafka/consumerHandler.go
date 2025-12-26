package kafka

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/IBM/sarama"
	"github.com/helays/utils/v2/close/vclose"
	"github.com/helays/utils/v2/logger/ulogs"
	"github.com/helays/utils/v2/tools"
	"github.com/helays/utils/v2/tools/backoff"
)

type ConsumerConfig struct {
	Topic    string        // 消费主题
	Offset   int64         // -1 OffsetNewest -2 OffsetOldest
	Interval time.Duration // 分区刷新周期

	// 自动重试配置
	MaxRetry   int           // 最大重试次数
	MinBackoff time.Duration // 最小避让时间
	MaxBackoff time.Duration // 最大避让时间
}

type onMessageFunc func(message *sarama.ConsumerMessage)

type ConsumerHandler struct {
	consumer sarama.Consumer
	ctx      context.Context
	opt      *ConsumerConfig
	mu       sync.Mutex // 读写锁
	pause    int32      // 1 暂停 0 运行状态
	// 现有的分区
	// 如果不存在就更新，如果取消了就删除
	partitions map[int32]context.CancelFunc
	onMessage  onMessageFunc
}

func NewConsumerHandler(ctx context.Context, consumer sarama.Consumer, opt *ConsumerConfig) (*ConsumerHandler, error) {
	c := &ConsumerHandler{
		consumer:   consumer,
		ctx:        ctx,
		opt:        opt,
		partitions: make(map[int32]context.CancelFunc),
	}

	if c.opt.Topic == "" {
		return nil, fmt.Errorf("未设置 topic")
	}

	c.opt.Interval = tools.AutoTimeDuration(c.opt.Interval, time.Second, time.Minute) // 分区刷新间隔 默认1分钟
	c.opt.MaxRetry = tools.Ternary(c.opt.MaxRetry < 1, 10, c.opt.MaxRetry)
	c.opt.MinBackoff = tools.AutoTimeDuration(c.opt.MinBackoff, time.Millisecond, time.Millisecond) // 最小避让时间 默认1微妙
	c.opt.MaxBackoff = tools.AutoTimeDuration(c.opt.MaxBackoff, time.Millisecond, 10*time.Second)   // 最大避让时间 默认10秒

	return c, nil
}

// GetPartitions 获取分区列表
// 这个函数是同步调用，依赖上层函数的锁，所以这个函数是并发安全的。
func (c *ConsumerHandler) getPartitions() []int32 {
	return tools.MapKeys(c.partitions)

}

// Pause 暂停消费
func (c *ConsumerHandler) Pause() {
	c.mu.Lock()
	defer c.mu.Unlock()
	p := map[string][]int32{
		c.opt.Topic: c.getPartitions(),
	}
	c.consumer.Pause(p)
	atomic.StoreInt32(&c.pause, 1)
}

func (c *ConsumerHandler) IsPause() bool {
	return atomic.LoadInt32(&c.pause) == 1
}

// Resume 恢复消费
func (c *ConsumerHandler) Resume() {
	c.mu.Lock()
	defer c.mu.Unlock()
	p := map[string][]int32{
		c.opt.Topic: c.getPartitions(),
	}
	c.consumer.Resume(p)
	atomic.StoreInt32(&c.pause, 0)
}

func (c *ConsumerHandler) Run(onMessage onMessageFunc) {
	if onMessage == nil {
		return
	}
	c.onMessage = onMessage
	c.refreshPartitions() // 先执行一次
	tck := time.NewTicker(c.opt.Interval)
	defer tck.Stop()

	for {
		select {
		case <-c.ctx.Done():
			c.mu.Lock()
			c.closePartitions()
			c.mu.Unlock()
			return
		case <-tck.C:
			c.refreshPartitions()
		}
	}
}

// 刷新分区周期
func (c *ConsumerHandler) refreshPartitions() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.IsPause() {
		return
	}
	partitions, err := c.consumer.Partitions(c.opt.Topic)
	if err != nil {
		ulogs.Errorf("获取分区列表失败 [topic:%s] %v", c.opt.Topic, err)
	}

	if len(partitions) == 0 {
		c.closePartitions()
		return
	}
	// 判断现有 partition 是否取消
	for partition, _ := range c.partitions {
		if !tools.Contains(partitions, partition) {
			c.closePartition(partition)
		}
	}
	// 判断是否有新增 partition
	for _, partition := range partitions {
		if _, ok := c.partitions[partition]; !ok {
			ctx, cancel := context.WithCancel(c.ctx)
			c.partitions[partition] = cancel
			go func(ctx context.Context, partition int32) {
				// 这个 goroutine 退出，就是说明分区消费者失败已经达到最大次数了。
				// 然后就应该删除当前分区的消费者数据，然后由上层Run里面的定时器，重新开始消费当前分区。
				defer func() {
					c.mu.Lock()
					defer c.mu.Unlock()
					c.closePartition(partition)
				}()
				b := backoff.NewBackoff(backoff.Exponential, c.opt.MinBackoff, c.opt.MaxBackoff, 2.0)
				for i := 0; i < c.opt.MaxRetry; i++ {
					select {
					case <-ctx.Done():
						return
					default:
						// ctx partition 已经通过闭包函数的参数传递进来，还有共享问题么？
						// 这个函数在正常消费过程中，理论上不会退出。
						partitionErr := c.partitionConsumer(ctx, partition)
						if partitionErr == nil {
							return
						}
						if partitionErr != nil {
							// 判断是否是上下文取消引起的报错，这种就退出不处理了。
							if errors.Is(partitionErr, context.Canceled) {
								return
							}
							time.Sleep(b.Next())
						}
					}
				}

			}(ctx, partition)

		}
	}
}

// 关闭所有分区消费者
func (c *ConsumerHandler) closePartitions() {
	for partition, cancel := range c.partitions {
		cancel()
		delete(c.partitions, partition)
	}
}

// 关闭指定分区消费者
func (c *ConsumerHandler) closePartition(partition int32) {
	if cancel, ok := c.partitions[partition]; ok {
		cancel()
		delete(c.partitions, partition)
	}
}

// 分区消费者
func (c *ConsumerHandler) partitionConsumer(ctx context.Context, partition int32) error {
	topic := c.opt.Topic
	ulogs.Infof("启动分区消费者 [topic: %s, partition: %d]", topic, partition)
	pc, err := c.consumer.ConsumePartition(topic, partition, c.opt.Offset)
	if err != nil {
		ulogs.Errorf("创建分区消费者失败 [topic: %s,partition: %d] %v", topic, partition, err)
		return err
	}
	defer vclose.Close(pc)
	// 创建避让时间自动递增的对象。
	b := backoff.NewBackoff(backoff.Exponential, c.opt.MinBackoff, c.opt.MaxBackoff, 2.0)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case message, ok := <-pc.Messages():
			if !ok {
				return nil // 通道关闭，退出循环
			}
			b.Reset() // 正常消费后，错误避让时间重置
			c.onMessage(message)
		case _err, ok := <-pc.Errors():
			if !ok {
				// 这个一般是通道关闭时候触发的
				return nil
			}
			err = _err.Unwrap()
			ulogs.Errorf("Kafka 消费者错误 [topic:%s, partition:%d] %v", topic, partition, err)
			time.Sleep(b.Next())
			if errors.Is(err, sarama.ErrOffsetOutOfRange) || strings.Contains(err.Error(), "offset out of range") {
				// 这两种问题，就突出当前，由上层函数决定是否要继续重试。
				return err
			}

		}
	}

}
