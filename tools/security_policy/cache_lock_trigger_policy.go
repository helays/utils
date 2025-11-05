package security_policy

import (
	"github.com/helays/utils/v2/map/safettl"
	"github.com/helays/utils/v2/map/syncMapWrapper"
	"github.com/helays/utils/v2/tools/mutex"
)

func NewTriggerCache(polices LockPolices) *TriggerCache {
	c := &TriggerCache{policies: polices}
	for _, policy := range polices {
		c.cache.Store(policy.Target, newTargetCache())
	}
	return c
}

type TriggerCache struct {
	policies LockPolices
	cache    syncMapWrapper.SyncMap[LockTarget, *targetCache]
}

type targetCache struct {
	triggerCount *safettl.PerKeyTTLMap[string, *mutex.SafeResourceRWMutex[int]] // 触发次数
	lockCount    *safettl.PerKeyTTLMap[string, *mutex.SafeResourceRWMutex[int]] // 锁定次数
}

func newTargetCache() *targetCache {
	return &targetCache{
		triggerCount: safettl.NewPerKeyTTLMapWithInterval[string, *mutex.SafeResourceRWMutex[int]](time.Minute), // gc频率，1分钟够用了。
		lockCount:    safettl.NewPerKeyTTLMapWithInterval[string, *mutex.SafeResourceRWMutex[int]](time.Minute),
	}
}
