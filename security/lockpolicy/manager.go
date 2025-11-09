package lockpolicy

import (
	"time"

	"github.com/helays/utils/v2/map/safemap"
	"github.com/helays/utils/v2/tools"
	"github.com/helays/utils/v2/tools/mutex"
)

type Manager struct {
	policyMap           safemap.SyncMap[LockTarget, *Policy] // 构建按 锁定目标 => 策略映射
	policies            *mutex.SafeResourceRWMutex[Policies] // 锁定策略配置
	independentPolicies *mutex.SafeResourceRWMutex[Policies] // 独立策略映射
	escalationChains    *mutex.SafeResourceRWMutex[Policies] // 升级链映射
	cache               safemap.SyncMap[LockTarget, *targetCache]
}

// NewManager 创建策略管理器
func NewManager(polices Policies) *Manager {
	m := &Manager{}
	m.policies = mutex.NewSafeResourceRWMutex(polices)
	m.independentPolicies = mutex.NewSafeResourceRWMutex(Policies{})
	m.escalationChains = mutex.NewSafeResourceRWMutex(Policies{})
	m.buildPolicy()
	m.setTargetCache(polices)
	return m
}

// UpdatePolices 更新策略
func (m *Manager) UpdatePolices(polices Policies) {
	m.policies.Write(polices)
	m.buildPolicy()
	m.setTargetCache(polices)
}

// 为锁定策略添加缓存实现
func (m *Manager) setTargetCache(polices Policies) {
	m.cache.DeleteAll()
	for _, policy := range polices {
		if policy.Trigger > 0 {
			// 添加缓存
			m.cache.Store(policy.Target, newTargetCache(&policy))
		}
	}
}

// 构建策略
func (m *Manager) buildPolicy() {
	m.policyMap.DeleteAll()

	var (
		filterPolicies     = make(Policies, 0)
		independentPolices = make(Policies, 0)
		escalationPolices  = make(Policies, 0)
		policies           = m.policies.Read()
	)
	for _, policy := range policies {
		if policy.Trigger < 1 {
			continue
		}
		policy.Valid()
		m.policyMap.Store(policy.Target, &policy) // 策略映射
		filterPolicies = append(filterPolicies, policy)
	}
	for _, policy := range filterPolicies {
		// 当前策略无升级目标，但还需要检测是否是其他策略的升级目标
		if policy.Escalation == nil && !tools.ContainsByField(filterPolicies, policy.GetTarget(), func(policy Policy) LockTarget { return policy.GetUpgradeTo() }) {
			independentPolices = append(independentPolices, policy)
		} else {
			escalationPolices = append(escalationPolices, policy)
		}
	}

	// 独立升级策略，直接根据优先级排序，降序
	independentPolices.Sort()
	m.independentPolicies.Write(independentPolices)

	// 升级策略，需要根据触发的顺序升序
	m.buildEscalationChain(escalationPolices)
}

// 构建升级策略链
func (m *Manager) buildEscalationChain(escalationPolices Policies) {
	// 找到升级链起点
	var startPolicies = make(Policies, 0)
	for _, policy := range escalationPolices {
		if policy.Escalation == nil {
			continue
		}
		// 如果当前策略目标，在其他策略的升级目标中无法查询到，就可以作为升级链起点。
		if !tools.ContainsByField(escalationPolices, policy.GetTarget(), func(policy Policy) LockTarget { return policy.GetUpgradeTo() }) {
			startPolicies = append(startPolicies, policy)
		}
	}
	if len(startPolicies) == 0 {
		m.escalationChains.Write(Policies{})
		return
	}

	var (
		chains        Policies                // 从起点开始构建升级链
		visited       = map[LockTarget]bool{} // 访问过的策略,避免循环引用
		currentPolicy = startPolicies[0]      // 取第一个起点开始构建升级链
	)

	for {
		// 检查是否已访问过，避免循环引用
		if visited[currentPolicy.Target] {
			break
		}
		// 添加到链中并标记已访问
		chains = append(chains, currentPolicy)
		visited[currentPolicy.Target] = true
		// 检查是否有升级目标
		if currentPolicy.Escalation == nil {
			break
		}
		nextTarget := currentPolicy.Escalation.UpgradeTo
		nextPolicy, ok := m.policyMap.Load(nextTarget)
		if !ok {
			break
		}
		currentPolicy = *nextPolicy
	}
	m.escalationChains.Write(chains)
}

// IsLocked 检查目标是否被锁定
func (m *Manager) IsLocked(targets Targets) (bool, *LockEvent) {
	for target, identifier := range targets {
		if c, ok := m.cache.Load(target); ok {
			isLocked, expire := c.IsLocked(identifier)
			if isLocked {
				event := &LockEvent{
					Target:     target,
					Identifier: identifier,
				}
				if policy, exist := m.policyMap.Load(target); exist {
					event.LockoutTime = policy.LockoutTime
					event.RemainingTime = expire.Sub(time.Now())
					event.Expire = expire
				}
				return true, event
			}
		}
	}
	return false, nil
}
func (m *Manager) RecordFailure(target LockTarget, identifier string, callbacks ...LockCallback) error {
	return m.RecordFailures(map[LockTarget]string{target: identifier}, callbacks...)
}
func (m *Manager) RecordFailures(targets Targets, callbacks ...LockCallback) error {
	// 待实现
	panic("implement me")
}
