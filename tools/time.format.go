package tools

import (
	"strings"
	"time"
)

// TimePrecision 表示时间精度
type TimePrecision int

const (
	PrecisionYear TimePrecision = iota
	PrecisionMonth
	PrecisionWeek
	PrecisionDay
	PrecisionHour
	PrecisionMinute
	PrecisionSecond
	PrecisionUnknown
)

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
	precision := GetPrecisionFromFormat(format)
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
