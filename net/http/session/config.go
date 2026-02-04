package session

import (
	"time"

	"helay.net/go/utils/v3/net/http/cookiekit"
)

// Options session 配置
type Options struct {
	DisableGc     bool             `json:"disable_gc" yaml:"disable_gc" ini:"disable_gc"`             // 是否禁用session gc
	GcProbability float64          `json:"gc_probability" yaml:"gc_probability" ini:"gc_probability"` // session gc 概率
	CheckInterval time.Duration    `json:"check_interval" yaml:"check_interval" ini:"check_interval"` // session 检测默认有效期
	Carrier       CookieCarrier    `json:"carrier" yaml:"carrier" ini:"carrier"`                      // session 载体，默认cookie
	Cookie        cookiekit.Config `json:"cookie" yaml:"cookie" ini:"cookie"`
}
