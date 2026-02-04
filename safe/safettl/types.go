package safettl

import (
	"sync"
	"time"
)

// itemWithExpiry 存储值和过期时间
type itemWithExpiry[V any] struct {
	value      V
	expiryTime time.Time
	ttl        time.Duration
	mu         sync.RWMutex
}

func (i *itemWithExpiry[V]) getExpiry() time.Time {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.expiryTime
}

func (i *itemWithExpiry[V]) setExpiry(d ...time.Duration) {
	i.mu.Lock()
	defer i.mu.Unlock()
	if len(d) > 0 {
		i.expiryTime = time.Now().Add(d[0])
	} else {
		i.expiryTime = time.Now().Add(i.ttl)
	}
}

// isExpired 检查是否过期
func (i *itemWithExpiry[V]) isExpired() bool {
	// 如果 expiryTime 是零值，表示永不过期
	if i.expiryTime.IsZero() {
		return false
	}
	return time.Now().After(i.expiryTime)
}
