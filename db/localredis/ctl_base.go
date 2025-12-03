// noinspection all
package localredis

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// Get 实现，增加过期检查
func (l *LocalCache) Get(ctx context.Context, key string) *redis.StringCmd {
	cmd := redis.NewStringCmd(ctx)

	// 检查是否过期
	if expiry, ok := l.expiryMap.Load(key); ok {
		if expiry.(time.Time).Before(time.Now()) {
			l.kvData.Delete(key)
			l.expiryMap.Delete(key)
			cmd.SetErr(redis.Nil) // 键过期视为不存在
			return cmd
		}
	}

	if val, ok := l.kvData.Load(key); ok {
		if str, ok := val.(string); ok {
			cmd.SetVal(str)
			return cmd
		}
	}

	cmd.SetErr(redis.Nil) // 键不存在
	return cmd
}

// Set 实现
func (l *LocalCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	l.kvData.Store(key, value)

	if expiration > 0 {
		// 存储过期时间
		l.expiryMap.Store(key, time.Now().Add(expiration))
	} else {
		// 如果没有设置过期时间，确保删除可能存在的过期时间
		l.expiryMap.Delete(key)
	}

	return redis.NewStatusCmd(ctx, "OK")
}

// Del 实现
func (l *LocalCache) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	cmd := redis.NewIntCmd(ctx)
	count := int64(0)

	for _, key := range keys {
		// 删除普通键
		if _, loaded := l.kvData.LoadAndDelete(key); loaded {
			count++
		}
		l.expiryMap.Delete(key)

		// 删除哈希表
		l.hMu.Lock()
		if _, ok := l.hData[key]; ok {
			delete(l.hData, key)
			count++
		}
		l.hMu.Unlock()
	}

	cmd.SetVal(count)
	return cmd
}

// HIncrBy 对哈希表中指定字段的值进行整数增量操作
func (l *LocalCache) HIncrBy(ctx context.Context, key, field string, incr int64) *redis.IntCmd {
	cmd := redis.NewIntCmd(ctx)

	l.hMu.Lock()
	h, ok := l.hData[key]
	if !ok {
		h = &hMap{data: make(map[string]interface{})}
		l.hData[key] = h
	}
	l.hMu.Unlock()

	h.mu.Lock()
	defer h.mu.Unlock()

	var current int64
	switch val := h.data[field].(type) {
	case nil:
		current = 0 // 字段不存在视为0
	case int:
		current = int64(val)
	case int64:
		current = val
	case float64:
		current = int64(val) // 浮点数截断
	case string:
		// 尝试解析字符串
		if num, err := strconv.ParseInt(val, 10, 64); err == nil {
			current = num
		} else {
			cmd.SetErr(fmt.Errorf("ERR hash value is not an integer"))
			return cmd
		}
	default:
		cmd.SetErr(fmt.Errorf("ERR hash value is not an integer"))
		return cmd
	}

	newVal := current + incr
	h.data[field] = newVal // 存储为int64类型
	cmd.SetVal(newVal)

	return cmd
}
