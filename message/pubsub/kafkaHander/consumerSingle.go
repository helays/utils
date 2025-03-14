package kafkaHander

import (
	"errors"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/helays/utils/message/pubsub"
	"github.com/helays/utils/tools/backoff"
	"strings"
	"time"
)

func (this *Instance) single(param pubsub.Params) {
	partitionList, err := this.consumer.Partitions(param.Topic)
	if err != nil {
		this.error("订阅发布组件订阅失败", "kafka consumer", fmt.Errorf("%s：%s", param.Topic, err.Error()))
		return
	}
	for _, partition := range partitionList {
		go this.partition(param.Topic, partition, sarama.OffsetNewest) // 默认从最新的offset开始消费
	}
}

// 分区消费
func (this *Instance) partition(topic string, partition int32, offset int64) {
	this.log("订阅发布组件", "kafka载体", "普通消费者", topic, "开始消费", "分区", partition)
	pc, err := this.consumer.ConsumePartition(topic, partition, offset)
	if err != nil {
		this.error("订阅发布组件订阅失败", "kafka consumer", fmt.Errorf("%s：%s", topic, err.Error()))
		return
	}
	defer pc.AsyncClose()
	b := backoff.NewBackoff(backoff.Exponential, time.Nanosecond, 10*time.Second, 2.0) // 失败等待时间指数递增，基数2.0
	for {
		select {
		case msg := <-pc.Messages(): // 收取数据
			b.Reset()
			if msg == nil {
				this.error("订阅发布组件订阅失败", "kafka consumer", topic, "消息为空", "分区", partition)
				continue
			}
			this.message <- msg
		case <-this.opts.Ctx.Done(): // 监听退出信号
			this.log("订阅发布组件", "kafka载体", "普通消费者", topic, "退出消费", "分区", partition)
			b.Reset()
			return
		case _err := <-pc.Errors(): // 监听错误
			if _err == nil {
				this.error("订阅发布组件订阅失败", "kafka consumer err is nil", topic, "分区", partition)
				time.Sleep(b.Next())
				continue
			}
			err = _err.Unwrap()
			this.error("订阅发布组件订阅失败", "kafka consumer", fmt.Errorf("%s：%s", topic, err.Error()), "分区", partition)
			time.Sleep(b.Next())
			if errors.Is(err, sarama.ErrOffsetOutOfRange) || strings.Contains(err.Error(), "offset out of range") {
				pc.AsyncClose()
				pc, err = this.consumer.ConsumePartition(topic, partition, offset)
				if err != nil {
					this.error("订阅发布组件订阅失败", "kafka consumer", "重新消费分区失败", fmt.Errorf("%s：%s", topic, err.Error()), "分区", partition)
					continue
				}
			}
		}
	}
}
