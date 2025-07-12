package kafka

import (
	"github.com/IBM/sarama"
	"github.com/helays/utils/v2/scram"
	"github.com/helays/utils/v2/tools"
	"time"
)

// 设置kafka
func (this *KafkaConfig) setConfig() (kafkaCfg *sarama.Config, err error) {
	kafkaCfg = sarama.NewConfig()
	kafkaCfg.Producer.Return.Successes = true
	kafkaCfg.Producer.Return.Errors = true
	if this.Version != "" {
		kafkaCfg.Version, err = sarama.ParseKafkaVersion(this.Version)
	}
	if this.Sasl {
		kafkaCfg.Net.SASL.Enable = true
		kafkaCfg.Net.SASL.User = this.User
		kafkaCfg.Net.SASL.Password = this.Password
		kafkaCfg.Net.SASL.Handshake = true
		if this.Mechanism != "" {
			kafkaCfg.Net.SASL.Mechanism = sarama.SASLMechanism(this.Mechanism)
			switch this.Mechanism {
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
func (this *KafkaConfig) consumerClientCfg() (*sarama.Config, error) {
	kafkaCfg, err := this.setConfig()
	if err != nil {
		return nil, err
	}
	kafkaCfg.Consumer.Offsets.Initial = tools.Ternary(this.Offset > -1 || this.Offset < -2, sarama.OffsetNewest, this.Offset)
	kafkaCfg.Consumer.Return.Errors = true
	return kafkaCfg, nil
}

// NewConsumerClient 创建消费者客户端
func (this *KafkaConfig) NewConsumerClient() (sarama.Consumer, error) {
	kafkaCfg, err := this.consumerClientCfg()
	if err != nil {
		return nil, err
	}
	return sarama.NewConsumer(this.Addrs, kafkaCfg)
}

// NewConsumerGroupClient 创建消费者组客户端
func (this *KafkaConfig) NewConsumerGroupClient() (sarama.ConsumerGroup, error) {
	kafkaCfg, err := this.consumerClientCfg()
	if err != nil {
		return nil, err
	}
	return sarama.NewConsumerGroup(this.Addrs, this.GroupName, kafkaCfg)
}

// 生产者配置文件
func (this *KafkaConfig) producerClientConfig() (*sarama.Config, error) {
	kafkaCfg, err := this.setConfig()
	if err != nil {
		return nil, err
	}
	kafkaCfg.Producer.Return.Successes = true
	kafkaCfg.Producer.Return.Errors = true
	kafkaCfg.Producer.RequiredAcks = sarama.WaitForAll                               // 等待所有同步副本确认
	kafkaCfg.Producer.Retry.Max = tools.Ternary(this.MaxRetry > 0, this.MaxRetry, 3) // 最大重试3次
	to := tools.AutoTimeDuration(this.Timeout, time.Second)
	if to > 0 {
		kafkaCfg.Producer.Timeout = to
	}

	return kafkaCfg, nil
}

// NewProducerSyncProducer 创建同步生产者客户端
func (this *KafkaConfig) NewProducerSyncProducer() (sarama.SyncProducer, error) {
	kafkaCfg, err := this.producerClientConfig()
	if err != nil {
		return nil, err
	}
	return sarama.NewSyncProducer(this.Addrs, kafkaCfg)
}

// NewProducerAsyncProducer 创建异步生产者客户端
func (this *KafkaConfig) NewProducerAsyncProducer() (sarama.AsyncProducer, error) {
	kafkaCfg, err := this.producerClientConfig()
	if err != nil {
		return nil, err
	}
	return sarama.NewAsyncProducer(this.Addrs, kafkaCfg)
}
