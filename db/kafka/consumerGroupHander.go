package kafka

import (
	"github.com/IBM/sarama"
)

type ConsumerGroupHander struct {
	Msg chan *sarama.ConsumerMessage
}

func (this *ConsumerGroupHander) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *ConsumerGroupHander) Cleanup(session sarama.ConsumerGroupSession) error {
	// Optional: implement this if you need to clean up any state for the session.
	return nil
}

func (h *ConsumerGroupHander) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok { // 分区重新平衡或分区关闭
				return nil
			}
			h.Msg <- message
			session.MarkMessage(message, "")
		case <-session.Context().Done():
			// 消费者组会话结束；返回上下文错误
			return session.Context().Err()
		}
	}
}
