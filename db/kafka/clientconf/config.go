package clientconf

import (
	"fmt"
	"time"

	"github.com/IBM/sarama"
	cfg_proxy "github.com/helays/utils/v2/config/cfg-proxy"
	"github.com/helays/utils/v2/scram"
	"github.com/helays/utils/v2/tools"
	"github.com/rcrowley/go-metrics"
)

// noinspection all
const defaultClientID = "sarama"

func (g *GSSAPIConfig) ToGSSAPIConfig() sarama.GSSAPIConfig {
	return sarama.GSSAPIConfig{
		AuthType:           g.AuthType,
		KeyTabPath:         g.KeyTabPath,
		CCachePath:         g.CCachePath,
		KerberosConfigPath: g.KerberosConfigPath,
		ServiceName:        g.ServiceName,
		Username:           g.Username,
		Password:           g.Password,
		Realm:              g.Realm,
		DisablePAFXFAST:    g.DisablePAFXFAST,
	}
}

// ToSarama 转为 sarama.Config
func (c *Config) ToSarama() (*sarama.Config, error) {
	cfg := &sarama.Config{}
	c.toAdmin(cfg)
	if err := c.toNet(cfg); err != nil {
		return nil, err
	}
	c.toMetadata(cfg)
	if err := c.toProducer(cfg); err != nil {
		return nil, err
	}
	if err := c.toConsumer(cfg); err != nil {
		return nil, err
	}
	cfg.ClientID = tools.Ternary(c.ClientID != "", c.ClientID, defaultClientID)
	cfg.ChannelBufferSize = tools.Ternary(c.ChannelBufferSize < 1, 256, c.ChannelBufferSize)
	cfg.ApiVersionsRequest = boolKit(c.ApiVersionsRequest, true)

	cfg.Version = sarama.DefaultVersion
	if c.Version != "" {
		var err error
		cfg.Version, err = sarama.ParseKafkaVersion(c.Version)
		if err != nil {
			return nil, fmt.Errorf("字段%s转Kafka版本失败：%v", "version", err)
		}
	}

	return cfg, nil
}

// 管理配置
func (c *Config) toAdmin(cfg *sarama.Config) {
	cfg.Admin.Retry.Max = tools.Ternary(c.Admin.Retry.Max < 0, 5, c.Admin.Retry.Max)
	cfg.Admin.Retry.Backoff = tools.AutoTimeDuration(c.Admin.Retry.Backoff, time.Millisecond, 100*time.Millisecond)
	cfg.Admin.Timeout = tools.AutoTimeDuration(c.Admin.Timeout, time.Second, 3*time.Second)
}

// 网络配置
func (c *Config) toNet(cfg *sarama.Config) error {
	cfg.Net.MaxOpenRequests = tools.Ternary(c.Net.MaxOpenRequests < 1, 5, c.Net.MaxOpenRequests)
	cfg.Net.DialTimeout = tools.AutoTimeDuration(c.Net.DialTimeout, time.Second, 30*time.Second)
	cfg.Net.ReadTimeout = tools.AutoTimeDuration(c.Net.ReadTimeout, time.Second, 30*time.Second)
	cfg.Net.WriteTimeout = tools.AutoTimeDuration(c.Net.WriteTimeout, time.Second, 30*time.Second)
	cfg.Net.ResolveCanonicalBootstrapServers = c.Net.ResolveCanonicalBootstrapServers

	if c.Net.TLS.Enable {
		var err error
		if cfg.Net.TLS.Config, err = c.Net.TLS.ToTLSConfig(); err != nil {
			return err
		}
		cfg.Net.TLS.Enable = true
	}

	cfg.Net.SASL.Enable = c.Net.SASL.Enable
	cfg.Net.SASL.Mechanism = c.Net.SASL.Mechanism

	switch cfg.Net.SASL.Mechanism {
	case sarama.SASLTypeSCRAMSHA256:
		cfg.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient { return &scram.XDGSCRAMClient{HashGeneratorFcn: scram.SHA256} }
	case sarama.SASLTypeSCRAMSHA512:
		cfg.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient { return &scram.XDGSCRAMClient{HashGeneratorFcn: scram.SHA512} }
	}

	cfg.Net.SASL.Version = sarama.SASLHandshakeV1
	if c.Net.SASL.Version != "" {
		v, err := tools.Any2int(c.Net.SASL.Version)
		if err != nil {
			return fmt.Errorf("字段%s转int失败：%v", "sasl.version", err)
		}
		cfg.Net.SASL.Version = int16(v)
	}

	cfg.Net.SASL.Handshake = boolKit(c.Net.SASL.Handshake, true)

	cfg.Net.SASL.AuthIdentity = c.Net.SASL.AuthIdentity
	cfg.Net.SASL.User = c.Net.SASL.User
	cfg.Net.SASL.Password = c.Net.SASL.Password
	cfg.Net.SASL.SCRAMAuthzID = c.Net.SASL.SCRAMAuthzID
	cfg.Net.SASL.GSSAPI = c.Net.SASL.GSSAPI.ToGSSAPIConfig()
	cfg.Net.KeepAlive = c.Net.KeepAlive
	if c.Net.Proxy != "" {
		py := cfg_proxy.Proxy{Addr: c.Net.Proxy}
		err := py.Valid()
		if err != nil {
			return err
		}
		cfg.Net.Proxy.Enable = true
		cfg.Net.Proxy.Dialer, err = py.AutoDialer()
		if err != nil {
			return fmt.Errorf("字段%s代理配置失败：%v", "net.proxy", err)
		}

	}
	return nil
}

// 元数据配置
func (c *Config) toMetadata(cfg *sarama.Config) {
	cfg.Metadata.Retry.Max = tools.Ternary(c.Metadata.Retry.Max < 1, 3, c.Metadata.Retry.Max)
	cfg.Metadata.Retry.Backoff = tools.AutoTimeDuration(c.Metadata.Retry.Backoff, time.Millisecond, 250*time.Millisecond)
	cfg.Metadata.RefreshFrequency = tools.AutoTimeDuration(c.Metadata.RefreshFrequency, time.Second, 10*time.Minute)
	cfg.Metadata.Timeout = c.Metadata.Timeout

	cfg.Metadata.AllowAutoTopicCreation = boolKit(c.Metadata.AllowAutoTopicCreation, true)
	cfg.Metadata.Full = boolKit(c.Metadata.Full, true)
	cfg.Metadata.SingleFlight = boolKit(c.Metadata.SingleFlight, true)

}

// 生产者配置
// noinspection all
func (c *Config) toProducer(cfg *sarama.Config) error {
	cfg.Producer.MaxMessageBytes = tools.Ternary(c.Producer.MaxMessageBytes < 1, 1024*1024, c.Producer.MaxMessageBytes)
	cfg.Producer.RequiredAcks = c.Producer.RequiredAcks
	cfg.Producer.Timeout = tools.AutoTimeDuration(c.Producer.Timeout, time.Second, 10*time.Second)
	cfg.Producer.Partitioner = sarama.NewRandomPartitioner
	cfg.Producer.Compression = c.Producer.Compression
	cfg.Producer.CompressionLevel = sarama.CompressionLevelDefault
	if c.Producer.CompressionLevel != "" {
		level, err := tools.Any2int(c.Producer.CompressionLevel)
		if err != nil {
			return fmt.Errorf("字段%s转int失败：%v", "producer.compression_level", err)
		}
		cfg.Producer.CompressionLevel = int(level)
	}
	cfg.Producer.Idempotent = c.Producer.Idempotent

	cfg.Producer.Transaction.ID = c.Producer.Transaction.ID
	cfg.Producer.Transaction.Timeout = tools.AutoTimeDuration(c.Producer.Transaction.Timeout, time.Second, time.Minute)
	cfg.Producer.Transaction.Retry.Max = tools.Ternary(c.Producer.Transaction.Retry.Max < 1, 50, c.Producer.Transaction.Retry.Max)
	cfg.Producer.Transaction.Retry.Backoff = tools.AutoTimeDuration(c.Producer.Transaction.Retry.Backoff, time.Millisecond, 100*time.Millisecond)

	cfg.Producer.Return.Errors = c.Producer.Return.Errors
	cfg.Producer.Return.Successes = c.Producer.Return.Successes

	cfg.Producer.Flush.Bytes = c.Producer.Flush.Bytes
	cfg.Producer.Flush.Messages = c.Producer.Flush.Messages
	cfg.Producer.Flush.Frequency = c.Producer.Flush.Frequency
	cfg.Producer.Flush.MaxMessages = c.Producer.Flush.MaxMessages

	cfg.Producer.Retry.Max = tools.Ternary(c.Producer.Retry.Max < 1, 3, c.Producer.Retry.Max)
	cfg.Producer.Retry.Backoff = tools.AutoTimeDuration(c.Producer.Retry.Backoff, time.Millisecond, 100*time.Millisecond)
	cfg.Producer.Retry.MaxBufferBytes = c.Producer.Retry.MaxBufferBytes
	cfg.Producer.Retry.MaxBufferLength = c.Producer.Retry.MaxBufferLength

	return nil

}

func (c *Config) toConsumer(cfg *sarama.Config) error {
	cfg.Consumer.Group.Session.Timeout = tools.AutoTimeDuration(c.Consumer.Group.Session.Timeout, time.Second, 10*time.Second)

	cfg.Consumer.Group.Heartbeat.Interval = tools.AutoTimeDuration(c.Consumer.Group.Heartbeat.Interval, time.Second, 3*time.Second)

	// 消费组平衡模式
	cfg.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRange()}
	if len(c.Consumer.Group.Rebalance.GroupStrategies) > 0 {
		var strategies []sarama.BalanceStrategy
		for _, strategy := range c.Consumer.Group.Rebalance.GroupStrategies {
			switch strategy {
			case sarama.RangeBalanceStrategyName:
				strategies = append(strategies, sarama.NewBalanceStrategyRange())
			case sarama.RoundRobinBalanceStrategyName:
				strategies = append(strategies, sarama.NewBalanceStrategyRoundRobin())
			case sarama.StickyBalanceStrategyName:
				strategies = append(strategies, sarama.NewBalanceStrategySticky())
			}
		}
		if len(strategies) > 0 {
			cfg.Consumer.Group.Rebalance.GroupStrategies = strategies
		}
	}
	cfg.Consumer.Group.Rebalance.Timeout = tools.AutoTimeDuration(c.Consumer.Group.Rebalance.Timeout, time.Second, 60*time.Second)
	cfg.Consumer.Group.Rebalance.Retry.Max = tools.Ternary(c.Consumer.Group.Rebalance.Retry.Max < 1, 4, c.Consumer.Group.Rebalance.Retry.Max)
	cfg.Consumer.Group.Rebalance.Retry.Backoff = tools.AutoTimeDuration(c.Consumer.Group.Rebalance.Retry.Backoff, time.Second, 2*time.Second)

	cfg.Consumer.Group.InstanceId = c.Consumer.Group.InstanceId
	cfg.Consumer.Group.ResetInvalidOffsets = boolKit(c.Consumer.Group.ResetInvalidOffsets, true)

	cfg.Consumer.Retry.Backoff = tools.AutoTimeDuration(c.Consumer.Retry.Backoff, time.Second, 2*time.Second)

	cfg.Consumer.Fetch.Min = tools.Ternary(c.Consumer.Fetch.Min < 1, 1, c.Consumer.Fetch.Min)
	cfg.Consumer.Fetch.Default = tools.Ternary(c.Consumer.Fetch.Default < 1, 1024*1024, c.Consumer.Fetch.Default)
	cfg.Consumer.Fetch.Max = c.Consumer.Fetch.Max

	cfg.Consumer.MaxWaitTime = tools.AutoTimeDuration(c.Consumer.MaxWaitTime, time.Millisecond, 500*time.Millisecond)
	cfg.Consumer.MaxProcessingTime = tools.AutoTimeDuration(c.Consumer.MaxProcessingTime, time.Millisecond, 100*time.Millisecond)
	cfg.Consumer.Return.Errors = c.Consumer.Return.Errors

	cfg.Consumer.Offsets.AutoCommit.Enable = boolKit(c.Consumer.Offsets.AutoCommit.Enable, true)
	cfg.Consumer.Offsets.AutoCommit.Interval = tools.AutoTimeDuration(c.Consumer.Offsets.AutoCommit.Interval, time.Second, time.Second)

	cfg.Consumer.Offsets.Initial = sarama.OffsetNewest
	if c.Consumer.Offsets.Initial != "" {
		initial, err := tools.Any2int(c.Consumer.Offsets.Initial)
		if err != nil {
			return fmt.Errorf("字段%s转int失败：%v", "consumer.offsets.initial", err)
		}
		cfg.Consumer.Offsets.Initial = initial
	}
	cfg.Consumer.Offsets.Retention = c.Consumer.Offsets.Retention
	cfg.Consumer.Offsets.Retry.Max = tools.Ternary(c.Consumer.Offsets.Retry.Max < 1, 3, c.Consumer.Offsets.Retry.Max)
	cfg.Consumer.IsolationLevel = c.Consumer.IsolationLevel
	cfg.MetricRegistry = metrics.NewRegistry()
	return nil

}
