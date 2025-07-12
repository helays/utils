package pubsub

import (
	"context"
	"github.com/helays/utils/v2/logger/zaploger"
	"path"
	"strings"
)

const (
	CarrierRedis    = "redis"    // redis 载体
	CarrierKafka    = "kafka"    // kafka 载体
	CarrierRabbitMQ = "rabbitmq" // rabbitmq 载体
	CarrierRocketMQ = "rocketmq" // rocketmq 载体
	CarrierEtcd     = "etcd"     // etcd 载体
)

const (
	ErrExit = "exit" // 系统退出
)

type PubSub struct {
	Handler Handler
}

// Handler 订阅发布组件
type Handler interface {
	Subscribe(param Params, cbs *Cbfunc) // 订阅消息
	Publish(param Params, msg any) error // 发布消息
}

// Options 配置
type Options struct {
	Ctx   context.Context
	Loger *zaploger.Logger
}

// Params 参数
type Params struct {
	Topic string // 订阅主题
	Key   string // 订阅的key 主要在kafka的时候可以用这个
}

// String 返回字符串
func (this Params) String(isPath bool) string {
	if isPath {
		topic := this.Topic
		if this.Key != "" {
			return path.Join(topic, this.Key)
		}
	}
	list := []string{this.Topic}
	if this.Key != "" {
		list = append(list, this.Key)
	}
	return strings.Join(list, "_")
}

// Cbfunc 回调函数
type Cbfunc struct {
	CbString func(msg string)
	CbByte   func(msg []byte)
	CbAny    func(msg any)
}

// Init 初始化订阅发布组件
func Init(sb Handler) *PubSub {
	return &PubSub{Handler: sb}
}
