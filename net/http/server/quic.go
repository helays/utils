package server

import "time"

type QUICConfig struct {
	HandshakeIdleTimeout             time.Duration `json:"handshake_idle_timeout" yaml:"handshake_idle_timeout"`                             // 握手阶段的空闲超时
	MaxIdleTimeout                   time.Duration `json:"max_idle_timeout" yaml:"max_idle_timeout"`                                         // 连接建立后的最大空闲超时
	InitialStreamReceiveWindow       uint64        `json:"initial_stream_receive_window" yaml:"initial_stream_receive_window"`               // 初始流接收窗口大小（默认 512KB）
	MaxStreamReceiveWindow           uint64        `json:"max_stream_receive_window" yaml:"max_stream_receive_window"`                       // 最大流接收窗口大小（默认 6MB）
	InitialConnectionReceiveWindow   uint64        `json:"initial_connection_receive_window" yaml:"initial_connection_receive_window"`       // 初始连接接收窗口大小（默认 512KB）
	MaxConnectionReceiveWindow       uint64        `json:"max_connection_receive_window" yaml:"max_connection_receive_window"`               // 最大连接接收窗口大小（默认 15MB）
	MaxIncomingStreams               int64         `json:"max_incoming_streams" yaml:"max_incoming_streams"`                                 // 允许对端同时打开的双向流最大数量
	MaxIncomingUniStreams            int64         `json:"max_incoming_uni_streams" yaml:"max_incoming_uni_streams"`                         // 允许对端同时打开的单向流最大数量
	KeepAlivePeriod                  time.Duration `json:"keep_alive_period" yaml:"keep_alive_period"`                                       // 发送保活数据包的时间间隔
	InitialPacketSize                uint16        `json:"initial_packet_size" yaml:"initial_packet_size"`                                   // 发送数据包的初始大小（也是下限）,通常自动进行 PMTU 发现
	DisablePathMTUDiscovery          bool          `json:"disable_path_mtu_discovery" yaml:"disable_path_mtu_discovery"`                     // 禁用路径 MTU 发现（RFC 8899）
	Allow0RTT                        bool          `json:"allow_0_rtt" yaml:"allow_0_rtt"`                                                   // 允许接受 0-RTT 连接尝试
	EnableDatagrams                  bool          `json:"enable_datagrams" yaml:"enable_datagrams"`                                         // 启用 QUIC 数据报支持（RFC 9221）
	EnableStreamResetPartialDelivery bool          `json:"enable_stream_reset_partial_delivery" yaml:"enable_stream_reset_partial_delivery"` // 启用带部分交付的流重置
}
