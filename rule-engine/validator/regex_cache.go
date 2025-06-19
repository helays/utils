package validator

import (
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
)

// RegexCache 封装正则表达式缓存
type RegexCache struct {
	cache   sync.Map   // 存储编译后的正则表达式
	count   int32      // 当前缓存数量
	maxSize int32      // 最大缓存数量
	mutex   sync.Mutex // 保护缓存操作
}

var globalRegexCache = NewRegexCache(10000) // 全局缓存实例

// NewRegexCache 创建新的缓存实例
func NewRegexCache(maxSize int32) *RegexCache {
	return &RegexCache{
		maxSize: maxSize,
	}
}

// Get 获取或编译正则表达式
func (rc *RegexCache) Get(pattern string) *regexp.Regexp {
	// 标准化pattern格式：将.field.*.subfield转为正则
	regexPattern := "^" + strings.ReplaceAll(
		strings.ReplaceAll(pattern, ".", "\\."),
		"*", "[^.]+",
	) + "$"

	// 第一层缓存检查（无锁）
	if cached, ok := rc.cache.Load(regexPattern); ok {
		return cached.(*regexp.Regexp)
	}

	// 缓存未命中，进入加锁流程
	rc.mutex.Lock()
	defer rc.mutex.Unlock()

	// 第二层缓存检查（双检查锁定）
	if cached, ok := rc.cache.Load(regexPattern); ok {
		return cached.(*regexp.Regexp)
	}

	// 执行缓存淘汰（如果超过最大大小）
	if atomic.LoadInt32(&rc.count) >= rc.maxSize {
		rc.evict(rc.maxSize / 2) // 淘汰一半缓存
	}

	// 编译新正则表达式
	regex := regexp.MustCompile(regexPattern)
	rc.cache.Store(regexPattern, regex)
	atomic.AddInt32(&rc.count, 1)
	return regex
}

// evict 执行缓存淘汰（LRU简化实现）
func (rc *RegexCache) evict(toRemove int32) {
	removed := int32(0)
	rc.cache.Range(func(key, value interface{}) bool {
		if removed < toRemove {
			rc.cache.Delete(key)
			atomic.AddInt32(&rc.count, -1)
			removed++
			return true // 继续迭代
		}
		return false // 停止迭代
	})
}

// Clear 清空缓存（释放资源）
func (rc *RegexCache) Clear() {
	rc.mutex.Lock()
	defer rc.mutex.Unlock()
	rc.cache.Range(func(key, value interface{}) bool {
		rc.cache.Delete(key)
		return true
	})
	atomic.StoreInt32(&rc.count, 0)
}

// GetGlobalCache 获取全局缓存实例
func GetGlobalCache() *RegexCache {
	return globalRegexCache
}
