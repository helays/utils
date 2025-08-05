package streamdecode_json

import (
	"fmt"
	"github.com/goccy/go-json"
	"github.com/helays/utils/v2/dataType/customWriter"
	"io"
)

func (i *Import) processJSON(handler JSONHandler) (totalSize int64, err error) {
	// 创建计数reader
	counter := &customWriter.SizeCounter{}
	teeReader := io.TeeReader(i.rd, counter)
	decoder := json.NewDecoder(teeReader)

	// 读取第一个token判断是数组还是对象
	t, err := decoder.Token()
	if err != nil {
		return 0, err
	}

	switch delim := t.(type) {
	case json.Delim:
		switch delim {
		case '{': // 由于提前消费了 { 字符，所以这里需要手动进行解析。
			obj := make(map[string]any)
			for decoder.More() {
				// 读取 key
				key, _err := decoder.Token()
				if _err != nil {
					return counter.TotalSize, err
				}
				keyStr, ok := key.(string)
				if !ok {
					return counter.TotalSize, fmt.Errorf("expected string key, got %T", key)
				}
				// 读取 value
				var val any
				if err = decoder.Decode(&val); err != nil {
					return counter.TotalSize, err
				}
				obj[keyStr] = val
			}
			if err = handler(i.ctx, obj); err != nil {
				return counter.TotalSize, err
			}
			return counter.TotalSize, nil
		case '[':
			for decoder.More() {
				select {
				case <-i.ctx.Done():
					return counter.TotalSize, i.ctx.Err()
				default:
					var obj map[string]any
					if err = decoder.Decode(&obj); err != nil {
						return counter.TotalSize, err
					}
					if err = handler(i.ctx, obj); err != nil {
						return counter.TotalSize, err
					}
				}
			}
			return counter.TotalSize, nil
		}
	}
	return counter.TotalSize, nil
}
