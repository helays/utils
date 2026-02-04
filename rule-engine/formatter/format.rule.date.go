package formatter

import (
	"fmt"
	"time"

	"github.com/araddon/dateparse"
	"helay.net/go/utils/v3/tools"
)

// 时间格式化
func (f FormatRule[T]) dateFormat(_src any) (any, error) {
	var (
		t   time.Time
		err error
		ok  bool
		src string
	)
	if t, ok = _src.(time.Time); !ok {
		src = tools.Any2string(_src)
		t, err = dateparse.ParseLocal(src)
	}
	// 首先尝试使用 https://github.com/araddon/dateparse 库

	if err == nil {
		if f.OutputRule != "" {
			return t.Format(f.OutputRule), nil
		}
		return t.Format(time.DateTime), err
	}
	for _, format := range f.InputRules {
		if format == "timestamp" {
			t, err = tools.AutoDetectTimestampString(src)
		} else {
			t, err = time.Parse(format, src)
		}
		if err != nil {
			continue
		}
		if f.OutputRule != "" {
			return t.Format(f.OutputRule), nil
		}
		return t.Format(time.DateTime), err
	}
	return src, nil
}

// 时间格式化
func (f FormatRule[T]) dateObjectFormat(_src any) (any, error) {
	var (
		t   time.Time
		err error
		ok  bool
		src = tools.Any2string(_src)
	)
	// 首先尝试使用 https://github.com/araddon/dateparse 库
	if t, ok = _src.(time.Time); !ok {
		t, err = dateparse.ParseLocal(src)
	}
	if err == nil {
		return t, nil
	}
	for _, format := range f.InputRules {
		if format == "timestamp" {
			t, err = tools.AutoDetectTimestampString(src)
		} else {
			t, err = time.Parse(format, src)
		}
		if err != nil {
			continue
		}
		return t, nil
	}
	return t, fmt.Errorf("时间解析失败：%s", _src)
}
