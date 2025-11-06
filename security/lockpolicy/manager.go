package lockpolicy

import (
	"github.com/helays/utils/v2/map/safemap"
	"github.com/helays/utils/v2/tools/mutex"
)

type Manager struct {
	policies *mutex.SafeResourceRWMutex[Policies] // 锁定策略配置
	cache    safemap.SyncMap[LockTarget, *targetCache]
}

// NewManager 创建策略管理器
func NewManager(policies Policies) *Manager {
	m := &Manager{}
	m.policies = mutex.NewSafeResourceRWMutex(policies)
	for _, policy := range policies {
		m.cache.Store(policy.Target, newTargetCache())
	}
	return m
}

// UpdatePolices 更新策略
func (m *Manager) UpdatePolices(policies Policies) {
	m.policies.Write(policies)
	// 判断是否有策略删除
	m.cache.Range(func(key LockTarget, value *targetCache) bool {
		deleted := true
		for _, policy := range policies {
			if policy.Target == key {
				deleted = false
			}
		}
		if deleted {
			// 删除缓存
			m.cache.Delete(key)
		}
		return true
	})
	for _, policy := range policies {
		// 添加缓存
		if _, ok := m.cache.Load(policy.Target); !ok {
			m.cache.Store(policy.Target, newTargetCache())
		}
	}
}
func (m *Manager) RecordFailure(target LockTarget, identifier string, callbacks ...LockCallback) error {
	return m.RecordFailures(map[LockTarget]string{target: identifier}, callbacks...)
}
func (m *Manager) RecordFailures(targets map[LockTarget]string, callbacks ...LockCallback) error {
	// 待实现
	panic("implement me")
}
