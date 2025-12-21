package lockpolicy

import (
	"context"
	"time"

	"github.com/helays/utils/v2/safe"
)

// targetCache 锁定目标缓存
// 缓存中无需锁定次数，锁定次数即升级后的触发次数
type targetCache struct {
	policy       *Policy                                        // 策略
	triggerCount *safe.Map[string, *safe.ResourceRWMutex[int]]  // 触发次数
	isLocked     *safe.Map[string, *safe.ResourceRWMutex[bool]] // 是否锁定

}

func newTargetCache(ctx context.Context, policy *Policy) *targetCache {
	return &targetCache{
		policy: policy,
		triggerCount: safe.NewMap[string, *safe.ResourceRWMutex[int]](ctx, safe.StringHasher{}, safe.MapConfig{
			EnableCleanup: true,
			ClearInterval: time.Minute / 2,
			TTL:           time.Minute,
		}),
		isLocked: safe.NewMap[string, *safe.ResourceRWMutex[bool]](ctx, safe.StringHasher{}, safe.MapConfig{
			EnableCleanup: true,
			ClearInterval: time.Minute / 2,
			TTL:           time.Minute,
		}),
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
		c = safe.NewResourceRWMutex[int](0)
		t.triggerCount.Store(key, c, t.policy.WindowTime) // 触发次数缓存，需要有窗口时间
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
	lock := safe.NewResourceRWMutex(true)
	t.isLocked.Store(key, lock, t.policy.LockoutTime)
}

func (t *targetCache) SetLockWithExpire(key string, expire time.Duration) {
	lock := safe.NewResourceRWMutex(true)
	t.isLocked.Store(key, lock, expire)
}
