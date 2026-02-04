package kafkaHander

import (
	"fmt"
	"helay.net/go/utils/v3/message/pubsub"
)

// 消费组模式，因特殊原因，这只能用相同的topic,根据不同的key来区分具体业务
func (this *Instance) group(param pubsub.Params) {
	var gh = consumerGroupHander{msg: this.message}
	// 列出所有topic
	this.log("订阅发布组件", "kafka consumerGroup", param.Topic)
	for {
		select {
		case <-this.opts.Ctx.Done():
			return
		default:
			if err := this.consumerGroup.Consume(this.opts.Ctx, []string{param.Topic}, &gh); err != nil {
				this.error("订阅发布组件订阅失败", "kafka consumerGroup", fmt.Errorf("%s：%s", param.Topic, err.Error()))
			}
		}
	}
}
