package ipmatch

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"net/netip"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/malfunkt/iprange"
	"go4.org/netipx"
	"helay.net/go/utils/v3/close/vclose"
	"helay.net/go/utils/v3/net/ipkit"
	"helay.net/go/utils/v3/safe"
	"helay.net/go/utils/v3/tools"
)

// 构建过程：
// 1. 解析规则 → 2. CIDR转Range  → 3. 排序Range → 4. 合并重复 → 5. 分离离散IP和连续范围

// 在添加规则的时候需要注意一下几点
// 支持ipv4 ipv6
// 对于 ipv4 可以用github.com/malfunkt/iprange 进行规则解析
// 对于 ipv6 可以用go4.org/netipx 进行解析。

// 构建过程注意一下几点
// 合并相同的规则，然后独立出离散ip和ip range
// cidr也需要转换成 range
// 对对最终的range 进行排序。

// 计算策略的时候需要注意一下几点。
// ip totals = 离散ip+range里面的数量
// ip totals 是否小于对应ip version 的阈值，如果小于阈值，则用map 直接查询
// 当ip totals大于阈值后，应该分别尝试 离散ip+range里面的数量， 获取取其中一部分，尽可能让map存储满。
// 剩下的range 进行排序，然后采用二分法查询

func NewIPMatcher(ctx context.Context, config *Config) (*IPMatcher, error) {
	m := &IPMatcher{
		config: config,
	}
	ipv4CacheTTL := tools.AutoTimeDuration(m.config.IPv4CacheTTL, time.Second, 30*time.Second)
	ipv6CacheTTL := tools.AutoTimeDuration(m.config.IPv6CacheTTL, time.Second, 10*time.Second)
	m.ipv4Cache = safe.NewMap[uint32, struct{}](ctx, safe.IntegerHasher[uint32]{}, safe.CacheConfig{
		EnableCleanup: true,
		ClearInterval: ipv4CacheTTL / 2,
		TTL:           ipv4CacheTTL,
	})
	m.ipv6Cache = safe.NewMap[[16]byte, struct{}](ctx, safe.Array16Hasher{}, safe.CacheConfig{
		EnableCleanup: true,
		ClearInterval: ipv6CacheTTL / 2,
		TTL:           ipv6CacheTTL,
	})
	m.temp = newIPTemp() // 构建时候的缓存先定义
	err := m.LoadRule()
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (m *IPMatcher) LoadRule() error {
	if m.config.Dynamic {
		m.mu.Lock()
		defer m.mu.Unlock()
	}
	ipv4RuleSet := m.config.IPv4RuleSet
	ipv6RuleSet := m.config.IPv6RuleSet

	for _, rule := range ipv4RuleSet.Rule {
		if err := m.AddIPv4Rule(rule); err != nil {
			return err
		}
	}
	for _, rule := range ipv6RuleSet.Rule {
		if err := m.AddIPv6Rule(rule); err != nil {
			return err
		}
	}

	// 读取文件列表

	// 读取ipv4的规则文件
	for _, file := range ipv4RuleSet.RuleFile {
		rules, err := m.loadFile(file)
		if err != nil {
			return err
		}
		for _, rule := range rules {
			if err = m.AddIPv4Rule(rule); err != nil {
				return err
			}
		}
	}

	// 读取ipv6的规则文件
	for _, file := range ipv6RuleSet.RuleFile {
		rules, err := m.loadFile(file)
		if err != nil {
			return err
		}
		for _, rule := range rules {
			if err = m.AddIPv6Rule(rule); err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *IPMatcher) loadFile(file string) ([]string, error) {
	filePath := tools.Fileabs(file)
	// 按行读取文件
	f, err := os.Open(filePath)
	defer vclose.Close(f)
	if err != nil {
		return nil, err
	}
	var rules = make([]string, 0)
	err = tools.ReadRowWithFile(f, func(scanner *bufio.Scanner) error {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			return nil
		}
		rules = append(rules, line)
		return nil
	})
	return rules, err
}

// AddIPv4Rule 添加ipv4规则
// 对于单个的ip，都定义为离散ip
// 对于多个的ip，都定义为连续范围
// 后面在构建的时候再进行排序合并
func (m *IPMatcher) AddIPv4Rule(ip string) error {
	if m.config.Dynamic {
		m.mu.Lock()
		defer m.mu.Unlock()
	}
	ip = strings.TrimSpace(ip)
	if ip == "" {
		return nil
	}
	var startIP, endIP net.IP
	if r, err := iprange.Parse(ip); err == nil {
		startIP = r.Min.To4()
		endIP = r.Max.To4()
	} else {
		start, end, _err := m.netIPxParse(ip)
		if _err != nil {
			return err
		}
		// 转成 net.IP
		startIP = start.AsSlice()
		endIP = end.AsSlice()
	}

	if startIP == nil || endIP == nil {
		return fmt.Errorf("规则 %s 不是有效的IPv4地址", ip)
	}
	// 直接使用范围信息，避免Expand()的内存开销
	start := ipkit.Ip2Int(startIP)
	end := ipkit.Ip2Int(endIP)
	if start > end {
		return fmt.Errorf("规则 %s 的起始IP大于结束IP", ip)
	}
	m.temp.ipv4Ranges = append(m.temp.ipv4Ranges, ipv4Range{Start: start, End: end})
	return nil
}

// AddIPv6Rule 添加ipv6规则
func (m *IPMatcher) AddIPv6Rule(ip string) error {
	if m.config.Dynamic {
		m.mu.Lock()
		defer m.mu.Unlock()
	}
	ip = strings.TrimSpace(ip)
	if ip == "" {
		return nil
	}
	start, end, err := m.netIPxParse(ip)
	if err != nil {
		return err
	}
	// 验证IPv6地址
	if !start.Is6() || !end.Is6() {
		return fmt.Errorf("IPv6地址 '%s' 不是有效的IPv6地址", ip)
	}
	if start.Compare(end) > 0 {
		return fmt.Errorf("IPv6地址范围 '%s' 的起始IP大于结束IP", ip)
	}
	m.temp.ipv6Ranges = append(m.temp.ipv6Ranges, ipv6Range{start: start, end: end})
	return nil
}

// netIPxParse 解析IP规则 采用netipx库解析
func (m *IPMatcher) netIPxParse(ip string) (netip.Addr, netip.Addr, error) {
	var start, end netip.Addr
	// 解析IP规则
	if strings.Contains(ip, "-") {
		// IP范围格式: 2001:db8::1-2001:db8::ffff
		ipRange, err := netipx.ParseIPRange(ip)
		if err != nil {
			return start, end, fmt.Errorf("无法解析IP地址范围 '%s' %v", ip, err)
		}
		start = ipRange.From()
		end = ipRange.To()
	} else if strings.Contains(ip, "/") {
		// CIDR格式: 2001:db8::/32
		prefix, err := netip.ParsePrefix(ip)
		if err != nil {
			return start, end, fmt.Errorf("未能解析IP CIDR '%s' %v", ip, err)
		}
		ipRange := netipx.RangeOfPrefix(prefix)
		start = ipRange.From()
		end = ipRange.To()
	} else {
		// 单个IP格式: 2001:db8::1
		addr, err := netip.ParseAddr(ip)
		if err != nil {
			return start, end, fmt.Errorf("无法解析地址 '%s': %w", ip, err)
		}
		start = addr
		end = addr
	}
	return start, end, nil
}

// Build 构建
func (m *IPMatcher) Build() {
	if m.config.Dynamic {
		m.mu.Lock()
		defer m.mu.Unlock()
	}
	m.enable = false
	defer m.clearTemp()

	// 构建正式存储
	m.storage = newIPStorage()

	// 处理ipv4
	if len(m.temp.ipv4Ranges) > 0 {
		sort.Slice(m.temp.ipv4Ranges, func(i, j int) bool {
			return m.temp.ipv4Ranges[i].Start < m.temp.ipv4Ranges[j].Start
		})
		ranges := m.mergeIPv4Ranges()
		m.optimizeIPv4Storage(ranges)
	}

	// 处理ipv6
	if len(m.temp.ipv6Ranges) > 0 {
		sort.Slice(m.temp.ipv6Ranges, func(i, j int) bool {
			return m.temp.ipv6Ranges[i].start.Less(m.temp.ipv6Ranges[j].start)
		})
		ranges := m.mergeIPv6Ranges()
		m.optimizeIPv6Storage(ranges)
	}
	if len(m.storage.ipv4Map) > 0 || len(m.storage.ipv6Map) > 0 || len(m.storage.ipv4Ranges) > 0 || len(m.storage.ipv6Ranges) > 0 {
		m.enable = true
	}
}

// 合并IPv4范围
func (m *IPMatcher) mergeIPv4Ranges() []ipv4Range {
	merged := make([]ipv4Range, 0)
	l := len(m.temp.ipv4Ranges)
	current := m.temp.ipv4Ranges[0]

	for i := 1; i < l; i++ {
		next := m.temp.ipv4Ranges[i]
		// 检查是否重叠或连续 (current.End + 1 >= next.Start)
		if current.End+1 >= next.Start {
			if next.End > current.End {
				current.End = next.End
			}
		} else {
			merged = append(merged, current)
			current = next
		}
	}
	merged = append(merged, current)

	return merged
}

// 合并IPv6范围
func (m *IPMatcher) mergeIPv6Ranges() []ipv6Range {
	merged := make([]ipv6Range, 0)
	l := len(m.temp.ipv6Ranges)
	current := m.temp.ipv6Ranges[0]

	for i := 1; i < l; i++ {
		next := m.temp.ipv6Ranges[i]
		// 检查是否重叠或连续
		if current.end.Next().Compare(next.start) >= 0 {
			if next.end.Compare(current.end) > 0 {
				current.end = next.end
			}
		} else {
			merged = append(merged, current)
			current = next
		}
	}
	merged = append(merged, current)
	return merged
}
