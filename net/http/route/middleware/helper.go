package middleware

import (
	"sort"
	"strconv"
	"strings"
)

type encodingWithWeight struct {
	encoding CompressionAlgorithm
	weight   float64
}

// parseAcceptEncoding 解析 Accept-Encoding 头
func parseAcceptEncoding(header string) CompressionAlgorithm {
	if header == "" {
		return ""
	}
	var encodings []encodingWithWeight
	parts := strings.Split(header, ",")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		encodingParts := strings.Split(part, ";")
		encodingStr := strings.TrimSpace(encodingParts[0])
		var encoding CompressionAlgorithm
		weight := 1.0 // 默认权重
		if encodingStr == "*" {
			encoding = Gzip
			weight = 0
		} else {
			encoding = CompressionAlgorithm(encodingStr)
			if _, ok := supportedCompressionAlgorithms[encoding]; !ok {
				continue
			}
		}

		// 解析 q 值
		for i := 1; i < len(encodingParts); i++ {
			param := strings.TrimSpace(encodingParts[i])
			if strings.HasPrefix(param, "q=") {
				if q, err := strconv.ParseFloat(param[2:], 64); err == nil {
					weight = q
				}
			}
		}

		encodings = append(encodings, encodingWithWeight{
			encoding: encoding,
			weight:   weight,
		})
	}

	// 按权重排序（从高到低）
	sort.Slice(encodings, func(i, j int) bool {
		return encodings[i].weight > encodings[j].weight
	})
	if len(encodings) > 0 {
		return encodings[0].encoding
	}
	return ""
}

func shouldCompress(contentType string) bool {
	if contentType == "" {
		return true
	}
	contentType = strings.ToLower(strings.TrimSpace(contentType))
	parts := strings.SplitN(contentType, "/", 2)
	switch parts[0] {
	case "image", "video", "audio", "font":
		return false
	}
	// 再次规范数据
	parts = strings.SplitN(contentType, ";", 2)
	if _, ok := excludeContentTypes[parts[0]]; ok {
		return false
	}
	return true
}
