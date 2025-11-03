package safettl

import "time"

// itemWithExpiry 存储值和过期时间
type itemWithExpiry[V any] struct {
	value      V
	expiryTime time.Time
	ttl        time.Duration
}
