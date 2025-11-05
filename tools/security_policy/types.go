package security_policy

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/helays/utils/v2/dataType"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

// LoginPolicy 登录安全
type LoginPolicy struct {
	MaxLoginAttempts    int  `json:"max_login_attempts"`    //最大登录尝试次数
	LockoutDuration     int  `json:"lockout_duration"`      // 锁定时长(分钟)
	SessionTimeout      int  `json:"session_timeout"`       // 会话超时(分钟)
	AllowMultiLocation  bool `json:"allow_multi_location"`  // 是否允许多地点登录
	MaxConcurrentLogins int  `json:"max_concurrent_logins"` // 最大并发登录数
}

// AccessPolicy 访问策略
type AccessPolicy struct {
	IPWhitelist       []string `json:"ip_whitelist"`        // IP白名单
	IPBlacklist       []string `json:"ip_blacklist"`        // IP黑名单
	AllowedRegions    []string `json:"allowed_regions"`     // 允许访问地区
	BusinessHoursOnly bool     `json:"business_hours_only"` // 仅限工作时间
}

// MFAPolicy 多因子策略
type MFAPolicy struct {
	Enabled         bool    `json:"enabled"`           // 是否启用MFA
	MFAType         MFAType `json:"mfa_type"`          // MFA类型
	Required        bool    `json:"required"`          // 是否强制要求MFA
	GracePeriod     int     `json:"grace_period"`      // 宽限期(天)
	BackupCodeCount int     `json:"backup_code_count"` // 备用验证码数量

	// 短信验证配置
	SMSEnabled        bool `json:"sms_enabled"`         // 启用短信验证
	SMSExpireMinutes  int  `json:"sms_expire_minutes"`  // 短信验证码过期时间
	SMSResendInterval int  `json:"sms_resend_interval"` // 重发间隔(秒)
	SMSDailyLimit     int  `json:"sms_daily_limit"`     // 每日发送限制

	// 邮箱验证配置
	EmailEnabled        bool `json:"email_enabled"`         // 启用邮箱验证
	EmailExpireMinutes  int  `json:"email_expire_minutes"`  // 邮箱验证过期时间
	EmailResendInterval int  `json:"email_resend_interval"` // 邮件重发间隔(秒)

	// 验证码配置
	Captcha CaptchaConfig `json:"captcha"` // 验证码

	// 生物识别
	BiometricEnabled bool     `json:"biometric_enabled"` // 启用生物识别
	BiometricTypes   []string `json:"biometric_types"`   // 支持的生物识别类型

	// 风险感知认证
	RiskBasedEnabled bool     `json:"risk_based_enabled"` // 启用风险感知认证
	HighRiskActions  []string `json:"high_risk_actions"`  // 高风险操作列表
}

// CaptchaConfig 验证码 安全策略配置
// 当 SessionTrigger > 0时，将启用连续n次验证码失败，锁定会话 SessionLockoutTime 时长
// 如果 SessionLockoutCount > 0,那么连续n次锁会话后，将锁定IP IPLockoutTime 时长。
type CaptchaConfig struct {
	CaptchaEnabled    bool        `json:"captcha_enabled" yaml:"captcha_enabled"` // 启用图形验证码
	CaptchaType       CaptchaType `json:"captcha_type" yaml:"captcha_type"`       // 验证码类型
	CaptchaLockPolicy LockPolicy  `json:"captcha_lock_policy" yaml:"captcha_lock_policy"`
}

// LockTarget 锁目标
type LockTarget string

func (t LockTarget) String() string {
	return string(t)
}

const (
	LockTargetSession LockTarget = "session" // 会话层锁定
	LockTargetIP      LockTarget = "ip"      // IP层锁定
	LockTargetUser    LockTarget = "user"    // 用户层锁定
)

// LockPolices 锁定时长策略
// 可以配置 会话层 连续失败n次，锁定IP
// IP连续锁定n次后，锁定用户
type LockPolices []LockTriggerPolicy

type LockTriggerPolicy struct {
	Target       LockTarget    `json:"target" yaml:"target"`                     // 锁定目标
	Trigger      int           `json:"ip_trigger" yaml:"ip_trigger"`             // 连续触发失败次数
	WindowTime   time.Duration `json:"ip_window_time" yaml:"ip_window_time"`     // 连续触发失败的窗口时间，多少时间内触发会累计缓存
	LockoutTime  time.Duration `json:"ip_lockout_time" yaml:"ip_lockout_time"`   // 连续失败Trigger后，目标的锁定时长
	LockoutCount int           `json:"ip_lockout_count" yaml:"ip_lockout_count"` //  锁定目标触发次数，用于升级锁定目标
}

// MFAType 多因子认证类型
type MFAType string

const (
	MFATypeNone     MFAType = "none"     // 无MFA
	MFATypeTOTP     MFAType = "totp"     // 时间型OTP
	MFATypeSMS      MFAType = "sms"      // 短信验证
	MFATypeEmail    MFAType = "email"    // 邮箱验证
	MFATypeBio      MFAType = "bio"      // 生物识别
	MFATypePush     MFAType = "push"     // 推送通知
	MFATypeWebAuthn MFAType = "webauthn" // Web认证
)

// CaptchaType 验证码类型
type CaptchaType string

const (
	CaptchaTypeImage  CaptchaType = "image"  // 图形验证码
	CaptchaTypeSlider CaptchaType = "slider" // 滑块验证
	CaptchaTypeClick  CaptchaType = "click"  // 点选验证
	CaptchaTypeSound  CaptchaType = "sound"  // 语音验证
)

type SecurityPolicyMeta struct {
	Password     PasswordPolicy `json:"password_policy"`
	LoginPolicy  LoginPolicy    `json:"login_policy"`
	AccessPolicy AccessPolicy   `json:"access_policy"`
	MFAPolicy    MFAPolicy      `json:"mfa_policy"`
}

func (t SecurityPolicyMeta) Value() (driver.Value, error) {
	return dataType.DriverValueWithJson(t)
}
func (t *SecurityPolicyMeta) Scan(val any) (err error) {
	return dataType.DriverScanWithJson(val, t)
}

func (t SecurityPolicyMeta) GormDataType() string {
	return "security_policy_meta"
}

func (SecurityPolicyMeta) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	return dataType.JsonDbDataType(db, field)
}

func (t SecurityPolicyMeta) GormValue(_ context.Context, db *gorm.DB) clause.Expr {
	byt, _ := json.Marshal(t)
	return dataType.MapGormValue(string(byt), db)
}
