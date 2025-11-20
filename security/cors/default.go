package cors

// StrictConfig 返回严格的CORS配置
func StrictConfig() Config {
	return Config{
		Enabled:          true,
		AllowOrigins:     []string{},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		ExposeHeaders:    []string{},
		MaxAge:           3600,
	}
}
