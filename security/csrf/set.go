package csrf

import "time"

func (c *Config) Enable() *Config {
	c.Enabled = true
	return c
}

func (c *Config) Disable() *Config {
	c.Enabled = false
	return c
}

func (c *Config) WithStrategy(strategy Strategy) *Config {
	c.Strategy = strategy
	return c
}

func (c *Config) WithTokenSource(tokenSource TokenSource) *Config {
	c.TokenSource = tokenSource
	return c
}

func (c *Config) WithTokenName(tokenName string) *Config {
	c.TokenName = tokenName
	return c
}

func (c *Config) WithTimeout(seconds time.Duration) *Config {
	c.Timeout = seconds
	return c
}

func (c *Config) WithTokenMode(tokenMode TokenMode) *Config {
	c.TokenMode = tokenMode
	return c
}

func (c *Config) WithSameSite(sameSite string) *Config {
	c.SameSite = sameSite
	return c
}

func (c *Config) WithSecure(secure bool) *Config {
	c.Secure = secure
	return c
}

func (c *Config) WithExemptMethods(methods ...string) *Config {
	c.ExemptMethods = methods
	return c
}

func (c *Config) Clone() *Config {
	if c == nil {
		return nil
	}
	clone := *c // 浅拷贝基础字段

	// 深拷贝切片
	if c.ExemptMethods != nil {
		clone.ExemptMethods = make([]string, len(c.ExemptMethods))
		copy(clone.ExemptMethods, c.ExemptMethods)
	}

	return &clone
}
