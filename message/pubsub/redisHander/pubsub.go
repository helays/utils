package redisHander

import (
	"context"
	"github.com/helays/utils/v2/logger/ulogs"
	"github.com/helays/utils/v2/message/pubsub"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// Instance redis实例
type Instance struct {
	rdb  redis.UniversalClient
	opts *pubsub.Options
}

// New redis实例
func New(rdb *redis.UniversalClient, opts *pubsub.Options) *Instance {
	return &Instance{
		rdb:  *rdb,
		opts: opts,
	}
}

// Subscribe 订阅消息
func (this *Instance) Subscribe(param pubsub.Params, cbs *pubsub.Cbfunc) {
	topic := param.String(false)
	subIde := this.rdb.Subscribe(this.opts.Ctx, topic)
	for {
		select {
		case <-this.opts.Ctx.Done():
			return
		default:
			msg, err := subIde.ReceiveMessage(this.opts.Ctx)
			if err != nil {
				if this.opts.Loger != nil {
					this.opts.Loger.Error(context.Background(), "redis订阅消息失败", zap.String("topic", topic), zap.String("错误信息", err.Error()))
				} else {
					ulogs.Error("redis订阅消息失败", topic, err)
				}
				continue
			}
			if cbs.CbString != nil {
				cbs.CbString(msg.Payload)
			} else if cbs.CbByte != nil {
				cbs.CbByte([]byte(msg.Payload))
			} else if cbs.CbAny != nil {
				cbs.CbAny(msg.Payload)
			}
		}
	}
}

// Publish 发布消息
func (this *Instance) Publish(param pubsub.Params, msg any) error {
	topic := param.String(false)

	return this.rdb.Publish(this.opts.Ctx, topic, msg).Err()
}
