package ipmatch

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

func NewIPMatcher(config Config) *IPMatcher {
	m := &IPMatcher{
		config: config,
	}
	m.temp = newIPStorage() // 构建时候的缓存先定义
	return m
}

func (m *IPMatcher) LoadRule() error {
	// TODO: 添加ipv4 规则
	panic("TODO：待实现")
}

func (m *IPMatcher) AddIPv4Rule(ip string) error {
	// TODO: 添加ipv4 规则
	panic("TODO：待实现")
}

func (m *IPMatcher) AddIPv6Rule(ip string) error {
	// TODO: 添加ipv6 规则
	panic("TODO：待实现")
}

func (m *IPMatcher) Build() error {
	defer m.clearTemp()
	// TODO: 构建实现
	panic("TODO：待实现")
}
