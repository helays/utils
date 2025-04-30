package nacos

import (
	"github.com/helays/utils/tools"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"time"
)

type Config struct {
	ServerConfigs []ServerConfig `json:"server_configs" yaml:"server_configs" ini:"server_configs"` // 服务器地址
	ClientConfig  *ClientConfig  `json:"client_config" yaml:"client_config" ini:"client_config"`    // 客户端配置
}

type ServerConfig struct {
	Scheme      string `json:"scheme" yaml:"scheme" ini:"scheme"`                   // nacos 服务协议，默认 http（2.0 版本非必填）
	ContextPath string `json:"context_path" yaml:"context_path" ini:"context_path"` // nacos 服务上下文路径，默认 /nacos（2.0 版本非必填）
	IpAddr      string `json:"ip_addr" yaml:"ip_addr" ini:"ip_addr"`                // nacos 服务地址
	Port        uint64 `json:"port" yaml:"port" ini:"port"`                         // nacos 服务端口
	GrpcPort    uint64 `json:"grpc_port" yaml:"grpc_port" ini:"grpc_port"`          // nacos 服务 gRPC 端口，默认 server_port + 1000（非必填）
}

type ClientConfig struct {
	TimeoutMs            uint64                   `json:"timeout_ms" yaml:"timeout_ms" ini:"timeout_ms"`                                        // 请求 Nacos 服务器的超时时间，默认 10000ms
	BeatInterval         int64                    `json:"beat_interval" yaml:"beat_interval" ini:"beat_interval"`                               // 发送心跳到服务器的时间间隔，默认 5000ms
	NamespaceId          string                   `json:"namespace_id" yaml:"namespace_id" ini:"namespace_id"`                                  // Nacos 的命名空间 ID，公共命名空间请留空
	AppName              string                   `json:"app_name" yaml:"app_name" ini:"app_name"`                                              // 应用名称
	AppKey               string                   `json:"app_key" yaml:"app_key" ini:"app_key"`                                                 // 客户端身份信息
	Endpoint             string                   `json:"endpoint" yaml:"endpoint" ini:"endpoint"`                                              // 获取 Nacos 服务器地址的端点
	RegionId             string                   `json:"region_id" yaml:"region_id" ini:"region_id"`                                           // KMS 的区域 ID
	AccessKey            string                   `json:"access_key" yaml:"access_key" ini:"access_key"`                                        // KMS 的访问密钥
	SecretKey            string                   `json:"secret_key" yaml:"secret_key" ini:"secret_key"`                                        // KMS 的密钥
	RamConfig            *RamConfig               `json:"ram_config" yaml:"ram_config" ini:"ram_config"`                                        // RAM 配置
	OpenKMS              bool                     `json:"open_kms" yaml:"open_kms" ini:"open_kms"`                                              // 是否开启 KMS，默认 false
	KMSVersion           constant.KMSVersion      `json:"kms_version" yaml:"kms_version" ini:"kms_version"`                                     // KMS 客户端版本
	KMSv3Config          *KMSv3Config             `json:"kmsv3_config" yaml:"kmsv3_config" ini:"kmsv3_config"`                                  // KMSv3 配置
	KMSConfig            *KMSConfig               `json:"kms_config" yaml:"kms_config" ini:"kms_config"`                                        // KMS 配置
	CacheDir             string                   `json:"cache_dir" yaml:"cache_dir" ini:"cache_dir"`                                           // 持久化 nacos 服务信息的目录，默认当前路径
	DisableUseSnapShot   bool                     `json:"disable_use_snap_shot" yaml:"disable_use_snap_shot" ini:"disable_use_snap_shot"`       // 开关，默认 false，表示当获取远程配置失败时使用本地缓存文件
	UpdateThreadNum      int                      `json:"update_thread_num" yaml:"update_thread_num" ini:"update_thread_num"`                   // 更新 nacos 服务信息的 goroutine 数量，默认 20
	NotLoadCacheAtStart  bool                     `json:"not_load_cache_at_start" yaml:"not_load_cache_at_start" ini:"not_load_cache_at_start"` // 启动时不加载 CacheDir 中的持久化 nacos 服务信息
	UpdateCacheWhenEmpty bool                     `json:"update_cache_when_empty" yaml:"update_cache_when_empty" ini:"update_cache_when_empty"` // 当从服务器获取空服务实例时更新缓存
	Username             string                   `json:"username" yaml:"username" ini:"username"`                                              // nacos 认证用户名
	Password             string                   `json:"password" yaml:"password" ini:"password"`                                              // nacos 认证密码
	LogDir               string                   `json:"log_dir" yaml:"log_dir" ini:"log_dir"`                                                 // 日志目录，默认当前路径
	LogLevel             string                   `json:"log_level" yaml:"log_level" ini:"log_level"`                                           // 日志级别，必须是 debug/info/warn/error，默认 info
	ContextPath          string                   `json:"context_path" yaml:"context_path" ini:"context_path"`                                  // nacos 服务器上下文路径
	AppendToStdout       bool                     `json:"append_to_stdout" yaml:"append_to_stdout" ini:"append_to_stdout"`                      // 是否将日志追加到标准输出
	LogSampling          *ClientLogSamplingConfig `json:"log_sampling" yaml:"log_sampling" ini:"log_sampling"`                                  // 日志采样配置
	LogRollingConfig     *ClientLogRollingConfig  `json:"log_rolling_config" yaml:"log_rolling_config" ini:"log_rolling_config"`                // 日志滚动配置
	TLSCfg               TLSConfig                `json:"tls_cfg" yaml:"tls_cfg" ini:"tls_cfg"`                                                 // TLS 配置
	AsyncUpdateService   bool                     `json:"async_update_service" yaml:"async_update_service" ini:"async_update_service"`          // 是否通过查询开启异步更新服务
	EndpointContextPath  string                   `json:"endpoint_context_path" yaml:"endpoint_context_path" ini:"endpoint_context_path"`       // 地址服务器端点上下文路径
	EndpointQueryParams  string                   `json:"endpoint_query_params" yaml:"endpoint_query_params" ini:"endpoint_query_params"`       // 地址服务器端点查询参数
	ClusterName          string                   `json:"cluster_name" yaml:"cluster_name" ini:"cluster_name"`                                  // 地址服务器集群名称
	AppConnLabels        map[string]string        `json:"app_conn_labels" yaml:"app_conn_labels" ini:"app_conn_labels"`                         // 应用连接标签
}

type ClientLogSamplingConfig struct {
	Initial    int           `json:"initial" yaml:"initial" ini:"initial"`          // 日志采样的初始值
	Thereafter int           `json:"thereafter" yaml:"thereafter" ini:"thereafter"` // 日志采样的后续值
	Tick       time.Duration `json:"tick" yaml:"tick" ini:"tick"`                   // 日志采样的时间间隔
}

type ClientLogRollingConfig struct {
	MaxSize    int  `json:"max_size" yaml:"max_size" ini:"max_size"`          // 日志文件轮转前的最大大小（MB），默认 100MB
	MaxAge     int  `json:"max_age" yaml:"max_age" ini:"max_age"`             // 保留旧日志文件的最大天数（基于文件名中的时间戳）
	MaxBackups int  `json:"max_backups" yaml:"max_backups" ini:"max_backups"` // 保留的旧日志文件最大数量
	LocalTime  bool `json:"local_time" yaml:"local_time" ini:"local_time"`    // 是否使用本地时间而非 UTC 时间
	Compress   bool `json:"compress" yaml:"compress" ini:"compress"`          // 轮转的日志文件是否压缩
}

type TLSConfig struct {
	Appointed          bool   `json:"appointed" yaml:"appointed" ini:"appointed"`                                  // 是否指定，如果为 false 则从环境变量获取
	Enable             bool   `json:"enable" yaml:"enable" ini:"enable"`                                           // 是否启用 TLS
	TrustAll           bool   `json:"trust_all" yaml:"trust_all" ini:"trust_all"`                                  // 是否信任所有服务器
	CaFile             string `json:"ca_file" yaml:"ca_file" ini:"ca_file"`                                        // 客户端验证服务器证书时使用的 CA 文件
	CertFile           string `json:"cert_file" yaml:"cert_file" ini:"cert_file"`                                  // 服务器验证客户端证书时使用的证书文件
	KeyFile            string `json:"key_file" yaml:"key_file" ini:"key_file"`                                     // 服务器验证客户端证书时使用的密钥文件
	ServerNameOverride string `json:"server_name_override" yaml:"server_name_override" ini:"server_name_override"` // 仅用于测试的服务器名称覆盖
}

type KMSv3Config struct {
	ClientKeyContent string `json:"client_key_content" yaml:"client_key_content" ini:"client_key_content"` // 客户端密钥内容
	Password         string `json:"password" yaml:"password" ini:"password"`                               // 密码
	Endpoint         string `json:"endpoint" yaml:"endpoint" ini:"endpoint"`                               // 端点
	CaContent        string `json:"ca_content" yaml:"ca_content" ini:"ca_content"`                         // CA 内容
}

type KMSConfig struct {
	Endpoint  string `json:"endpoint" yaml:"endpoint" ini:"endpoint"`       // 端点
	OpenSSL   string `json:"openssl" yaml:"openssl" ini:"openssl"`          // OpenSSL 配置
	CaContent string `json:"ca_content" yaml:"ca_content" ini:"ca_content"` // CA 内容
}

type RamConfig struct {
	SecurityToken         string `json:"security_token" yaml:"security_token" ini:"security_token"`                            // 安全令牌
	SignatureRegionId     string `json:"signature_region_id" yaml:"signature_region_id" ini:"signature_region_id"`             // 签名区域 ID
	RamRoleName           string `json:"ram_role_name" yaml:"ram_role_name" ini:"ram_role_name"`                               // RAM 角色名称
	RoleArn               string `json:"role_arn" yaml:"role_arn" ini:"role_arn"`                                              // 角色 ARN
	Policy                string `json:"policy" yaml:"policy" ini:"policy"`                                                    // 策略
	RoleSessionName       string `json:"role_session_name" yaml:"role_session_name" ini:"role_session_name"`                   // 角色会话名称
	RoleSessionExpiration int    `json:"role_session_expiration" yaml:"role_session_expiration" ini:"role_session_expiration"` // 角色会话过期时间
	OIDCProviderArn       string `json:"oidc_provider_arn" yaml:"oidc_provider_arn" ini:"oidc_provider_arn"`                   // OIDC 提供者 ARN
	OIDCTokenFilePath     string `json:"oidc_token_file_path" yaml:"oidc_token_file_path" ini:"oidc_token_file_path"`          // OIDC 令牌文件路径
	CredentialsURI        string `json:"credentials_uri" yaml:"credentials_uri" ini:"credentials_uri"`                         // 凭证 URI
	SecretName            string `json:"secret_name" yaml:"secret_name" ini:"secret_name"`                                     // 密钥名称
}

func (this *Config) NewClient() (naming_client.INamingClient, error) {
	var serverConfigs []constant.ServerConfig
	for _, v := range this.ServerConfigs {
		serverConfigs = append(serverConfigs,
			constant.ServerConfig{Scheme: v.Scheme, ContextPath: v.ContextPath, IpAddr: v.IpAddr, Port: v.Port, GrpcPort: v.GrpcPort})
	}
	clientConfig := this.setClientConfig()
	this.setRamConfig(clientConfig.RamConfig)
	this.setKMSv3Config(clientConfig.KMSv3Config)
	this.setKMSConfig(clientConfig.KMSConfig)
	this.setLogSampling(clientConfig.LogSampling)
	this.setLogRollingConfig(clientConfig.LogRollingConfig)
	return clients.NewNamingClient(vo.NacosClientParam{
		ClientConfig:  clientConfig,
		ServerConfigs: serverConfigs,
	})
}

// 设置 client 配置
func (this *Config) setClientConfig() *constant.ClientConfig {
	return &constant.ClientConfig{
		TimeoutMs:            this.ClientConfig.TimeoutMs,
		BeatInterval:         this.ClientConfig.BeatInterval,
		NamespaceId:          this.ClientConfig.NamespaceId,
		AppName:              this.ClientConfig.AppName,
		AppKey:               this.ClientConfig.AppKey,
		Endpoint:             this.ClientConfig.Endpoint,
		RegionId:             this.ClientConfig.RegionId,
		AccessKey:            this.ClientConfig.AccessKey,
		SecretKey:            this.ClientConfig.SecretKey,
		OpenKMS:              this.ClientConfig.OpenKMS,
		KMSVersion:           this.ClientConfig.KMSVersion,
		CacheDir:             this.ClientConfig.CacheDir,
		DisableUseSnapShot:   this.ClientConfig.DisableUseSnapShot,
		UpdateThreadNum:      this.ClientConfig.UpdateThreadNum,
		NotLoadCacheAtStart:  this.ClientConfig.NotLoadCacheAtStart,
		UpdateCacheWhenEmpty: this.ClientConfig.UpdateCacheWhenEmpty,
		Username:             this.ClientConfig.Username,
		Password:             this.ClientConfig.Password,
		LogDir:               this.ClientConfig.LogDir,
		LogLevel:             this.ClientConfig.LogLevel,
		ContextPath:          this.ClientConfig.ContextPath,
		AppendToStdout:       this.ClientConfig.AppendToStdout,
		TLSCfg:               this.setTLSCfg(),
		AsyncUpdateService:   this.ClientConfig.AsyncUpdateService,
		EndpointContextPath:  this.ClientConfig.EndpointContextPath,
		EndpointQueryParams:  this.ClientConfig.EndpointQueryParams,
		ClusterName:          this.ClientConfig.ClusterName,
		AppConnLabels:        this.ClientConfig.AppConnLabels,
	}
}

// 设置 RAM 配置
func (this *Config) setRamConfig(v *constant.RamConfig) {
	if this.ClientConfig.RamConfig == nil {
		return
	}
	v.SecurityToken = this.ClientConfig.RamConfig.SecurityToken
	v.SignatureRegionId = this.ClientConfig.RamConfig.SignatureRegionId
	v.RamRoleName = this.ClientConfig.RamConfig.RamRoleName
	v.RoleArn = this.ClientConfig.RamConfig.RoleArn
	v.Policy = this.ClientConfig.RamConfig.Policy
	v.RoleSessionName = this.ClientConfig.RamConfig.RoleSessionName
	v.RoleSessionExpiration = this.ClientConfig.RamConfig.RoleSessionExpiration
	v.OIDCProviderArn = this.ClientConfig.RamConfig.OIDCProviderArn
	v.OIDCTokenFilePath = this.ClientConfig.RamConfig.OIDCTokenFilePath
	v.CredentialsURI = this.ClientConfig.RamConfig.CredentialsURI
	v.SecretName = this.ClientConfig.RamConfig.SecretName
}

// 设置 KMS v3 配置
func (this *Config) setKMSv3Config(v *constant.KMSv3Config) {
	if this.ClientConfig.KMSv3Config == nil {
		return
	}
	v.ClientKeyContent = this.ClientConfig.KMSv3Config.ClientKeyContent
	v.Password = this.ClientConfig.KMSv3Config.Password
	v.Endpoint = this.ClientConfig.KMSv3Config.Endpoint
	v.CaContent = this.ClientConfig.KMSv3Config.CaContent
}

// 设置 KMS 配置
func (this *Config) setKMSConfig(v *constant.KMSConfig) {
	if this.ClientConfig.KMSConfig == nil {
		return
	}
	v.Endpoint = this.ClientConfig.KMSConfig.Endpoint
	v.OpenSSL = this.ClientConfig.KMSConfig.OpenSSL
	v.CaContent = this.ClientConfig.KMSConfig.CaContent
}

func (this *Config) setLogSampling(v *constant.ClientLogSamplingConfig) {
	if this.ClientConfig.LogSampling == nil {
		return
	}
	v.Initial = this.ClientConfig.LogSampling.Initial
	v.Thereafter = this.ClientConfig.LogSampling.Thereafter
	v.Tick = tools.AutoTimeDuration(this.ClientConfig.LogSampling.Tick, time.Second)
}

func (this *Config) setLogRollingConfig(v *constant.ClientLogRollingConfig) {
	if this.ClientConfig.LogRollingConfig == nil {
		return
	}
	v.MaxSize = this.ClientConfig.LogRollingConfig.MaxSize
	v.MaxAge = this.ClientConfig.LogRollingConfig.MaxAge
	v.MaxBackups = this.ClientConfig.LogRollingConfig.MaxBackups
	v.LocalTime = this.ClientConfig.LogRollingConfig.LocalTime
	v.Compress = this.ClientConfig.LogRollingConfig.Compress
}

func (this *Config) setTLSCfg() constant.TLSConfig {
	return constant.TLSConfig{
		Appointed:          this.ClientConfig.TLSCfg.Appointed,
		Enable:             this.ClientConfig.TLSCfg.Enable,
		TrustAll:           this.ClientConfig.TLSCfg.TrustAll,
		CaFile:             this.ClientConfig.TLSCfg.CaFile,
		CertFile:           this.ClientConfig.TLSCfg.CertFile,
		KeyFile:            this.ClientConfig.TLSCfg.KeyFile,
		ServerNameOverride: this.ClientConfig.TLSCfg.ServerNameOverride,
	}
}
