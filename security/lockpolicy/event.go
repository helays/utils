package lockpolicy

import "time"

// LockType 锁定类型
type LockType string

const (
	LockTypeIndependent LockType = "independent" // 独立锁触发锁定
	LockTypeDirect      LockType = "direct"      // 直接触发锁定
	LockTypeEscalation  LockType = "escalation"  // 升级触发锁定
)

// LockEvent 锁定事件
type LockEvent struct {
	Target        LockTarget    `json:"target"`
	Identifier    string        `json:"identifier"`
	LockType      LockType      `json:"lock_type"`      // 锁定类型：direct-直接触发, escalation-升级触发, memory-记忆效应
	LockoutTime   time.Duration `json:"lockout_time"`   // 锁定时长
	RemainingTime time.Duration `json:"remaining_time"` // 剩余锁定时间
	Reason        string        `json:"reason"`         // 锁定原因
	Timestamp     time.Time     `json:"timestamp"`      // 锁定时间
	Expire        time.Time     `json:"expire"`         // 过期时间
}

// LockCallback 锁定回调函数
type LockCallback func(event LockEvent)

func (e LockEvent) Error() string {
	return e.Reason
}
