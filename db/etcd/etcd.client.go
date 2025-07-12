package etcd

import (
	"context"
	tlsconfig "github.com/helays/utils/v2/crypto/tls.config"
	"github.com/helays/utils/v2/tools"
	"go.etcd.io/etcd/client/v3"
	"time"
)

type Config struct {
	Endpoints             []string         `json:"endpoints" yaml:"endpoints" ini:"endpoints,omitempty"`                                 // etcd地址
	DialTimeout           time.Duration    `json:"dial_timeout" yaml:"dial_timeout" ini:"dial_timeout"`                                  // 连接超时时间
	AutoSyncInterval      time.Duration    `json:"auto_sync_interval" yaml:"auto_sync_interval" ini:"auto_sync_interval"`                // 自动同步间隔
	DialKeepAliveTime     time.Duration    `json:"dial_keep_alive_time" yaml:"dial_keep_alive_time" ini:"dial_keep_alive_time"`          // 连接保活时间，客户端与服务器之间的 keep-alive 探测间隔时间。设置为 0 表示禁用 keep-alive。
	DialKeepAliveTimeout  time.Duration    `json:"dial_keep_alive_timeout" yaml:"dial_keep_alive_timeout" ini:"dial_keep_alive_timeout"` // keep-alive 探测的超时时间。如果在这个时间内没有收到服务器的响应，连接将被关闭。
	MaxCallSendMsgSize    int              `json:"max_call_send_msg_size" yaml:"max_call_send_msg_size" ini:"max_call_send_msg_size"`    // 客户端发送消息的最大大小。设置为 0 表示没有限制。
	MaxCallRecvMsgSize    int              `json:"max_call_recv_msg_size" yaml:"max_call_recv_msg_size" ini:"max_call_recv_msg_size"`    // 客户端接收消息的最大大小。设置为 0 表示没有限制。
	EnableTLS             bool             `json:"enable_tls" yaml:"enable_tls" ini:"enable_tls"`                                        // 是否使用TLS
	TLS                   tlsconfig.Config `json:"tls" yaml:"tls" ini:"tls"`                                                             // TLS配置
	Username              string           `json:"username" yaml:"username" ini:"username"`
	Password              string           `json:"password" yaml:"password" ini:"password"`
	RejectOldCluster      bool             `json:"reject_old_cluster" yaml:"reject_old_cluster" ini:"reject_old_cluster"`                // 拒绝旧集群
	PermitWithoutStream   bool             `json:"permit_without_stream" yaml:"permit_without_stream" ini:"permit_without_stream"`       // 如果设置为 true，即使没有活跃的流，客户端也会保持连接。
	MaxUnaryRetries       uint             `json:"max_unary_retries" yaml:"max_unary_retries" ini:"max_unary_retries"`                   // 单次 gRPC 调用的最大重试次数。
	BackoffWaitBetween    time.Duration    `json:"backoff_wait_between" yaml:"backoff_wait_between" ini:"backoff_wait_between"`          // 重试之间的等待时间
	BackoffJitterFraction float64          `json:"backoff_jitter_fraction" yaml:"backoff_jitter_fraction" ini:"backoff_jitter_fraction"` // 重试等待时间的抖动因子，用于避免重试风暴。
}

// NewClient 创建etcd客户端
func (this Config) NewClient(ctx context.Context) (*clientv3.Client, error) {
	c := clientv3.Config{
		Endpoints:             this.Endpoints,
		AutoSyncInterval:      tools.AutoTimeDuration(this.AutoSyncInterval, time.Second),
		DialTimeout:           tools.AutoTimeDuration(this.DialTimeout, time.Second, 10*time.Second),
		DialKeepAliveTime:     tools.AutoTimeDuration(this.DialKeepAliveTime, time.Second),
		DialKeepAliveTimeout:  tools.AutoTimeDuration(this.DialKeepAliveTimeout, time.Second),
		MaxCallSendMsgSize:    this.MaxCallSendMsgSize,
		MaxCallRecvMsgSize:    this.MaxCallRecvMsgSize,
		Username:              this.Username,
		Password:              this.Password,
		RejectOldCluster:      this.RejectOldCluster,
		Context:               ctx,
		PermitWithoutStream:   this.PermitWithoutStream,
		MaxUnaryRetries:       this.MaxUnaryRetries,
		BackoffWaitBetween:    tools.AutoTimeDuration(this.BackoffWaitBetween, time.Second),
		BackoffJitterFraction: this.BackoffJitterFraction,
	}

	// 使用tls
	if this.EnableTLS {
		// 使用tls
		tlsConfig, err := this.TLS.NewTLSConfig()
		if err != nil {
			return nil, err
		}
		c.TLS = tlsConfig
	}
	return clientv3.New(c)
}
