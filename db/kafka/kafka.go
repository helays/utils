package kafka

import (
	"time"

	"github.com/IBM/sarama"
	"helay.net/go/utils/v3/scram"
	"helay.net/go/utils/v3/tools"
)

// 设置 kafka
// noinspection all
func (kc *KafkaConfig) setConfig() (kafkaCfg *sarama.Config, err error) {
	kafkaCfg = sarama.NewConfig()
	kafkaCfg.Producer.Return.Successes = true
	kafkaCfg.Producer.Return.Errors = true
	if kc.Version != "" {
		kafkaCfg.Version, err = sarama.ParseKafkaVersion(kc.Version)
	}
	if kc.Sasl {
		kafkaCfg.Net.SASL.Enable = true
		kafkaCfg.Net.SASL.User = kc.User
		kafkaCfg.Net.SASL.Password = kc.Password
		kafkaCfg.Net.SASL.Handshake = true
		if kc.Mechanism != "" {
			kafkaCfg.Net.SASL.Mechanism = sarama.SASLMechanism(kc.Mechanism)
			switch kc.Mechanism {
			case sarama.SASLTypeSCRAMSHA256:
				kafkaCfg.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient { return &scram.XDGSCRAMClient{HashGeneratorFcn: scram.SHA256} }
			case sarama.SASLTypeSCRAMSHA512:
				kafkaCfg.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient { return &scram.XDGSCRAMClient{HashGeneratorFcn: scram.SHA512} }
			}
		}
	}

	return
}

// 消费者配置
// noinspection all
func (kc *KafkaConfig) consumerClientCfg() (*sarama.Config, error) {
	kafkaCfg, err := kc.setConfig()
	if err != nil {
		return nil, err
	}
	kafkaCfg.Consumer.Offsets.Initial = tools.Ternary(kc.Offset > -1 || kc.Offset < -2, sarama.OffsetNewest, kc.Offset)
	kafkaCfg.Consumer.Return.Errors = true
	return kafkaCfg, nil
}

// NewConsumerClient 创建消费者客户端
// noinspection all
func (kc *KafkaConfig) NewConsumerClient() (sarama.Consumer, error) {
	kafkaCfg, err := kc.consumerClientCfg()
	if err != nil {
		return nil, err
	}
	return sarama.NewConsumer(kc.Addrs, kafkaCfg)
}

// NewConsumerGroupClient 创建消费者组客户端
// noinspection all
func (kc *KafkaConfig) NewConsumerGroupClient() (sarama.ConsumerGroup, error) {
	kafkaCfg, err := kc.consumerClientCfg()
	if err != nil {
		return nil, err
	}
	return sarama.NewConsumerGroup(kc.Addrs, kc.GroupName, kafkaCfg)
}

// 生产者配置文件
// noinspection all
func (kc *KafkaConfig) producerClientConfig() (*sarama.Config, error) {
	kafkaCfg, err := kc.setConfig()
	if err != nil {
		return nil, err
	}
	kafkaCfg.Producer.Return.Successes = true
	kafkaCfg.Producer.Return.Errors = true
	kafkaCfg.Producer.RequiredAcks = sarama.WaitForAll                           // 等待所有同步副本确认
	kafkaCfg.Producer.Retry.Max = tools.Ternary(kc.MaxRetry > 0, kc.MaxRetry, 3) // 最大重试3次
	to := tools.AutoTimeDuration(kc.Timeout, time.Second)
	if to > 0 {
		kafkaCfg.Producer.Timeout = to
	}

	return kafkaCfg, nil
}

// NewProducerSyncProducer 创建同步生产者客户端
// noinspection all
func (kc *KafkaConfig) NewProducerSyncProducer() (sarama.SyncProducer, error) {
	kafkaCfg, err := kc.producerClientConfig()
	if err != nil {
		return nil, err
	}
	return sarama.NewSyncProducer(kc.Addrs, kafkaCfg)
}

// NewProducerAsyncProducer 创建异步生产者客户端
// noinspection all
func (kc *KafkaConfig) NewProducerAsyncProducer() (sarama.AsyncProducer, error) {
	kafkaCfg, err := kc.producerClientConfig()
	if err != nil {
		return nil, err
	}
	return sarama.NewAsyncProducer(kc.Addrs, kafkaCfg)
}
