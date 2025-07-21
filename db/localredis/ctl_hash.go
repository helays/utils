package localredis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
)

// HGet 实现
func (l *LocalCache) HGet(ctx context.Context, key, field string) *redis.StringCmd {
	cmd := redis.NewStringCmd(ctx)

	l.hMu.RLock()
	h, ok := l.hData[key]
	l.hMu.RUnlock()

	if !ok {
		cmd.SetErr(redis.Nil)
		return cmd
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	if val, ok := h.data[field]; ok {
		// 自动转换为字符串
		switch v := val.(type) {
		case string:
			cmd.SetVal(v)
		case []byte:
			cmd.SetVal(string(v))
		case int, int64, float64:
			cmd.SetVal(fmt.Sprintf("%v", v))
		default:
			cmd.SetVal(fmt.Sprintf("%v", v))
		}
	} else {
		cmd.SetErr(redis.Nil)
	}
	return cmd
}

// HSet 实现
func (l *LocalCache) HSet(ctx context.Context, key string, values ...interface{}) *redis.IntCmd {
	cmd := redis.NewIntCmd(ctx)

	// 处理map参数
	if len(values) == 1 {
		if m, ok := values[0].(map[string]interface{}); ok {
			values = make([]interface{}, 0, len(m)*2)
			for k, v := range m {
				values = append(values, k, v)
			}
		}
	}

	if len(values)%2 != 0 {
		cmd.SetErr(ErrInvalidArgNum)
		return cmd
	}

	l.hMu.Lock()
	h, ok := l.hData[key]
	if !ok {
		h = &hMap{data: make(map[string]interface{})} // 修改初始化
		l.hData[key] = h
	}
	l.hMu.Unlock()

	h.mu.Lock()
	defer h.mu.Unlock()

	count := int64(0)
	for i := 0; i < len(values); i += 2 {
		if field, ok := values[i].(string); ok {
			h.data[field] = values[i+1] // 直接存储原始值
			count++
		}
	}

	cmd.SetVal(count)
	return cmd
}

// HGetAll 实现
func (l *LocalCache) HGetAll(ctx context.Context, key string) *redis.MapStringStringCmd {
	cmd := redis.NewMapStringStringCmd(ctx)

	l.hMu.RLock()
	h, ok := l.hData[key]
	l.hMu.RUnlock()

	if !ok {
		cmd.SetErr(redis.Nil) // 键不存在时返回 redis.Nil
		return cmd
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	result := make(map[string]string, len(h.data))
	for field, value := range h.data {
		// 统一转换为字符串
		switch v := value.(type) {
		case string:
			result[field] = v
		case []byte:
			result[field] = string(v)
		case int, int8, int16, int32, int64,
			uint, uint8, uint16, uint32, uint64:
			result[field] = fmt.Sprintf("%d", v)
		case float32, float64:
			result[field] = strconv.FormatFloat(v.(float64), 'f', -1, 64)
		case bool:
			if v {
				result[field] = "1"
			} else {
				result[field] = "0"
			}
		case nil:
			result[field] = ""
		default:
			// 其他类型使用 JSON 序列化
			jsonData, err := json.Marshal(v)
			if err != nil {
				result[field] = fmt.Sprintf("%v", v)
			} else {
				result[field] = string(jsonData)
			}
		}
	}

	cmd.SetVal(result)
	return cmd
}

// HMGet 实现
func (l *LocalCache) HMGet(ctx context.Context, key string, fields ...string) *redis.SliceCmd {
	cmd := redis.NewSliceCmd(ctx)
	result := make([]interface{}, len(fields))

	l.hMu.RLock()
	h, ok := l.hData[key]
	l.hMu.RUnlock()

	if !ok {
		// 键不存在时所有字段返回nil
		for i := range result {
			result[i] = nil
		}
		cmd.SetVal(result)
		return cmd
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	for i, field := range fields {
		if val, exists := h.data[field]; exists {
			// 类型转换逻辑（与HGetAll保持一致）
			switch v := val.(type) {
			case string:
				result[i] = v
			case []byte:
				result[i] = string(v)
			case int, int8, int16, int32, int64,
				uint, uint8, uint16, uint32, uint64:
				result[i] = fmt.Sprintf("%d", v)
			case float32, float64:
				result[i] = strconv.FormatFloat(v.(float64), 'f', -1, 64)
			case bool:
				if v {
					result[i] = "1"
				} else {
					result[i] = "0"
				}
			case nil:
				result[i] = ""
			default:
				jsonData, err := json.Marshal(v)
				if err != nil {
					result[i] = fmt.Sprintf("%v", v)
				} else {
					result[i] = string(jsonData)
				}
			}
		} else {
			result[i] = nil
		}
	}

	cmd.SetVal(result)
	return cmd
}
