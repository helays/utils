package cors

import (
	"net/http"
	"strconv"
	"strings"

	"helay.net/go/utils/v3/tools"
)

func (c *Config) Apply(w http.ResponseWriter, origin string) bool {
	if !c.Enabled {
		return true
	}
	// 同源请求，直接放行
	if origin == "" {
		return true
	}
	if !c.isOriginAllowed(origin) {
		if c.Strict {
			w.WriteHeader(http.StatusForbidden)
			return false
		}
		return true
	}
	// 设置 Allow-Origin
	if c.AllowOrigins[0] == "*" && !c.AllowCredentials {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	} else {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}
	if c.AllowCredentials {
		w.Header().Set("Access-Control-Allow-Credentials", "true")
	}
	if len(c.ExposeHeaders) > 0 {
		w.Header().Set("Access-Control-Expose-Headers", strings.Join(c.ExposeHeaders, ","))
	}
	return true
}

// HandlePreflight 预检请求处理
func (c *Config) HandlePreflight(w http.ResponseWriter, origin string) {
	if !c.Apply(w, origin) {
		return
	}
	// 预检请求专用的头
	if len(c.AllowMethods) > 0 {
		w.Header().Set("Access-Control-Allow-Methods", strings.Join(c.AllowMethods, ","))
	}

	if len(c.AllowHeaders) > 0 {
		w.Header().Set("Access-Control-Allow-Headers", strings.Join(c.AllowHeaders, ","))
	}

	if c.MaxAge > 0 {
		w.Header().Set("Access-Control-Max-Age", strconv.Itoa(c.MaxAge))
	}
	w.WriteHeader(http.StatusNoContent)
}

// 判断origin是否在允许的列表中
func (c *Config) isOriginAllowed(origin string) bool {
	if len(c.AllowOrigins) == 0 {
		return false
	}
	if c.AllowOrigins[0] == "*" && !c.AllowCredentials {
		return true
	}
	return tools.Contains(c.AllowOrigins, origin)
}
