package kafkaClose

import (
	"github.com/IBM/sarama"
	"github.com/helays/utils/v2/logger/ulogs"
)

// CloseKafkaPartition 关闭分区
func CloseKafkaPartition(partition sarama.PartitionConsumer) {
	if partition == nil {
		return
	}
	ulogs.Checkerr(partition.Close(), "CloseKafkaPartition")
}

// CloseKafkaConsumerGroup 关闭消费者组
func CloseKafkaConsumerGroup(group sarama.ConsumerGroup) {
	if group != nil {
		ulogs.Checkerr(group.Close(), "CloseKafkaConsumerGroup 执行失败")
	}
}

func CloseKafkaConsumer(consumer sarama.Consumer) {
	if consumer != nil {
		ulogs.Checkerr(consumer.Close(), "CloseKafkaConsumer 执行失败")
	}
}

func CloseKafkaSyncProducer(producer sarama.SyncProducer) {
	if producer != nil {
		ulogs.Checkerr(producer.Close(), "CloseKafkaProducer 执行失败")
	}
}

func CloseKafkaAsyncProducer(producer sarama.AsyncProducer) {
	if producer != nil {
		producer.AsyncClose()
	}
}
