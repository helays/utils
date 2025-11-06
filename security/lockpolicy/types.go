package lockpolicy

import "time"

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

// Policy 单条锁定策略
type Policy struct {
	Target      LockTarget    `json:"target" yaml:"target"`                   // 锁定目标
	Trigger     int           `json:"ip_trigger" yaml:"ip_trigger"`           // 连续触发失败次数
	WindowTime  time.Duration `json:"ip_window_time" yaml:"ip_window_time"`   // 连续触发失败的窗口时间，多少时间内触发会累计缓存
	LockoutTime time.Duration `json:"ip_lockout_time" yaml:"ip_lockout_time"` // 连续失败Trigger后，目标的锁定时长
	Priority    int           `json:"priority" yaml:"priority"`               // 优先级 ，值越大，优先级越高，独立模式下使用
	// 升级配置
	// 未配置升级规则，切 连续触发失败次数 >0，作为独立策略
	Escalation *EscalationRule `json:"escalation,omitempty" yaml:"escalation,omitempty"`
}

// EscalationRule 升级规则
type EscalationRule struct {
	// 基于锁定次数的升级
	LockoutCount int           `json:"lockout_count" yaml:"lockout_count"` // 自身锁定多少次后触发升级
	TimeWindow   time.Duration `json:"time_window" yaml:"time_window"`     // 锁定次数统计窗口

	// 升级目标
	UpgradeTo   LockTarget    `json:"upgrade_to" yaml:"upgrade_to"`     // 升级到哪个目标
	UpgradeTime time.Duration `json:"upgrade_time" yaml:"upgrade_time"` // 升级后的锁定时长
	// 记忆效应
	MemoryEffect bool `json:"memory_effect" yaml:"memory_effect"` // 是否启用记忆效应
}

// Policies 策略集合
type Policies []Policy
