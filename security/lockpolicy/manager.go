package lockpolicy

import (
	"fmt"
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

// RestoreLock 恢复锁定
// 用于程序重启后，从数据库载入所有锁定信息
// 仅仅恢复锁定目标，失败次数允许丢失。
func (m *Manager) RestoreLock(target LockTarget, identifier string, expire time.Time) {
	// 从当前时间计算剩余的锁定时间
	remaining := expire.Sub(time.Now())
	if remaining <= 0 {
		return
	}
	if cache, ok := m.cache.Load(target); ok {
		cache.SetLockWithExpire(identifier, remaining)
	}
}

// Clear 处理成功有，可以将失败缓存进行一个删除操作
func (m *Manager) Clear(targets Targets) {
	for target, identifier := range targets {
		if c, ok := m.cache.Load(target); ok {
			c.DeleteLock(identifier)         // 删除锁，主要是避免gc未触发，清理过期数据
			c.DeleteTriggerCount(identifier) // 删除触发次数
		}
	}

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
func (m *Manager) RecordFailure(target LockTarget, identifier string, callbacks ...LockCallback) (bool, *LockEvent) {
	return m.RecordFailures(map[LockTarget]string{target: identifier}, callbacks...)
}
func (m *Manager) RecordFailures(targets Targets, callbacks ...LockCallback) (bool, *LockEvent) {
	locked, event := m.recordIndependentPolicies(targets, callbacks...)
	if locked {
		return locked, event
	}

	return m.recordEscalationPolicies(targets, callbacks...)
}

// 记录独立策略
func (m *Manager) recordIndependentPolicies(targets Targets, callbacks ...LockCallback) (bool, *LockEvent) {
	var (
		isLocked bool
		event    *LockEvent
	)
	m.independentPolicies.ReadWith(func(policies Policies) {
		for _, policy := range policies {
			identifier, ok := targets[policy.Target]
			if !ok {
				continue
			}
			cache, ok := m.cache.Load(policy.Target)
			if !ok {
				continue
			}
			count := cache.SetTriggerCount(identifier)
			if count >= policy.Trigger {
				cache.SetLock(identifier)
				cache.DeleteTriggerCount(identifier) // 触发锁定后，需要重置连续错误次数
				isLocked = true
				event = m.recordEvent(identifier, &policy, LockTypeIndependent, callbacks...)
				break
			}
		}
	})

	return isLocked, event
}

func (m *Manager) recordEvent(identifier string, policy *Policy, lockType LockType, callbacks ...LockCallback) *LockEvent {
	event := &LockEvent{
		Target:        policy.Target,
		Identifier:    identifier,
		LockType:      lockType,
		LockoutTime:   policy.LockoutTime, // 锁定时长
		RemainingTime: policy.LockoutTime, // 剩余锁定时间
		Reason:        fmt.Sprintf("连续错误次数%d次，触发策略%s", policy.Trigger, policy.Target),
		Expire:        time.Now().Add(policy.LockoutTime), // 过期时间
		Timestamp:     time.Now(),                         // 锁定时间
		Policy:        *policy,
	}
	for _, callback := range callbacks {
		callback(*event)
	}
	return event
}

func (m *Manager) recordEscalationPolicies(targets Targets, callbacks ...LockCallback) (bool, *LockEvent) {
	var (
		isLocked bool
		event    *LockEvent
	)
	m.escalationChains.ReadWith(func(policies Policies) {
		pl := policies.Len()
		// 对于升级链，需要从后往回开始处理
		for i := pl - 1; i >= 0; i-- {
			policy := policies[i]
			identifier, ok := targets[policy.Target]
			if !ok {
				continue
			}
			cache, ok := m.cache.Load(policy.Target)
			if !ok {
				continue
			}
			// 如果当前缓存，触发次数是0 或者没有，也跳过厝里
			if cache.GetTriggerCount(identifier) <= 0 {
				continue
			}
			// 当前策略 如果有触发次数，还需要判断是否是升级链的第一个或者看上一个，是否启用记忆效应，如果启用了记忆效应，才直接再当前策略进行次数累计。
			if i > 0 {
				prevPolicy := policies[i-1]
				// 未启用记忆效应，则继续向前找
				if !prevPolicy.Escalation.MemoryEffect {
					continue
				}
			}
			// 升级链的第一个策略或者启用记忆效应，则进行次数累计
			count := cache.SetTriggerCount(identifier)
			// 触发次数达到阈值，则进行锁定。
			if count >= policy.Trigger {
				caches := []*cacheSlice{{cache, identifier}}
				// 判断 当前策略是否是升级链的最后一个、
				if i < (pl - 1) {
					caches = append(caches, m.updateLock(i+1, policies, targets)...)
				}
				// 直接用caches的最后一个策略进行锁定即可
				current := caches[len(caches)-1]
				current.cache.SetLock(current.identifier)
				cache.DeleteTriggerCount(current.identifier) // 触发锁定后，需要重置连续错误次数
				isLocked = true
				lkType := tools.Ternary(len(caches) == 1, LockTypeDirect, LockTypeEscalation)
				event = m.recordEvent(current.identifier, current.cache.policy, lkType, callbacks...)
				break
			}
		}
	})
	return isLocked, event
}

func (m *Manager) updateLock(idx int, chains Policies, targets Targets) []*cacheSlice {
	caches := make([]*cacheSlice, 0)
	for _, policy := range chains[idx:] {
		identifier, ok := targets[policy.Target]
		if !ok {
			continue
		}
		cache, ok := m.cache.Load(policy.Target)
		if !ok {
			continue
		}
		count := cache.SetTriggerCount(identifier)
		if count >= policy.Trigger {
			caches = append(caches, &cacheSlice{
				cache:      cache,
				identifier: identifier,
			})
		}
	}
	return caches
}
