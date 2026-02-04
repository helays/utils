package clientconf

import (
	"database/sql/driver"
	"time"

	"github.com/IBM/sarama"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"helay.net/go/utils/v3/dataType"
	"helay.net/go/utils/v3/net/tlsconfig"
)

type Retry struct {
	Max     int           `yaml:"max" json:"max"`         // 最大重试次数
	Backoff time.Duration `yaml:"backoff" json:"backoff"` // 失败请求重试之间的退避时间，默认100ms。
	// AsyncProducer#retryHandler 中 `input` 和 `retries` 通道之间
	// 桥接缓冲区的最大长度。
	// 限制是为了防止此缓冲区溢出或导致 OOM。
	// 默认为 0 表示无限制。
	// 0 到 4096 之间的任何值都会推送到 4096。
	// 零或负值表示无限制。
	MaxBufferLength int `yaml:"max_buffer_length" json:"max_buffer_length"`
	// AsyncProducer#retryHandler 中 `input` 和 `retries` 通道之间
	// 桥接缓冲区中消息的最大总字节大小。
	// 此限制防止缓冲区消耗过多内存。
	// 默认为 0 表示无限制。
	// 0 到 32 MB 之间的任何值都会推送到 32 MB。
	// 零或负值表示无限制。
	MaxBufferBytes int64 `yaml:"max_buffer_bytes" json:"max_buffer_bytes"`
}

// Admin 是管理性 Kafka 客户端适用的 ClusterAdmin 属性的命名空间。
type Admin struct {
	Retry   Retry         `yaml:"retry" json:"retry"`
	Timeout time.Duration `yaml:"timeout" json:"timeout"`
}
type GSSAPIConfig struct {
	AuthType           int    `yaml:"auth_type" json:"auth_type"`
	KeyTabPath         string `yaml:"key_tab_path" json:"key_tab_path"`
	CCachePath         string `yaml:"ccache_path" json:"ccache_path"`
	KerberosConfigPath string `yaml:"kerberos_config_path" json:"kerberos_config_path"`
	ServiceName        string `yaml:"service_name" json:"service_name"`
	Username           string `yaml:"username" json:"username"`
	Password           string `yaml:"password" json:"password"`
	Realm              string `yaml:"realm" json:"realm"`
	DisablePAFXFAST    bool   `yaml:"disable_pafxfast" json:"disable_pafxfast"`
}

// SASL 配置
// 基于 SASL 的代理身份验证。虽然有多种 SASL 认证方法，
// 但当前实现仅限于明文（SASL/PLAIN）认证。
// noinspection all
type SASL struct {
	Enable bool `yaml:"enable" json:"enable"` // 是否在连接到代理时使用 SASL 认证（默认为 false）。
	// SASLMechanism 是启用的 SASL 机制的名称。
	// 可能的值：OAUTHBEARER, PLAIN（默认为 PLAIN）。
	Mechanism sarama.SASLMechanism `yaml:"mechanism" json:"mechanism"`
	// Version 是使用的 SASL 协议版本
	// Kafka > 1.x 应使用 V1，除了 Azure EventHub 使用 V0
	Version string `yaml:"version" json:"version"`
	// 如果启用，是否首先发送 Kafka SASL 握手
	//（默认为 true）。只有在使用非 Kafka SASL 代理时才应将其设置为 false。
	Handshake string `yaml:"handshake" json:"handshake"`
	// AuthIdentity 是用于 SASL/PLAIN 认证的可选授权身份（authzid）
	//（如果与 User 不同），当认证用户被允许充当提供的替代用户时。
	// 详见 RFC4616。
	AuthIdentity string `yaml:"auth_identity" json:"auth_identity"`
	// User 是用于 SASL/PLAIN 或 SASL/SCRAM 认证的身份（authcid）。
	User string `yaml:"user" json:"user"`
	// 用于 SASL/PLAIN 认证的密码
	Password string `yaml:"password" json:"password"`
	// 用于 SASL/SCRAM 认证的 authz id
	SCRAMAuthzID string       `yaml:"scram_authz_id" json:"scram_authz_id"`
	GSSAPI       GSSAPIConfig `yaml:"gssapi" json:"gssapi"`
}

// Net 是 Broker 使用的网络级别属性的命名空间，
// 由 Client/Producer/Consumer 共享。
type Net struct {
	// 连接在阻塞发送之前允许的未完成请求数量（默认 5）。
	// 如果 Producer.Idempotent 被禁用，可以提高吞吐量但不保证消息顺序，参见：
	// https://kafka.apache.org/protocol#protocol_network
	// https://kafka.apache.org/28/documentation.html#producerconfigs_max.in.flight.requests.per.connection
	MaxOpenRequests int `yaml:"max_open_requests" json:"max_open_requests"`
	// 以下三个配置类似于 JVM kafka 中的 `socket.timeout.ms` 设置。
	// 它们都默认为 30 秒。
	DialTimeout  time.Duration `yaml:"dial_timeout" json:"dial_timeout"`   // 等待初始连接的时间。
	ReadTimeout  time.Duration `yaml:"read_timeout" json:"read_timeout"`   // 等待响应的时间。
	WriteTimeout time.Duration `yaml:"write_timeout" json:"write_timeout"` // 等待传输的时间。
	// ResolveCanonicalBootstrapServers 将每个引导代理地址
	// 转换为一组 IP，然后对每个 IP 进行反向查找以获取其
	// 规范主机名。这个主机名列表随后替换原始地址列表。
	// 类似于 JVM 客户端中的 `client.dns.lookup` 选项，
	// 这在 GSSAPI 中特别有用，因为它允许提供别名记录
	// 而不是单个代理主机名。默认为 false。
	ResolveCanonicalBootstrapServers bool `yaml:"resolve_canonical_bootstrap_servers" json:"resolve_canonical_bootstrap_servers"`
	// tls 配置
	TLS  tlsconfig.TLSConfig `yaml:"tls" json:"tls"`
	SASL SASL                `yaml:"sasl" json:"sasl"`
	// KeepAlive 指定活动网络连接的保活周期（默认为 0）。
	// 如果为零或正数，则启用保活。
	// 如果为负数，则禁用保活。
	KeepAlive time.Duration `yaml:"keep_alive" json:"keep_alive"`
	Proxy     string        `yaml:"proxy" json:"proxy"`
}

// Metadata 是 Client 使用的元数据管理属性的命名空间，
// 由 Producer/Consumer 共享。
type Metadata struct {
	Retry Retry `yaml:"retry" json:"retry"`
	// 在后台刷新集群元数据的频率。
	// 默认为 10 分钟。设置为 0 以禁用。类似于
	// JVM 版本中的 `topic.metadata.refresh.interval.ms`。
	RefreshFrequency time.Duration `yaml:"refresh_frequency" json:"refresh_frequency"`
	// 是维护所有主题的完整元数据集，还是仅维护
	// 迄今为止所需的最小集。完整集更简单
	// 且通常更方便，但如果有很多主题和分区，
	// 可能会占用大量内存。默认为 true。
	Full string `yaml:"full" json:"full"`
	// 等待成功元数据响应的时间。
	// 默认禁用，这意味着对不可达集群的元数据请求
	//（所有代理都不可达或无响应）最多可能需要
	// `Net.[Dial|Read]Timeout * BrokerCount * (Metadata.Retry.Max + 1) + Metadata.Retry.Backoff * Metadata.Retry.Max`
	// 才会失败。
	Timeout time.Duration `yaml:"timeout" json:"timeout"`
	// 是否允许在元数据刷新中自动创建主题。如果设置为 true，
	// 代理可能会自动创建我们请求的尚不存在的主题，
	// 如果它被配置为这样做（`auto.create.topics.enable` 为 true）。默认为 true。
	AllowAutoTopicCreation string `yaml:"allow_auto_topic_creation" json:"allow_auto_topic_creation"`
	// SingleFlight 控制是在给定时间发送单个元数据刷新请求，
	// 还是允许任何人并发刷新元数据。
	// 如果设置为 true，并且客户端需要从不同的 goroutine 刷新元数据，
	// 请求将被批处理，以便一次发送单个刷新。
	// 更多详情请参见 https://github.com/IBM/sarama/issues/3224。
	// SingleFlight 默认为 true。
	SingleFlight string `yaml:"single_flight" json:"single_flight"`
}

type Transaction struct {
	// 在事务中用于通过重启识别生产者实例
	ID string `yaml:"id" json:"id"`
	// 事务可以保持未解决（既未提交也未中止）的时间量
	// 默认为 1 分钟
	Timeout time.Duration `yaml:"timeout" json:"timeout"`
	Retry   Retry         `yaml:"retry" json:"retry"`
}

// Return 指定将填充哪些通道。如果它们设置为 true，
// 您必须从相应的通道读取以防止死锁。然而，
// 如果此配置用于创建 `SyncProducer`，则两者都必须设置为 true，
// 并且您不应从通道读取，因为生产者会在内部处理。
type Return struct {
	Successes bool `yaml:"successes" json:"successes"` // 如果启用，成功传递的消息将在 Successes 通道上返回（默认禁用）。
	Errors    bool `yaml:"errors" json:"errors"`       // 如果启用，传递失败的消息将在 Errors 通道上返回，包括错误（默认启用）。
}

// Flush 以下配置选项控制消息批处理和发送到代理的频率。
// 默认情况下，消息尽可能快地发送，并且在当前批次传输过程中
// 接收到的所有消息都放入后续批次中。
type Flush struct {
	// 触发刷新的最佳努力字节数。使用
	// 全局 sarama.MaxRequestSize 设置硬上限。
	Bytes int `yaml:"bytes" json:"bytes"`
	// 触发刷新的最佳努力消息数。使用
	// `MaxMessages` 设置硬上限。
	Messages int `yaml:"messages" json:"messages"`
	// 刷新的最佳努力频率。相当于
	// JVM producer 的 `queue.buffering.max.ms` 设置。
	Frequency time.Duration `yaml:"frequency" json:"frequency"`
	// 生产者在单个代理请求中发送的最大消息数。
	// 默认为 0 表示无限制。类似于
	// JVM producer 中的 `queue.buffering.max.messages`。
	MaxMessages int `yaml:"max_messages" json:"max_messages"`
}

// Producer 是与生产消息相关的配置的命名空间，
// 由 Producer 使用。
type Producer struct {
	// 消息的最大允许大小（默认为 1000000）。应
	// 设置等于或小于代理的 `message.max.bytes`。
	MaxMessageBytes int `yaml:"max_message_bytes" json:"max_message_bytes"`
	// 需要从代理获得的确认可靠性级别（默认为 WaitForLocal）。
	// 相当于 JVM producer 的 `request.required.acks` 设置。
	RequiredAcks sarama.RequiredAcks `yaml:"required_acks" json:"required_acks"`
	// 代理等待收到 RequiredAcks 数量的最长时间
	//（默认为 10 秒）。这仅在
	// RequiredAcks 设置为 WaitForAll 或数字 > 1 时相关。仅支持
	// 毫秒分辨率，纳秒将被截断。相当于
	// JVM producer 的 `request.timeout.ms` 设置。
	Timeout time.Duration `yaml:"timeout" json:"timeout"`
	// 用于消息的压缩类型（默认为无压缩）。
	// 类似于 JVM producer 的 `compression.codec` 设置。
	Compression sarama.CompressionCodec `yaml:"compression" json:"compression"`
	// 用于消息的压缩级别。含义取决于
	// 实际使用的压缩类型，默认为编解码器的默认压缩级别。
	CompressionLevel string `yaml:"compression_level" json:"compression_level"`
	// 如果启用，生产者将确保每条消息恰好写入一份副本。
	Idempotent bool `yaml:"idempotent" json:"idempotent"`
	// Transaction 指定
	Transaction Transaction `yaml:"transaction" json:"transaction"`
	Return      Return      `yaml:"return" json:"return"`
	Flush       Flush       `yaml:"flush" json:"flush"`
	Retry       Retry       `yaml:"retry" json:"retry"`
}

type Session struct {
	// 当使用 Kafka 的组管理功能时，用于检测消费者失败的超时时间。
	// 消费者定期发送心跳以向代理指示其活跃状态。
	// 如果在此会话超时到期之前代理未收到心跳，
	// 则代理将从组中移除该消费者并启动重新平衡。
	// 注意，该值必须在代理配置中允许的范围内，
	// 由 `group.min.session.timeout.ms` 和 `group.max.session.timeout.ms` 配置（默认 10 秒）
	Timeout time.Duration `yaml:"timeout" json:"timeout"`
}

type Heartbeat struct {
	// 当使用 Kafka 的组管理功能时，向消费者协调器发送心跳的预期间隔时间。
	// 心跳用于确保消费者的会话保持活跃，并促进新消费者加入或离开组时的重新平衡。
	// 该值必须设置为低于 Consumer.Group.Session.Timeout，但通常不应高于该值的三分之一。
	// 可以设置得更低以控制正常重新平衡的预期时间（默认 3 秒）
	Interval time.Duration `yaml:"interval" json:"interval"`
}

type Rebalance struct {
	// GroupStrategies 是客户端消费者组平衡策略的优先级有序列表，
	// 将提供给协调器。所有组成员支持的第一个策略将由领导者选择。
	// 默认：[ NewBalanceStrategyRange() ]
	GroupStrategies []string `yaml:"group_strategies" json:"group_strategies"`
	// 重新平衡开始后，每个工作线程加入组的最大允许时间。
	// 这基本上是所有任务刷新任何待处理数据和提交偏移量所需的时间限制。
	// 如果超时，工作线程将从组中移除，这将导致偏移量提交失败（默认 60 秒）。
	Timeout time.Duration `yaml:"timeout" json:"timeout"`
	Retry   Retry         `yaml:"retry" json:"retry"`
}

type Group struct {
	Session   Session   `yaml:"session" json:"session"`
	Heartbeat Heartbeat `yaml:"heartbeat" json:"heartbeat"`
	Rebalance Rebalance `yaml:"rebalance" json:"rebalance"`
	// 支持 KIP-345
	InstanceId string `yaml:"instance_id" json:"instance_id"`
	// 如果为 true，当获取的消费者偏移量超出可用偏移量范围时，
	// 消费者偏移量将自动重置为配置的初始值。超出范围
	// 可能发生在数据已从服务器删除，或者在副本尚未拥有所有数据的
	// 复制不足情况下。自动重置偏移量可能很危险，尤其是在后一种情况下。
	// 默认为 true 以保持现有行为。
	ResetInvalidOffsets string `yaml:"reset_invalid_offsets" json:"reset_invalid_offsets"`
}

// Fetch 是控制每次请求检索多少字节的命名空间。
type Fetch struct {
	// 请求中获取的最小消息字节数 - 代理
	// 将等待直到至少有这么多字节可用。默认值为 1，
	// 因为 0 会导致在没有消息可用时消费者旋转。
	// 相当于 JVM 的 `fetch.min.bytes`。
	Min int32 `yaml:"min" json:"min"`
	// 每次请求中从代理获取的默认消息字节数（默认 1MB）。
	// 这应大于大多数消息的大小，否则消费者将花费大量时间
	// 协商大小而不是实际消费。类似于 JVM 的 `fetch.message.max.bytes`。
	Default int32 `yaml:"default" json:"default"`
	// 单次请求中从代理获取的最大消息字节数。
	// 大于此值的消息将返回 ErrMessageTooLarge 且不可消费，
	// 因此您必须确保此值至少与最大消息一样大。默认为 0（无限制）。
	// 类似于 JVM 的 `fetch.message.max.bytes`。
	// 全局 `sarama.MaxResponseSize` 仍然适用。
	Max int32 `yaml:"max" json:"max"`
}

type AutoCommit struct {
	// 是否自动将更新的偏移量提交回代理。
	//（默认启用）。
	Enable string `yaml:"enable" json:"enable"`
	// 提交更新偏移量的频率。除非启用自动提交，否则无效（默认 1 秒）
	Interval time.Duration `yaml:"interval" json:"interval"`
}

// Offsets 指定如何以及何时提交消费的偏移量的配置。
// 这当前需要手动使用 OffsetManager，但最终将自动化。
type Offsets struct {
	// AutoCommit 指定自动提交消息的配置。
	AutoCommit AutoCommit `yaml:"auto_commit" json:"auto_commit"`
	// 如果之前没有提交偏移量，则使用的初始偏移量。
	// 应为 OffsetNewest 或 OffsetOldest。默认为 OffsetNewest。
	Initial string `yaml:"initial" json:"initial"`
	// 已提交偏移量的保留时间。如果为零，则禁用
	//（此时将使用代理上的 `offsets.retention.minutes` 选项）。
	// Kafka 仅支持到毫秒的精度；纳秒将被截断。需要 Kafka 代理版本 0.9.0 或更高。
	//（默认为 0：禁用）。
	Retention time.Duration `yaml:"retention" json:"retention"`
	Retry     Retry         `yaml:"retry" json:"retry"`
}

// Consumer 是与消费消息相关的配置的命名空间，
// 由 Consumer 使用。
type Consumer struct {
	// Group 是配置消费者组的命名空间。
	Group Group `yaml:"group" json:"group"`
	Retry Retry `yaml:"retry" json:"retry"`
	Fetch Fetch `yaml:"fetch" json:"fetch"`
	// 代理在返回少于 Consumer.Fetch.Min 字节之前等待其变为可用的最长时间。
	// 默认值为 250ms，因为 0 会导致在没有事件可用时消费者旋转。
	// 对于大多数情况，100-500ms 是一个合理的范围。Kafka 仅支持到毫秒的精度；
	// 纳秒将被截断。相当于 JVM 的 `fetch.max.wait.ms`。
	MaxWaitTime time.Duration `yaml:"max_wait_time" json:"max_wait_time"`
	// 消费者预期用户处理消息所需的最长时间。
	// 如果写入 Messages 通道的时间超过此时间，则该分区将停止获取更多消息，
	// 直到可以继续。
	// 注意，由于 Messages 通道是缓冲的，实际宽限时间为
	// (MaxProcessingTime * ChannelBufferSize)。默认为 100ms。
	// 如果在 expiryTicker 的两个计时周期之间没有消息写入 Messages 通道，
	// 则检测到超时。
	// 使用计时器而不是计时器来检测超时通常会导致对计时器函数的调用大大减少，
	// 如果发送大量消息且超时不频繁，这可能会显著提高性能。
	// 使用计时器而不是计时器的缺点是超时精度较低。
	// 也就是说，有效超时可能在 `MaxProcessingTime` 和 `2 * MaxProcessingTime` 之间。
	// 例如，如果 `MaxProcessingTime` 为 100ms，则两次发送消息之间 180ms 的延迟
	// 可能不会被识别为超时。
	MaxProcessingTime time.Duration `yaml:"max_processing_time" json:"max_processing_time"`
	// Return 指定将填充哪些通道。如果它们设置为 true，
	// 您必须从中读取以防止死锁。
	Return  Return  `yaml:"return" json:"return"`
	Offsets Offsets `yaml:"offsets" json:"offsets"`
	// IsolationLevel 支持 2 种模式：
	// 	- 使用 `ReadUncommitted`（默认）消费并返回消息通道中的所有消息
	//	- 使用 `ReadCommitted` 隐藏属于已中止事务的消息
	IsolationLevel sarama.IsolationLevel `yaml:"isolation_level" json:"isolation_level"`
}

// noinspection all
type Config struct {
	Admin    Admin    `yaml:"admin" json:"admin"`
	Net      Net      `yaml:"net" json:"net"`
	Metadata Metadata `yaml:"metadata" json:"metadata"`
	Producer Producer `yaml:"producer" json:"producer"`
	Consumer Consumer `yaml:"consumer" json:"consumer"`
	// 用户提供的字符串，随每个请求发送给代理，用于日志记录、
	// 调试和审计目的。默认为 "sarama"，但您应该
	// 可能将其设置为特定于应用程序的内容。
	ClientID string `yaml:"client_id" json:"client_id"`
	// 此客户端的机架标识符。这可以是任何字符串值，
	// 指示此客户端的物理位置。
	// 它与代理配置 'broker.rack' 对应
	RackID string `yaml:"rack_id" json:"rack_id"`
	// 内部和外部通道中缓冲的事件数。这允许生产者和消费者
	// 在用户代码工作时在后台继续处理一些消息，大大提高了吞吐量。
	// 默认为 256。
	ChannelBufferSize int `yaml:"channel_buffer_size" json:"channel_buffer_size"`
	// ApiVersionsRequest 决定 Sarama 是否应在其初始连接过程中
	// 向每个代理发送 ApiVersionsRequest 消息。这默认为 `true` 以匹配官方 Java 客户端
	// 和大多数第三方客户端。
	ApiVersionsRequest string `yaml:"api_versions_request" json:"api_versions_request"`
	// Sarama 将假定其运行的 Kafka 版本。
	// 默认为支持的最旧稳定版本。由于 Kafka 提供向后兼容性，
	// 将其设置为比您拥有的版本更旧的版本不会破坏任何内容，
	// 尽管它可能会阻止您使用最新功能。将其设置为比实际运行的版本更高的版本
	// 可能导致随机故障。
	Version string `yaml:"version" json:"version"`
}

// noinspection all
func (c Config) Value() (driver.Value, error) {
	return dataType.DriverValueWithJson(c)
}

// noinspection all
func (c *Config) Scan(val interface{}) error {
	return dataType.DriverScanWithJson(val, c)
}

func (c *Config) GormDataType() string {
	return "json"
}

// noinspection all
func (Config) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	return dataType.JsonDbDataType(db, field)
}
