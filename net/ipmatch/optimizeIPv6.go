package ipmatch

import (
	"math/big"
	"net/netip"
	"sort"

	"github.com/helays/utils/v2/tools"
)

// optimizeIPv6Storage 优化IPv6存储
func (m *IPMatcher) optimizeIPv6Storage(ranges []ipv6Range) {
	totalIPs := &big.Int{}
	singleIPs := &big.Int{}
	ipRanges := &big.Int{}

	for _, r := range ranges {
		rangeSize := calculateIPv6RangeSize(r.start, r.end)
		totalIPs.Add(totalIPs, rangeSize)

		if r.start.Compare(r.end) == 0 {
			singleIPs.Add(singleIPs, big.NewInt(1))
		} else {
			ipRanges.Add(ipRanges, big.NewInt(1))
		}
	}

	threshold := tools.Ternary(m.config.IPv6MapThreshold > 0, m.config.IPv6MapThreshold, DefaultIPv6MapThreshold)
	thresholdBig := big.NewInt(int64(threshold))

	if totalIPs.Cmp(thresholdBig) <= 0 {
		m.expandAllIPv6ToMap(ranges)
	} else {
		m.optimizeIPv6MixedStrategy(ranges, thresholdBig)
	}
}

// expandAllIPv6ToMap 将所有IPv6范围展开为离散IP
func (m *IPMatcher) expandAllIPv6ToMap(ranges []ipv6Range) {
	m.storage.ipv6Map = make(map[[16]byte]struct{})
	for _, r := range ranges {
		current := r.start
		for {
			m.storage.ipv6Map[current.As16()] = struct{}{}

			if current.Compare(r.end) >= 0 {
				break
			}
			current = current.Next()
		}
	}
	// 清空范围存储，只使用map
	m.storage.ipv6Ranges = nil
}

// optimizeIPv6MixedStrategy IPv6混合策略优化
func (m *IPMatcher) optimizeIPv6MixedStrategy(ranges []ipv6Range, threshold *big.Int) {
	// 重新初始化map（用于存储部分IP）
	m.storage.ipv6Map = make(map[[16]byte]struct{})

	usedCapacity := &big.Int{}
	ipRanges := make([]ipv6Range, 0)

	// 分离单个IP和范围IP
	for _, r := range ranges {
		if r.start.Compare(r.end) == 0 {
			m.storage.ipv6Map[r.start.As16()] = struct{}{}
			usedCapacity.Add(usedCapacity, big.NewInt(1))
		} else {
			ipRanges = append(ipRanges, r)
		}
	}

	// 单IP已经超过阈值，range类型就不拆分展开了
	if usedCapacity.Cmp(threshold) > 0 {
		m.storage.ipv6Ranges = ipRanges
		return
	}

	// 对ipRanges按照范围大小进行排序，小范围的优先
	sort.Slice(ipRanges, func(i, j int) bool {
		sizeI := calculateIPv6RangeSize(ipRanges[i].start, ipRanges[i].end)
		sizeJ := calculateIPv6RangeSize(ipRanges[j].start, ipRanges[j].end)
		return sizeI.Cmp(sizeJ) < 0
	})

	// 根据剩余容量依次展开范围
	remainingRanges := make([]ipv6Range, 0)
	remainingCapacity := &big.Int{}
	remainingCapacity.Sub(threshold, usedCapacity)
	ipRangesLen := len(ipRanges)
	for idx, r := range ipRanges {
		rangeSize := calculateIPv6RangeSize(r.start, r.end)

		if rangeSize.Cmp(remainingCapacity) <= 0 {
			// 可以完整展开整个范围
			current := r.start
			for {
				m.storage.ipv6Map[current.As16()] = struct{}{}

				if current.Compare(r.end) >= 0 {
					break
				}
				current = current.Next()
			}
			usedCapacity.Add(usedCapacity, rangeSize)
			remainingCapacity.Sub(remainingCapacity, rangeSize)
		} else if remainingCapacity.Cmp(big.NewInt(0)) > 0 {
			// 部分展开：只展开剩余容量能容纳的部分
			current := r.start
			count := &big.Int{}

			for count.Cmp(remainingCapacity) < 0 && current.Compare(r.end) <= 0 {
				m.storage.ipv6Map[current.As16()] = struct{}{}
				current = current.Next()
				count.Add(count, big.NewInt(1))
			}

			usedCapacity.Add(usedCapacity, remainingCapacity)

			// 剩余部分保留为范围
			if current.Compare(r.end) <= 0 {
				remainingRanges = append(remainingRanges, ipv6Range{
					start: current,
					end:   r.end,
				})
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

	// 对剩余范围重新排序和合并
	if len(remainingRanges) == 0 {
		return
	}

	// 按起始IP排序（为了二分查找）
	sort.Slice(remainingRanges, func(i, j int) bool {
		return remainingRanges[i].start.Less(remainingRanges[j].start)
	})
	m.storage.ipv6Ranges = remainingRanges
}

// calculateIPv6RangeSize 计算IPv6范围的大小
func calculateIPv6RangeSize(start, end netip.Addr) *big.Int {
	if start.Compare(end) > 0 {
		return big.NewInt(0)
	}

	// 将IPv6地址转换为big.Int进行计算
	startBytes := start.As16()
	endBytes := end.As16()

	startInt := new(big.Int).SetBytes(startBytes[:])
	endInt := new(big.Int).SetBytes(endBytes[:])

	// 计算范围大小：end - start + 1
	size := new(big.Int).Sub(endInt, startInt)
	size.Add(size, big.NewInt(1))

	return size
}
