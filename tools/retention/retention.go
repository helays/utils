// Package retention 是用在自动清理数据的场景。比如文件、数据表等的自动清理
// 做好配置后，提供一个名称清单，然后讲根据配置的规则进行识别
// 最后通过Run函数，可进行对应的处理
package retention

import (
	"sort"
	"strconv"
	"strings"
	"time"
)

// Criteria 定义排序依据类型
type Criteria string

const (
	ByTime   Criteria = "time"   // 按时间排序
	ByNumber Criteria = "number" // 按数字排序
	ByName   Criteria = "name"   // 按名称排序
)

// Order 定义排序顺序类型
type Order bool

const (
	Descending Order = false // 降序
	Ascending  Order = true  // 升序
)

// Config 保留管理器配置
type Config struct {
	MaxRetain      int      // 最大保留数量
	Prefix         string   // 项目前缀
	Delimiter      string   // 分隔符，默认为"_"
	TimeFormat     string   // 时间格式，默认为"20060102"
	Order          Order    // 排序顺序
	Criteria       Criteria // 排序依据
	IgnoreSuffixes []string // 要忽略的后缀列表
}

// Manager 保留管理器
type Manager struct {
	config     Config
	items      []item
	extractKey func(string) interface{}
}

type item struct {
	name string
	key  interface{}
}

// New 创建新的保留管理器
func New(config Config) *Manager {
	if config.Delimiter == "" {
		config.Delimiter = "_"
	}
	if config.TimeFormat == "" {
		config.TimeFormat = "20060102"
	}

	return &Manager{
		config: config,
	}
}

// WithCustomExtractor 设置自定义键提取器
func (m *Manager) WithCustomExtractor(fn func(string) any) *Manager {
	m.extractKey = fn
	return m
}

// AddItems 添加要管理的项目
func (m *Manager) AddItems(names []string) *Manager {
	for _, name := range names {
		it := m.parseItem(name)
		if it != nil {
			m.items = append(m.items, *it)
		}
	}
	return m
}

// Run 执行保留管理，传入删除回调函数
func (m *Manager) Run(onDelete func(string) error) error {
	if m.config.MaxRetain < 1 || len(m.items) <= m.config.MaxRetain {
		return nil
	}

	m.sortItems()

	for _, _item := range m.items[m.config.MaxRetain:] {
		if err := onDelete(_item.name); err != nil {
			return err
		}
	}

	return nil
}

func (m *Manager) parseItem(name string) *item {
	if m.extractKey != nil {
		return &item{
			name: name,
			key:  m.extractKey(name),
		}
	}

	if !strings.HasPrefix(name, m.config.Prefix+m.config.Delimiter) {
		return nil
	}
	baseName := m.removeSuffixes(name)
	suffix := baseName[len(m.config.Prefix+m.config.Delimiter):]

	switch m.config.Criteria {
	case ByTime:
		if t, err := time.Parse(m.config.TimeFormat, suffix); err == nil {
			return &item{
				name: name,
				key:  t,
			}
		}
	case ByNumber:
		if num, err := strconv.ParseInt(suffix, 10, 64); err == nil {
			return &item{
				name: name,
				key:  num,
			}
		}
	}

	// 默认按名称排序
	return &item{
		name: name,
		key:  name,
	}
}

func (m *Manager) removeSuffixes(name string) string {
	for _, suffix := range m.config.IgnoreSuffixes {
		if strings.HasSuffix(name, suffix) {
			return strings.TrimSuffix(name, suffix)
		}
	}
	return name
}

func (m *Manager) sortItems() {
	sort.Slice(m.items, func(i, j int) bool {
		a, b := m.items[i], m.items[j]

		less := compare(a.key, b.key)

		if m.config.Order == Descending {
			return !less
		}
		return less
	})
}

func compare(a, b interface{}) bool {
	switch aVal := a.(type) {
	case time.Time:
		bVal := b.(time.Time)
		return aVal.Before(bVal)
	case int64:
		bVal := b.(int64)
		return aVal < bVal
	case string:
		bVal := b.(string)
		return aVal < bVal
	default:
		return false
	}
}
