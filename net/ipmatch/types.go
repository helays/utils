package ipmatch

import (
	"net/netip"
	"sync"
	"time"

	"github.com/helays/utils/v2/safe/safettl"
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
	IPv4MapThreshold uint64 `json:"ipv4_map_threshold" yaml:"ipv4_map_threshold" ini:"ipv4_map_threshold"` // IPv4 Map模式阈值
	IPv6MapThreshold uint64 `json:"ipv6_map_threshold" yaml:"ipv6_map_threshold" ini:"ipv6_map_threshold"` // IPv6 Map模式阈值

	IPv4RuleSet IPRuleSet `json:"ipv4_rule_set" yaml:"ipv4_rule_set" ini:"ipv4_rule_set"`
	IPv6RuleSet IPRuleSet `json:"ipv6_rule_set" yaml:"ipv6_rule_set" ini:"ipv6_rule_set"`

	IPv4CacheTTL time.Duration `json:"ipv4_cache_ttl" yaml:"ipv4_cache_ttl" ini:"ipv4_cache_ttl"`
	IPv6CacheTTL time.Duration `json:"ipv6_cache_ttl" yaml:"ipv6_cache_ttl" ini:"ipv6_cache_ttl"`
}

type IPRuleSet struct {
	Rule     []string `json:"rule" yaml:"rule" ini:"rule"`                // 规则
	RuleFile []string `json:"rule_file" yaml:"rule_file" ini:"rule_file"` // 规则文件
}

// 默认阈值常量
const (
	DefaultIPv4MapThreshold = 100000 // IPv4默认阈值：10万IP
	DefaultIPv6MapThreshold = 50000  // IPv6默认阈值：5万IP（IPv6范围通常较大）
)

// IPMatcher 主匹配器结构
type IPMatcher struct {
	enable bool // 构建完成后，如果有ip，应该改成true。
	config *Config

	mu sync.RWMutex // 读写锁，当启用动态添加规则时，这个应该使用

	// 用于二分查询中的加速
	ipv4Cache *safettl.Map[uint32, struct{}]   // IPv4缓存
	ipv6Cache *safettl.Map[[16]byte, struct{}] // IPv6缓存

	// 正式存储
	storage *ipStorage

	// 临时存储，构建过程用
	temp *ipTemp
}

// ipStorage 存储结构
type ipStorage struct {
	ipv4Map map[uint32]struct{}   // IPv4: uint32 -> bool
	ipv6Map map[[16]byte]struct{} // IPv6: [16]byte -> bool

	ipv4Ranges []ipv4Range // IPv4连续范围
	ipv6Ranges []ipv6Range // IPv6连续范围
}

type ipTemp struct {
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

func newIPTemp() *ipTemp {
	return &ipTemp{
		ipv4Ranges: make([]ipv4Range, 0), // 建议也初始化
		ipv6Ranges: make([]ipv6Range, 0), // 建议也初始化
	}
}

func (m *IPMatcher) clearTemp() {
	m.temp = nil
}

func (m *IPMatcher) Close() {
	m.ipv4Cache.Close()
	m.ipv6Cache.Close()
}
