package tools

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// AutoDetectTimestampString 解析字符串形式的时间戳并转换为 time.Time 对象
// timestampStr 是时间戳字符串
// 这个函数在最近一百年内还是挺好使的哈，过了2286-11-21 01:46:39 就不管用了
func AutoDetectTimestampString(timestampStr string) (time.Time, error) {
	// 尝试将字符串转换为 int64 类型的时间戳
	var timestamp int64
	var err error
	if strings.ContainsAny(timestampStr, ".") {
		// 如果包含小数点，可能是浮点数时间戳（秒），需要特殊处理
		fTimestamp, convErr := strconv.ParseFloat(timestampStr, 64)
		if convErr != nil {
			return time.Time{}, fmt.Errorf("cannot parse the provided timestamp: %v", convErr)
		}
		timestamp = int64(fTimestamp)
	} else {
		timestamp, err = strconv.ParseInt(timestampStr, 10, 64)
		if err != nil {
			return time.Time{}, fmt.Errorf("cannot parse the provided timestamp: %v", err)
		}
	}

	length := len(Any2string(timestamp))
	switch {
	case length <= 10:
		// 如果时间戳的位数不超过10位，可能是秒级时间戳
		return time.Unix(timestamp, 0), nil
	case length <= 13:
		// 如果时间戳的位数不超过13位，可能是毫秒级时间戳
		return time.UnixMilli(timestamp), nil
	default:
		// 否则，认为是纳秒级时间戳
		seconds := timestamp / 1e9
		nanos := timestamp % 1e9
		if seconds < 0 || nanos < 0 || nanos >= 1e9 {
			return time.Time{}, fmt.Errorf("invalid timestamp value")
		}
		return time.Unix(seconds, nanos), nil
	}
}
