package dataType

import (
	"database/sql/driver"
	"mime"
	"net/http"
	"strings"
)

// ContentType 表示检测到的内容类型
type ContentType int

func (ct ContentType) Value() (driver.Value, error) {
	return int(ct), nil
}

// String 返回内容类型的字符串表示
func (ct ContentType) String() string {
	switch ct {
	case JSON:
		return "JSON"
	case PlainText:
		return "PlainText"
	case CSV:
		return "CSV"
	case Excel:
		return "Excel"
	default:
		return "Unknown"
	}
}

const (
	Unknown ContentType = iota
	JSON
	PlainText
	CSV
	Excel
)

// DetectContentType 从 HTTP 请求头中检测内容类型
func DetectContentType(header http.Header) ContentType {
	contentType := header.Get("Content-Type")
	if contentType == "" {
		return Unknown
	}

	// 解析 Content-Type，可能包含 charset 等参数
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return Unknown
	}

	switch {
	case mediaType == "application/json" || mediaType == "text/json":
		return JSON
	case mediaType == "text/plain":
		return PlainText
	case mediaType == "text/csv":
		return CSV
	case mediaType == "application/vnd.ms-excel":
		return Excel
	case mediaType == "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":
		return Excel
	case strings.HasPrefix(mediaType, "application/vnd.ms-excel"):
		return Excel
	default:
		return Unknown
	}
}
