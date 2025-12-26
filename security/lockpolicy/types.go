package lockpolicy

import (
	"sort"
	"time"

	"github.com/cespare/xxhash/v2"
)

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

// 为 LockTarget 专门创建的 Hasher
type lockTargetHasher struct{}

func (h lockTargetHasher) Hash(key LockTarget) uint64 {
	return xxhash.Sum64String(string(key))
}

// Policy 单条锁定策略
type Policy struct {
	Target      LockTarget    `json:"target" yaml:"target"`             // 锁定目标
	Trigger     int           `json:"trigger" yaml:"trigger"`           // 连续触发失败次数
	WindowTime  time.Duration `json:"window_time" yaml:"window_time"`   // 连续触发失败的窗口时间，多少时间内触发会累计缓存
	LockoutTime time.Duration `json:"lockout_time" yaml:"lockout_time"` // 连续失败Trigger后，目标的锁定时长
	Priority    int           `json:"priority" yaml:"priority"`         // 优先级 ，值越大，优先级越高，独立模式下使用
	// 升级配置
	// 未配置升级规则，切 连续触发失败次数 >0，作为独立策略
	Escalation *EscalationRule `json:"escalation,omitempty" yaml:"escalation"`
}

func (p *Policy) GetTarget() LockTarget {
	return p.Target
}

func (p *Policy) GetUpgradeTo() LockTarget {
	if p.Escalation != nil {
		return p.Escalation.UpgradeTo
	}
	return ""
}

// EscalationRule 升级规则
type EscalationRule struct {
	UpgradeTo LockTarget `json:"upgrade_to" yaml:"upgrade_to"` // 升级到哪个目标
	// 是否启用记忆效应
	// 比如升级链是 会话锁 IP锁，
	// 当会话锁定次数满足 IP锁的锁定条件后，触发IP锁，当下一次操作错误时，直接IP错误次数累加，而不是重新从会话锁开始升级上来。
	MemoryEffect bool `json:"memory_effect" yaml:"memory_effect"`
}

// Policies 策略集合
type Policies []Policy

func (p Policies) Len() int {
	return len(p)
}

// Less 策略排序 降序
func (p Policies) Less(i, j int) bool {
	return p[i].Priority > p[j].Priority
}

func (p Policies) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

// Sort 对策略进行排序（按优先级降序）- 使用指针接收者
func (p *Policies) Sort() {
	sort.Sort(p)
}

// Targets 锁定记录与检测的传参类型定义
type Targets map[LockTarget]string

type cacheSlice struct {
	cache      *targetCache
	identifier string
}
