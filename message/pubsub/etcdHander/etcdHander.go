package etcdHander

import (
	"context"
	"fmt"
	"github.com/helays/utils/v2/logger/ulogs"
	"github.com/helays/utils/v2/logger/zaploger"
	"github.com/helays/utils/v2/message/pubsub"
	"github.com/helays/utils/v2/tools"
	"go.etcd.io/etcd/client/v3"
)

type Instance struct {
	etcdClient *clientv3.Client
	opts       *pubsub.Options
}

func New(etcd *clientv3.Client, opts *pubsub.Options) *Instance {
	return &Instance{
		etcdClient: etcd,
		opts:       opts,
	}
}

// Subscribe 订阅消息
func (this *Instance) Subscribe(param pubsub.Params, cbs *pubsub.Cbfunc) {
	topic := param.String(true)
	rch := this.etcdClient.Watch(this.opts.Ctx, topic)
	for wresp := range rch {
		for _, ev := range wresp.Events {
			switch ev.Type {
			case clientv3.EventTypePut:
				if cbs.CbString != nil {
					cbs.CbString(string(ev.Kv.Value))
				} else if cbs.CbByte != nil {
					cbs.CbByte(ev.Kv.Value)
				} else if cbs.CbAny != nil {
					cbs.CbAny(ev.Kv.Value)
				}
			}
		}
	}
}

func (this *Instance) Publish(param pubsub.Params, msg any) error {
	topic := param.String(true)
	sendMsg := tools.Any2string(msg)
	resp, err := this.etcdClient.Put(this.opts.Ctx, topic, sendMsg)
	if err != nil {
		return fmt.Errorf("etcdHander 发布消息失败: %v", err)
	}
	this.log("etcdHander 发布消息成功", "主题", topic, "版本", resp.Header.GetRevision(), "发送内容", sendMsg)
	return nil
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
