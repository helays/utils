package lockpolicy

import (
	"time"

	"github.com/helays/utils/v2/map/safettl"
	"github.com/helays/utils/v2/tools/mutex"
)

type targetCache struct {
	triggerCount *safettl.PerKeyTTLMap[string, *mutex.SafeResourceRWMutex[int]] // 触发次数
	lockCount    *safettl.PerKeyTTLMap[string, *mutex.SafeResourceRWMutex[int]] // 锁定次数
}

func newTargetCache() *targetCache {
	return &targetCache{
		triggerCount: safettl.NewPerKeyTTLMapWithInterval[string, *mutex.SafeResourceRWMutex[int]](time.Minute),
		lockCount:    safettl.NewPerKeyTTLMapWithInterval[string, *mutex.SafeResourceRWMutex[int]](time.Minute),
	}
}
