package ipmatch

import (
	"sort"

	"helay.net/go/utils/v3/tools"
)

func (m *IPMatcher) optimizeIPv4Storage(ranges []ipv4Range) {
	totalIPs := uint64(0)
	singleIPs := uint64(0)
	ipRanges := uint64(0)

	// 统计规则内的IP数量
	for _, r := range ranges {
		rangeSize := r.End - r.Start + 1
		totalIPs += uint64(rangeSize)
		if r.Start == r.End {
			singleIPs++
		} else {
			ipRanges++
		}
	}
	threshold := tools.Ternary(m.config.IPv4MapThreshold > 0, m.config.IPv4MapThreshold, DefaultIPv4MapThreshold)
	if totalIPs <= threshold {
		m.expandAllIPv4ToMap(ranges)
	} else {
		m.optimizeIPv4MixedStrategy(ranges, threshold)
	}
}

// expandAllIPv4ToMap 将所有IPv4范围展开为离散IP
func (m *IPMatcher) expandAllIPv4ToMap(ranges []ipv4Range) {
	m.storage.ipv4Map = make(map[uint32]struct{})
	for _, r := range ranges {
		for i := r.Start; i <= r.End; i++ {
			m.storage.ipv4Map[i] = struct{}{}
		}
	}
	// 清空范围存储，只使用map
	m.storage.ipv4Ranges = nil
}

// optimizeIPv4MixedStrategy 混合策略优化
func (m *IPMatcher) optimizeIPv4MixedStrategy(ranges []ipv4Range, threshold uint64) {
	// 重新初始化map（用于存储部分IP）
	m.storage.ipv4Map = make(map[uint32]struct{})

	usedCapacity := uint64(0)
	ipRanges := make([]ipv4Range, 0)
	// 分离单个IP和范围IP
	for _, r := range ranges {
		if r.Start == r.End {
			m.storage.ipv4Map[r.Start] = struct{}{}
			usedCapacity++
		} else {
			ipRanges = append(ipRanges, r)
		}
	}
	// 单IP已经超过阈值，range类型就不拆分展开了
	if usedCapacity > threshold {
		m.storage.ipv4Ranges = ipRanges
		return
	}
	// 对ipRanges 按照范围大小进行排序，小范围的优先
	sort.Slice(ipRanges, func(i, j int) bool {
		sizeI := ipRanges[i].End - ipRanges[i].Start + 1
		sizeJ := ipRanges[j].End - ipRanges[j].Start + 1
		return sizeI < sizeJ
	})
	// 步骤5：根据剩余容量依次展开范围
	remainingRanges := make([]ipv4Range, 0)
	remainingCapacity := threshold - usedCapacity
	ipRangesLen := len(ipRanges)
	for idx, r := range ipRanges {
		rangeSize := uint64(r.End - r.Start + 1)
		if rangeSize <= remainingCapacity {
			// 可以完整展开整个范围
			for i := r.Start; i <= r.End; i++ {
				m.storage.ipv4Map[i] = struct{}{}
			}
			usedCapacity += rangeSize
			remainingCapacity -= rangeSize
		} else if remainingCapacity > 0 {
			// 部分展开：只展开剩余容量能容纳的部分
			end := r.Start + uint32(remainingCapacity) - 1
			// 展开部分IP
			for i := r.Start; i <= end; i++ {
				m.storage.ipv4Map[i] = struct{}{}
			}
			usedCapacity = threshold // map已满

			// 剩余部分保留为范围
			if end < r.End {
				remainingRanges = append(remainingRanges, ipv4Range{Start: end + 1, End: r.End})
			}
			// 添加剩余未处理的范围
			if idx+1 < ipRangesLen {
				remainingRanges = append(remainingRanges, ipRanges[idx+1:]...)
			}
			break
		} else {
			// 没有剩余容量，整个范围保留
			remainingRanges = append(remainingRanges, r)
		}
	}

	// 步骤6：对剩余范围重新排序和合并
	if len(remainingRanges) == 0 {
		return
	}
	// 按起始IP排序（为了二分查找）
	sort.Slice(remainingRanges, func(i, j int) bool {
		return remainingRanges[i].Start < remainingRanges[j].Start
	})

	m.storage.ipv4Ranges = remainingRanges
}
