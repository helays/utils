package tools

import (
	"fmt"
	"strings"
	"time"
)

// TimePrecision 表示时间精度
type TimePrecision int

func (tp TimePrecision) TimeFormat(t time.Time) string {
	switch tp {
	case PrecisionWeek:
		return FormatWeek(t)
	case PrecisionUnknown:
		return ""
	default:
		return t.Format(GetFormatFromPrecision(tp))
	}
}

// IsSameTimePrecision 方法实现时间精度比较
func (tp TimePrecision) IsSameTimePrecision(t1, t2 time.Time) bool {
	// 容错处理：检查时间是否为零值
	if t1.IsZero() || t2.IsZero() {
		return false
	}

	// 优化：先比较可能更快的年、月、日等
	switch tp {
	case PrecisionYear:
		return t1.Year() == t2.Year()
	case PrecisionMonth:
		y1, m1, _ := t1.Date()
		y2, m2, _ := t2.Date()
		return y1 == y2 && m1 == m2
	case PrecisionWeek:
		// 使用ISO周比较，更准确处理跨年周
		y1, w1 := t1.ISOWeek()
		y2, w2 := t2.ISOWeek()
		return y1 == y2 && w1 == w2
	case PrecisionDay:
		// 优化：直接比较年月日
		y1, m1, d1 := t1.Date()
		y2, m2, d2 := t2.Date()
		return y1 == y2 && m1 == m2 && d1 == d2
	case PrecisionHour:
		// 优化：先比较日，如果不相同则直接返回
		if !PrecisionDay.IsSameTimePrecision(t1, t2) {
			return false
		}
		return t1.Truncate(time.Hour).Equal(t2.Truncate(time.Hour))
	case PrecisionMinute:
		// 优化：先比较小时
		if !PrecisionHour.IsSameTimePrecision(t1, t2) {
			return false
		}
		return t1.Truncate(time.Minute).Equal(t2.Truncate(time.Minute))
	case PrecisionSecond:
		// 优化：先比较分钟
		if !PrecisionMinute.IsSameTimePrecision(t1, t2) {
			return false
		}
		return t1.Truncate(time.Second).Equal(t2.Truncate(time.Second))
	default:
		return false
	}
}

const (
	_                TimePrecision = iota // 忽略 0
	PrecisionYear                         // 1
	PrecisionMonth                        // 2
	PrecisionWeek                         // 3
	PrecisionDay                          // 4
	PrecisionHour                         // 5
	PrecisionMinute                       // 6
	PrecisionSecond                       // 7
	PrecisionUnknown                      // 8
)

var TimePrecisionChinese = map[TimePrecision]string{
	PrecisionYear:   "年",
	PrecisionMonth:  "月",
	PrecisionWeek:   "周",
	PrecisionDay:    "日",
	PrecisionHour:   "小时",
	PrecisionMinute: "分钟",
	PrecisionSecond: "秒",
}

// GetPrecisionFromFormat 根据时间格式字符串获取时间精度
func GetPrecisionFromFormat(format string) TimePrecision {
	// 检查是否包含周数相关的格式
	if strings.Contains(format, "W") || strings.Contains(format, "Monday") {
		return PrecisionWeek
	}

	// 检查格式中是否包含秒级精度
	if strings.Contains(format, "05") {
		return PrecisionSecond
	}

	// 检查格式中是否包含分钟级精度
	if strings.Contains(format, "04") {
		return PrecisionMinute
	}

	// 检查格式中是否包含小时级精度
	if strings.Contains(format, "15") {
		return PrecisionHour
	}

	// 检查格式中是否包含天级精度
	if strings.Contains(format, "02") || strings.Contains(format, "_2") || strings.Contains(format, "Mon") || strings.Contains(format, "Monday") {
		return PrecisionDay
	}

	// 检查格式中是否包含月级精度
	if strings.Contains(format, "01") || strings.Contains(format, "1") || strings.Contains(format, "Jan") || strings.Contains(format, "January") {
		return PrecisionMonth
	}

	// 检查格式中是否包含年级精度
	if strings.Contains(format, "2006") || strings.Contains(format, "06") {
		return PrecisionYear
	}

	return PrecisionUnknown
}

// GetFormatFromPrecision 根据时间精度获取标准的时间格式字符串
func GetFormatFromPrecision(precision TimePrecision) string {
	switch precision {
	case PrecisionSecond:
		return "20060102150405"
	case PrecisionMinute:
		return "200601021504"
	case PrecisionHour:
		return "2006010215"
	case PrecisionDay:
		return "20060102"
	case PrecisionMonth:
		return "200601"
	case PrecisionYear:
		return "2006"
	default:
		return "" // 或返回默认格式，如 "2006-01-02"
	}
}

// FormatWeek 将时间格式化为202501（年+周数）格式
func FormatWeek(t time.Time) string {
	year, week := t.ISOWeek()
	return fmt.Sprintf("%04d%02d", year, week)
}

// GetDurationFromPrecision 根据时间精度返回对应的时间段
func GetDurationFromPrecision(precision TimePrecision) time.Duration {
	switch precision {
	case PrecisionYear:
		return 365 * 24 * time.Hour // 近似值
	case PrecisionMonth:
		return 30 * 24 * time.Hour // 近似值
	case PrecisionWeek:
		return 7 * 24 * time.Hour
	case PrecisionDay:
		return 24 * time.Hour
	case PrecisionHour:
		return time.Hour
	case PrecisionMinute:
		return time.Minute
	case PrecisionSecond:
		return time.Second
	default:
		return time.Hour // 默认按小时
	}
}

// GetNextPeriodStart 根据时间格式模板获取下一个周期的起始时间
func GetNextPeriodStart(t time.Time, format string) time.Time {
	precision := GetPrecisionFromFormat(format)

	switch precision {
	case PrecisionYear:
		return time.Date(t.Year()+1, 1, 1, 0, 0, 0, 0, t.Location())
	case PrecisionMonth:
		nextMonth := t.Month() + 1
		nextYear := t.Year()
		if nextMonth > 12 {
			nextMonth = 1
			nextYear++
		}
		return time.Date(nextYear, nextMonth, 1, 0, 0, 0, 0, t.Location())
	case PrecisionWeek:
		weekday := int(t.Weekday())
		if weekday == 0 { // Sunday
			weekday = 7
		}
		return t.AddDate(0, 0, 8-weekday).Truncate(24 * time.Hour)
	case PrecisionDay:
		return t.AddDate(0, 0, 1).Truncate(24 * time.Hour)
	case PrecisionHour:
		return t.Truncate(time.Hour).Add(time.Hour)
	case PrecisionMinute:
		return t.Truncate(time.Minute).Add(time.Minute)
	case PrecisionSecond:
		return t.Truncate(time.Second).Add(time.Second)
	default:
		return t
	}
}

// GetNextPeriodTime 获取下一个周期的当前时间点（保持相同时间位置）
func GetNextPeriodTime(t time.Time, format string) time.Time {
	return GetNextPeriodTimeByPrecision(t, GetPrecisionFromFormat(format))
}

func GetNextPeriodTimeByPrecision(t time.Time, precision TimePrecision) time.Time {
	switch precision {
	case PrecisionYear:
		return t.AddDate(1, 0, 0)
	case PrecisionMonth:
		return t.AddDate(0, 1, 0)
	case PrecisionWeek:
		return t.AddDate(0, 0, 7)
	case PrecisionDay:
		return t.AddDate(0, 0, 1)
	case PrecisionHour:
		return t.Add(time.Hour)
	case PrecisionMinute:
		return t.Add(time.Minute)
	case PrecisionSecond:
		return t.Add(time.Second)
	default:
		return t
	}
}
