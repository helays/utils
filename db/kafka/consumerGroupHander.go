package kafka

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/helays/utils/logger/ulogs"
	"github.com/helays/utils/tools/backoff"
	"time"
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

// ConsumerGroupHandler 消费者组消费，通用处理器
func ConsumerGroupHandler(client sarama.ConsumerGroup, topics []string, ctx context.Context, handler *ConsumerGroupHander) {
	b := backoff.NewBackoff(backoff.Exponential, time.Nanosecond, 10*time.Second, 2.0) // 失败等待时间指数递增，基数2.0
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if err := client.Consume(ctx, topics, handler); err != nil {
				ulogs.Error(err, "消费组消费失败", "topic", topics)
				time.Sleep(b.Next())
			} else {
				b.Reset()
			}

		}
	}
}
