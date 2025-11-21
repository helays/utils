package ipmatch

func (m *IPMatcher) Contains(ip string) bool {
	if !m.enable {
		return true
	}
	// TODO: 待实现
	panic("TODO：待实现")
}
