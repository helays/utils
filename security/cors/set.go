package cors

// WithOrigins 设置允许的源
func (c *Config) WithOrigins(origins ...string) *Config {
	c.AllowOrigins = origins
	return c
}

// WithMethods 设置允许的方法
func (c *Config) WithMethods(methods ...string) *Config {
	c.AllowMethods = methods
	return c
}

// WithHeaders 设置允许的请求头
func (c *Config) WithHeaders(headers ...string) *Config {
	c.AllowHeaders = headers
	return c
}

// WithCredentials 设置是否允许凭据
func (c *Config) WithCredentials(allow bool) *Config {
	c.AllowCredentials = allow
	return c
}

// WithExposeHeaders 设置暴露的响应头
func (c *Config) WithExposeHeaders(headers ...string) *Config {
	c.ExposeHeaders = headers
	return c
}

// WithMaxAge 设置预检请求缓存时间
func (c *Config) WithMaxAge(seconds int) *Config {
	c.MaxAge = seconds
	return c
}

// Enable 启用CORS
func (c *Config) Enable() *Config {
	c.Enabled = true
	return c
}

// Disable 禁用CORS
func (c *Config) Disable() *Config {
	c.Enabled = false
	return c
}

// Clone 创建配置的深拷贝
func (c *Config) Clone() *Config {
	clone := *c

	// 切片需要深拷贝
	if c.AllowOrigins != nil {
		clone.AllowOrigins = make([]string, len(c.AllowOrigins))
		copy(clone.AllowOrigins, c.AllowOrigins)
	}

	if c.AllowMethods != nil {
		clone.AllowMethods = make([]string, len(c.AllowMethods))
		copy(clone.AllowMethods, c.AllowMethods)
	}

	if c.AllowHeaders != nil {
		clone.AllowHeaders = make([]string, len(c.AllowHeaders))
		copy(clone.AllowHeaders, c.AllowHeaders)
	}

	if c.ExposeHeaders != nil {
		clone.ExposeHeaders = make([]string, len(c.ExposeHeaders))
		copy(clone.ExposeHeaders, c.ExposeHeaders)
	}

	return &clone
}
