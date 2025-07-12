package kafkaHander

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/helays/utils/v2/logger/ulogs"
	"github.com/helays/utils/v2/logger/zaploger"
	"github.com/helays/utils/v2/message/pubsub"
	"github.com/helays/utils/v2/tools"
	"sync"
)

type Instance struct {
	opts          *pubsub.Options
	producer      sarama.SyncProducer  // kafka producer 同步生产者
	isGroup       bool                 // 是否消费组模式
	consumer      sarama.Consumer      // kafka consumer 消费者
	consumerGroup sarama.ConsumerGroup // kafka consumer 消费者组
	message       chan *sarama.ConsumerMessage
	topics        sync.Map // 用于注册 topics监听的安全集合
	pubHandler    sync.Map // 用于注册 topic_key 的消息处理函数
	ctx           context.Context
	cancel        context.CancelFunc
}

// New 创建kafkaHander实例
func New(opts *pubsub.Options, k ...any) (*Instance, error) {
	if len(k) < 1 {
		return nil, fmt.Errorf("kafkaHander 参数错误：缺失消费、生产客户端")
	}
	ins := Instance{
		opts: opts,
	}
	ins.ctx, ins.cancel = context.WithCancel(context.Background())
	for _, v := range k {
		switch t := v.(type) {
		case sarama.SyncProducer:
			ins.producer = t
		case sarama.Consumer:
			ins.consumer = t
			ins.isGroup = false
		case sarama.ConsumerGroup:
			//ins.consumerGroup = t
			//ins.isGroup = true
			return nil, fmt.Errorf("kafkaHander 参数错误：发布订阅不支持kafka消费组")
		}
	}
	if ins.producer == nil {
		return nil, fmt.Errorf("kafkaHander 参数错误：缺失生产客户端")
	}
	if ins.consumer == nil && ins.consumerGroup == nil {
		return nil, fmt.Errorf("kafkaHander 参数错误：缺失消费客户端")
	}
	ins.message = make(chan *sarama.ConsumerMessage, 100) // 创建消息队列
	return &ins, nil
}

// Publish 发布消息
func (this *Instance) Publish(param pubsub.Params, msg any) error {
	// 这里生成消息
	byt, err := tools.Any2bytes(msg)
	if err != nil {
		return err
	}
	partition, offset, err := this.producer.SendMessage(&sarama.ProducerMessage{
		Topic: param.Topic,
		Key:   sarama.StringEncoder(param.Key),
		Value: sarama.ByteEncoder(byt),
	})
	if err != nil {
		return fmt.Errorf("kafkaHander 发布消息失败: %v", err)
	}
	this.debug("kafkaHander 发布消息成功", "主题", param.Topic, "key", param.Key, "分区", partition, "偏移", offset)
	return nil
}

// Subscribe 订阅消息
// 关于kafka订阅，如果topic一样，需要进行合并，通过key来判断是否是同一个消息
// 如果topic不一样，就分开消费
func (this *Instance) Subscribe(param pubsub.Params, cbs *pubsub.Cbfunc) {
	this.pubHandler.Store(fmt.Sprintf("%s_%s", param.Topic, param.Key), cbs) // 注册消息处理函数
	_, ok := this.topics.Load(param.Topic)
	if ok {
		// 如果topic已经监听，就不继续了
		return
	}
	this.topics.Store(param.Topic, true) // 初始化当前 topic

	go this.msgHander()
	if this.isGroup {
		this.group(param)
	} else {
		this.single(param)
	}
}

func (this *Instance) msgHander() {
	for {
		select {
		case <-this.opts.Ctx.Done():
			return
		case msg := <-this.message:
			topic := msg.Topic
			key := sarama.StringEncoder(msg.Key)
			this.log("订阅发布组件", "kafka载体", "topic", topic, "key", key, "offset", msg.Offset, msg.Timestamp)
			_t, ok := this.pubHandler.Load(fmt.Sprintf("%s_%s", topic, key))
			if !ok {
				continue
			}
			cbs := _t.(*pubsub.Cbfunc)
			if cbs.CbString != nil {
				cbs.CbString(string(msg.Value))
			} else if cbs.CbByte != nil {
				cbs.CbByte(msg.Value)
			} else if cbs.CbAny != nil {
				cbs.CbAny(msg.Value)
			}
		}
	}
}

func (this *Instance) log(title string, args ...any) {
	if this.opts.Loger == nil {
		ulogs.Log(append([]any{title}, args...)...)
	} else {
		this.opts.Loger.Info(context.Background(), title, zaploger.Auto2Field(args...))
	}
}

func (this *Instance) error(title string, args ...any) {
	if this.opts.Loger == nil {
		ulogs.Error(append([]any{title}, args...)...)
	} else {
		this.opts.Loger.Error(context.Background(), title, zaploger.Auto2Field(args...))
	}
}

func (this *Instance) debug(title string, args ...any) {
	if this.opts.Loger == nil {
		ulogs.Debug(append([]any{title}, args...)...)
	} else {
		this.opts.Loger.Debug(context.Background(), title, zaploger.Auto2Field(args...))
	}
}
