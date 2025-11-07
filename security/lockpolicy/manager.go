package lockpolicy

import (
	"github.com/helays/utils/v2/map/safemap"
	"github.com/helays/utils/v2/tools"
	"github.com/helays/utils/v2/tools/mutex"
)

type Manager struct {
	policyMap           safemap.SyncMap[LockTarget, *Policy]
	policies            *mutex.SafeResourceRWMutex[Policies] // 锁定策略配置
	independentPolicies *mutex.SafeResourceRWMutex[Policies] // 独立策略映射
	escalationChains    *mutex.SafeResourceRWMutex[Policies] // 升级链映射
	cache               safemap.SyncMap[LockTarget, *targetCache]
}

// NewManager 创建策略管理器
func NewManager(policies Policies) *Manager {
	m := &Manager{}
	m.policies = mutex.NewSafeResourceRWMutex(policies)
	m.independentPolicies = mutex.NewSafeResourceRWMutex(Policies{})
	m.escalationChains = mutex.NewSafeResourceRWMutex(Policies{})
	m.buildPolicy()
	for _, policy := range policies {
		m.cache.Store(policy.Target, newTargetCache())
	}
	return m
}

// UpdatePolices 更新策略
func (m *Manager) UpdatePolices(policies Policies) {
	m.policies.Write(policies)
	m.buildPolicy()
	// 判断是否有策略删除
	m.cache.Range(func(key LockTarget, value *targetCache) bool {
		if tools.ContainsByField(policies, key, func(policy Policy) LockTarget { return policy.Target }) {
			m.cache.Delete(key) // 删除缓存
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

// 构建策略
func (m *Manager) buildPolicy() {
	m.policyMap.DeleteAll()
	m.buildIndependent()
	m.buildEscalationChain()
}

// 构建独立策略
func (m *Manager) buildIndependent() {
	var polices Policies
	for _, policy := range m.policies.Read() {
		m.policyMap.Store(policy.Target, &policy) // 策略映射
		policy.Valid()
		if policy.Trigger > 0 && policy.Escalation != nil {
			polices = append(polices, policy)
		}
	}
	m.independentPolicies.Write(polices)
}

// 构建升级策略链
func (m *Manager) buildEscalationChain() {

}

func (m *Manager) RecordFailure(target LockTarget, identifier string, callbacks ...LockCallback) error {
	return m.RecordFailures(map[LockTarget]string{target: identifier}, callbacks...)
}
func (m *Manager) RecordFailures(targets Targets, callbacks ...LockCallback) error {
	// 待实现
	panic("implement me")
}
