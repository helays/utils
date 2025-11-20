package csrf

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"net/http"
	"time"

	"github.com/helays/utils/v2/crypto/xxhashkit"
	"github.com/helays/utils/v2/dataType"
	"github.com/helays/utils/v2/tools"
	"github.com/vmihailenco/msgpack/v5"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type Strategy string

func (s Strategy) String() string {
	return string(s)
}

const (
	StrategyNone  Strategy = "none"  // 不需要CSRF保护
	StrategyToken Strategy = "token" // 使用Token验证
	// StrategyDoubleTap 双重提交Cookie（Double Submit Cookie） 是一种CSRF防护技术，它的核心思想是：
	// 服务端生成一个随机Token,通过Cookie发送给浏览器
	// 前端在请求中携带这个token，通常放在header或者query中，
	// 服务端验证请求中的Token和Cookie中的Token是否一致
	// /api/csrf-token接口在生成token时，需要根据对应path的配置来处理，另外需要设置两个cookie,一个可读，一个用于验证；可读的不设置具体值
	StrategyDoubleTap Strategy = "double_tap" // 双重提交Cookie
)

// TokenMode 令牌模式
// 支持三种模式
// 每次请求前都需要获取一次 （需要和path绑定）
// 会话级别
// 带有效期 （需要绑定到path）
type TokenMode string

const (
	TokenModePerRequest TokenMode = "per_request" // 每次请求前获取（最安全）
	TokenModeSession    TokenMode = "session"     // 会话级Token（登录后全局使用）
	TokenModeTimed      TokenMode = "timed"       // 带有效期Token
)

type Config struct {
	Enabled       bool          `json:"enabled,omitempty" yaml:"enabled" ini:"enabled"`                      // 是否启用CSRF保护
	Strategy      Strategy      `json:"strategy,omitempty" yaml:"strategy" ini:"strategy"`                   // 防护策略
	Timeout       time.Duration `json:"timeout,omitempty" yaml:"timeout" ini:"timeout"`                      // Token超时时间(秒)
	TokenMode     TokenMode     `json:"token_mode,omitempty" yaml:"token_mode" ini:"token_mode"`             // Token模式
	SameSite      http.SameSite `json:"same_site,omitempty" yaml:"same_site" ini:"same_site"`                // SameSite策略: strict/lax/none
	Secure        bool          `json:"secure,omitempty" yaml:"secure" ini:"secure"`                         // 是否仅HTTPS
	ExemptMethods []string      `json:"exempt_methods,omitempty" yaml:"exempt_methods" ini:"exempt_methods"` // 豁免的HTTP方法
}

func (c Config) Value() (driver.Value, error) {
	return dataType.DriverValueWithJson(c)
}
func (c *Config) Scan(val any) (err error) {
	return dataType.DriverScanWithJson(val, c)
}

func (c Config) GormDataType() string {
	return "json"
}

func (Config) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	return dataType.JsonDbDataType(db, field)
}

func (c Config) GormValue(_ context.Context, db *gorm.DB) clause.Expr {
	byt, _ := json.Marshal(c)
	return dataType.MapGormValue(string(byt), db)
}

func (c Config) GobEncode() ([]byte, error) {
	return msgpack.Marshal(c)
}

func (c *Config) GobDecode(data []byte) error {
	return msgpack.Unmarshal(data, c)
}

// IsValid 验证配置是否有效
func (c *Config) IsValid() bool {
	if c == nil {
		return false
	}
	if c.Enabled {
		if c.Strategy == "" {
			return false
		}
	}
	return true
}

// ShouldValidate 验证是否需要进行CSRF验证
func (c *Config) ShouldValidate(method ...string) bool {
	if c == nil || !c.Enabled {
		return false
	}
	// 不需要进行 http method验证
	if len(method) < 1 {
		return true
	}
	// 判断是否在豁免列表中
	return !tools.Contains(c.ExemptMethods, method[0])
}

// GetTokenBinding 这里需要用xxhash对 path进行hash编码。
func (c *Config) GetTokenBinding(path string) string {
	prefix := "csrf_"
	switch c.TokenMode {
	case TokenModePerRequest, TokenModeTimed:
		return prefix + xxhashkit.XXHashString(path)
	case TokenModeSession:
		return prefix + "global"
	default:
		return prefix + "global"
	}
}
