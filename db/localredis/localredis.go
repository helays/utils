// noinspection all
package localredis

import (
	"sync"
	"time"
)

// LocalCache 实现redis.UniversalClient接口的高效本地缓存
type LocalCache struct {
	kvData    sync.Map         // 普通键值存储
	hData     map[string]*hMap // 哈希表存储
	hMu       sync.RWMutex     // 哈希表顶层锁
	expiryMap sync.Map         // 存储键的过期时间
	closeChan chan struct{}    // 用于关闭后台清理goroutine
}

// hMap 哈希表结构
type hMap struct {
	data map[string]any
	mu   sync.RWMutex
}

// noinspection all
func NewLocalCache() *LocalCache {
	lc := &LocalCache{
		hData:     make(map[string]*hMap),
		closeChan: make(chan struct{}),
	}
	go lc.cleanupExpiredKeys() // 启动后台清理goroutine
	return lc
}

// Close 关闭缓存，停止后台goroutine
func (l *LocalCache) Close() error {
	close(l.closeChan)
	return nil
}

// cleanupExpiredKeys 定期清理过期键
func (l *LocalCache) cleanupExpiredKeys() {
	ticker := time.NewTicker(1 * time.Minute) // 每分钟检查一次
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			now := time.Now()
			l.kvData.Range(func(key, value interface{}) bool {
				if expiry, ok := l.expiryMap.Load(key); ok {
					if expiry.(time.Time).Before(now) {
						l.kvData.Delete(key)
						l.expiryMap.Delete(key)
					}
				}
				return true
			})
		case <-l.closeChan:
			return
		}
	}
}
