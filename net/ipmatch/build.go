package ipmatch

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/helays/utils/v2/close/vclose"
	"github.com/helays/utils/v2/net/ipkit"
	"github.com/helays/utils/v2/tools"
	"github.com/malfunkt/iprange"
	"go4.org/netipx"
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

func NewIPMatcher(config *Config) (*IPMatcher, error) {
	m := &IPMatcher{
		config: config,
	}
	m.temp = newIPStorage() // 构建时候的缓存先定义
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
	r, err := iprange.Parse(ip)
	if err != nil {
		return fmt.Errorf("IPV4规则%s解析失败：%v", ip, err)
	}
	startIP := r.Min.To4()
	endIP := r.Max.To4()
	if startIP == nil || endIP == nil {
		return fmt.Errorf("规则 %s 不是有效的IPv4地址", ip)
	}
	// 直接使用范围信息，避免Expand()的内存开销
	start := ipkit.Ip2Int(startIP)
	end := ipkit.Ip2Int(endIP)
	if start > end {
		return fmt.Errorf("规则 %s 的起始IP大于结束IP", ip)
	}
	if start == end {
		m.temp.ipv4Map[start] = struct{}{}
	} else {
		m.temp.ipv4Ranges = append(m.temp.ipv4Ranges, ipv4Range{Start: start, End: end})
	}
	return nil
}

func (m *IPMatcher) AddIPv6Rule(ip string) error {
	if m.config.Dynamic {
		m.mu.Lock()
		defer m.mu.Unlock()
	}
	ip = strings.TrimSpace(ip)
	if ip == "" {
		return nil
	}
	netipx.ParseIPRange()
	netipx.ParsePrefixOrAddr()
	// TODO: 添加ipv6 规则
	panic("TODO：待实现")
}

func (m *IPMatcher) Build() error {
	if m.config.Dynamic {
		m.mu.Lock()
		defer m.mu.Unlock()
	}
	defer m.clearTemp()
	// TODO: 构建实现
	panic("TODO：待实现")
}
