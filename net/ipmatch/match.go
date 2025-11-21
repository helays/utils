package ipmatch

func (m *IPMatcher) Contains(ip string) bool {
	if !m.enable {
		return true
	}
	if m.config.Dynamic {
		m.mu.RLock()
		defer m.mu.RUnlock()
	}
	// TODO: 待实现
	panic("TODO：待实现")
}
