package csrf

const (
	// Token相关常量

	DefaultTokenName  = "X-CSRF-Token"      // HTTP Header/Form/Query中的字段名
	DefaultCookieName = "csrf_token"        // HttpOnly Cookie名称
	DefaultStatusName = "csrf_token_status" // 状态Cookie名称
)
