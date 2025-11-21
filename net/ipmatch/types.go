package ipmatch

import (
	"net/netip"
	"sync"
)

// IPVersion IP版本类型
type IPVersion string

const (
	IPv4 IPVersion = "ipv4"
	IPv6 IPVersion = "ipv6"
)

// Config IP匹配器配置
type Config struct {
	Dynamic          bool   `json:"dynamic" yaml:"dynamic" ini:"dynamic"`                                  // 是否支持动态添加规则
	Ipv4MapThreshold uint64 `json:"ipv4_map_threshold" yaml:"ipv4_map_threshold" ini:"ipv4_map_threshold"` // IPv4 Map模式阈值
	Ipv6MapThreshold uint64 `json:"ipv6_map_threshold" yaml:"ipv6_map_threshold" ini:"ipv6_map_threshold"` // IPv6 Map模式阈值

	IPv4RuleSet IPRuleSet `json:"ipv4_rule_set" yaml:"ipv4_rule_set" ini:"ipv4_rule_set"`
	IPv6RuleSet IPRuleSet `json:"ipv6_rule_set" yaml:"ipv6_rule_set" ini:"ipv6_rule_set"`
}

type IPRuleSet struct {
	Allow     []string `json:"allow" yaml:"allow" ini:"allow"`
	AllowFile []string `json:"allow_file" yaml:"allow_file" ini:"allow_file"`
	Deny      []string `json:"deny" yaml:"deny" ini:"deny"`
	DenyFile  []string `json:"deny_file" yaml:"deny_file" ini:"deny_file"`
}

// DefaultConfig 返回默认配置
func DefaultConfig() Config {
	return Config{
		Ipv4MapThreshold: DefaultIPv4MapThreshold,
		Ipv6MapThreshold: DefaultIPv6MapThreshold,
		Dynamic:          false,
	}
}

// 默认阈值常量
const (
	DefaultIPv4MapThreshold = 100000 // IPv4默认阈值：10万IP
	DefaultIPv6MapThreshold = 50000  // IPv6默认阈值：5万IP（IPv6范围通常较大）
)

// IPMatcher 主匹配器结构
type IPMatcher struct {
	enable bool // 构建完成后，如果有ip，应该改成true。
	config Config

	mu sync.RWMutex // 读写锁，当启用动态添加规则时，这个应该使用

	// 统计数据
	ipv4Count BuildStats // IPv4统计
	ipv6Count BuildStats // IPv6统计

	// 正式存储
	storage *ipStorage

	// 临时存储，构建过程用
	temp *ipStorage
}

// BuildStats 构建统计信息
type BuildStats struct {
	TotalIPs  uint64 // 总IP数量：离散IP+cidr里面的IP数量+连续IP里面包含的数量
	IPRanges  int    // 连续范围数量：规则中的连续IP范围数量
	SingleIPs int    // 离散IP数量
}

// ipStorage 存储结构
type ipStorage struct {
	ipv4Map map[uint32]struct{}   // IPv4: uint32 -> bool
	ipv6Map map[[16]byte]struct{} // IPv6: [16]byte -> bool

	ipv4Ranges []ipv4Range // IPv4连续范围
	ipv6Ranges []ipv6Range // IPv6连续范围
}

type ipv4Range struct {
	Start uint32
	End   uint32
}

type ipv6Range struct {
	start netip.Addr
	end   netip.Addr
}

func newIPStorage() *ipStorage {
	return &ipStorage{
		ipv4Map:    make(map[uint32]struct{}),
		ipv6Map:    make(map[[16]byte]struct{}),
		ipv4Ranges: make([]ipv4Range, 0), // 建议也初始化
		ipv6Ranges: make([]ipv6Range, 0), // 建议也初始化
	}
}

func (m *IPMatcher) clearTemp() {
	m.temp = nil
}
