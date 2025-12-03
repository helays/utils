package request

import (
	"fmt"
	"net"
	"net/http"
	"net/url"

	"github.com/helays/utils/v2/tools"
)

// GetHeaderValueFormatInt 获取多个header字段，并转为int
func GetHeaderValueFormatInt(header http.Header, fields ...string) (map[string]int, error) {
	resp := make(map[string]int)
	for _, field := range fields {
		_v := header.Get(field)
		v, err := tools.Any2int(_v)
		if err != nil {
			return nil, fmt.Errorf("不合法的字段[%s]值 %s %v", field, _v, err)
		}
		resp[field] = int(v)
	}
	return resp, nil
}

func GetQueryValueFromRequest2Int(r *http.Request, fields string, defaultValue ...int) (int, bool) {
	return GetQueryValueFromQuery2Int(r.URL.Query(), fields, defaultValue...)
}

func GetQueryValueFromQuery2Int(qs url.Values, fields string, defaultValue ...int) (int, bool) {
	val := qs.Get(fields)
	if val == "" {
		if len(defaultValue) > 0 {
			return defaultValue[0], true
		}
		return 0, false
	}
	_v, err := tools.Any2int(val)
	if err != nil {
		return 0, true
	}
	return int(_v), true
}
