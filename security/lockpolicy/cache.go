package lockpolicy

import (
	"time"

	"github.com/helays/utils/v2/map/safettl"
	"github.com/helays/utils/v2/tools/mutex"
)

// targetCache 锁定目标缓存
// 缓存中无需锁定次数，锁定次数即升级后的触发次数
type targetCache struct {
	policy       *Policy                                                         // 策略
	triggerCount *safettl.PerKeyTTLMap[string, *mutex.SafeResourceRWMutex[int]]  // 触发次数
	isLocked     *safettl.PerKeyTTLMap[string, *mutex.SafeResourceRWMutex[bool]] // 是否锁定

}

func newTargetCache(policy *Policy) *targetCache {
	return &targetCache{
		policy:       policy,
		triggerCount: safettl.NewPerKeyTTLMapWithInterval[string, *mutex.SafeResourceRWMutex[int]](time.Minute),
		isLocked:     safettl.NewPerKeyTTLMapWithInterval[string, *mutex.SafeResourceRWMutex[bool]](time.Minute),
	}
}

// GetTriggerCount 获取触发次数
func (t *targetCache) GetTriggerCount(key string) int {
	c, ok := t.triggerCount.Load(key)
	if !ok {
		return 0
	}
	return c.Read()
}

// SetTriggerCount 设置触发次数
func (t *targetCache) SetTriggerCount(key string) int {
	c, ok := t.triggerCount.Load(key)
	if !ok {
		c = mutex.NewSafeResourceRWMutex[int](0)
		t.triggerCount.StoreWithTTL(key, c, t.policy.WindowTime) // 触发次数缓存，需要有窗口时间
	}
	next := c.Read() + 1
	c.Write(next)
	t.triggerCount.Refresh(key, t.policy.WindowTime)
	return next
}

// IsLocked 获取是否锁定
func (t *targetCache) IsLocked(key string) (bool, time.Time) {
	c, expiry, ok := t.isLocked.LoadWithExpiry(key)
	if !ok {
		return false, expiry
	}

	return c.Read(), expiry
}

// DeleteLock 删除锁定
func (t *targetCache) DeleteLock(key string) {
	t.isLocked.Delete(key)
}

// DeleteTriggerCount 删除触发次数
func (t *targetCache) DeleteTriggerCount(key string) {
	t.triggerCount.Delete(key)
}

// SetLock 设置锁定
func (t *targetCache) SetLock(key string) {
	lock := mutex.NewSafeResourceRWMutex(true)
	t.isLocked.StoreWithTTL(key, lock, t.policy.LockoutTime)
}
