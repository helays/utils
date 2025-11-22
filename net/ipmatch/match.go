package ipmatch

import (
	"net/netip"

	"github.com/helays/utils/v2/net/ipkit"
)

func (m *IPMatcher) Contains(ip string) bool {
	if !m.enable {
		return true
	}
	if m.config.Dynamic {
		m.mu.RLock()
		defer m.mu.RUnlock()
	}
	// 使用 netip 解析IP地址
	addr, err := netip.ParseAddr(ip)
	if err != nil {
		return false
	}
	if addr.Is4() {
		return m.containsIPv4(addr)
	} else if addr.Is6() {
		return m.containsIPv6(addr)
	} else {
		return false
	}
}

// containsIPv4 判断IPv4地址是否在IPv4匹配表中
func (m *IPMatcher) containsIPv4(addr netip.Addr) bool {
	ipUint := ipkit.IPAddr2Int(addr)

	if _, ok := m.storage.ipv4Map[ipUint]; ok {
		return true
	}
	if len(m.storage.ipv4Ranges) == 0 {
		return false
	}
	// 缓存加速查询，如果存在还会刷新缓存
	if _, ok := m.ipv4Cache.LoadAndRefresh(ipUint); ok {
		return true
	}
	left, right := 0, len(m.storage.ipv4Ranges)-1
	for left <= right {
		mid := (left + right) / 2
		r := m.storage.ipv4Ranges[mid]

		if ipUint < r.Start {
			right = mid - 1
		} else if ipUint > r.End {
			left = mid + 1
		} else {
			m.ipv4Cache.Store(ipUint, struct{}{})
			return true
		}
	}

	return false
}

func (m *IPMatcher) containsIPv6(addr netip.Addr) bool {
	k := addr.As16()
	// 检查map
	if _, exists := m.storage.ipv6Map[k]; exists {
		return true
	}

	// 在范围内二分查找
	if len(m.storage.ipv6Ranges) == 0 {
		return false
	}

	if _, ok := m.ipv6Cache.LoadAndRefresh(k); ok {
		return true
	}

	left, right := 0, len(m.storage.ipv6Ranges)-1
	for left <= right {
		mid := (left + right) / 2
		r := m.storage.ipv6Ranges[mid]

		if addr.Less(r.start) {
			right = mid - 1
		} else if r.end.Less(addr) {
			left = mid + 1
		} else {
			m.ipv6Cache.Store(k, struct{}{})
			return true
		}
	}

	return false
}
