package csrf

import "time"

// PerRequestConfig 每次请求前都需要获取一次token 的预设配置
func PerRequestConfig() *Config {
	return &Config{
		Enabled:       true,
		Strategy:      StrategyToken,
		Timeout:       5 * time.Minute,
		TokenMode:     TokenModePerRequest,
		Secure:        true,
		ExemptMethods: []string{"GET", "HEAD", "OPTIONS"},
	}
}

// SessionTokenConfig 使用session存储token的预设配置
func SessionTokenConfig() *Config {
	return &Config{
		Enabled:       true,
		Strategy:      StrategyToken,
		Timeout:       24 * time.Hour,
		TokenMode:     TokenModeSession,
		Secure:        true,
		ExemptMethods: []string{"GET", "HEAD", "OPTIONS"},
	}
}

// TimedTokenConfig 使用定时存储token的预设配置
func TimedTokenConfig(timeout time.Duration) *Config {
	return &Config{
		Enabled:       true,
		Strategy:      StrategyToken,
		Timeout:       timeout,
		TokenMode:     TokenModeTimed,
		Secure:        true,
		ExemptMethods: []string{"GET", "HEAD", "OPTIONS"},
	}
}

// PermissiveConfig 宽松配置（用于内部API）
func PermissiveConfig() *Config {
	return &Config{
		Enabled:       false,
		Strategy:      StrategyNone,
		ExemptMethods: []string{},
	}
}
