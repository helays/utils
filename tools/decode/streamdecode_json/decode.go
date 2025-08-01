package streamdecode_json

import (
	"context"
	"github.com/goccy/go-json"
	"github.com/helays/utils/v2/dataType/customWriter"
	"io"
)

type JSONHandler func(ctx context.Context, obj map[string]interface{}) error

func ProcessJSON(ctx context.Context, r io.Reader, handler JSONHandler) (totalSize int64, err error) {
	// 创建计数reader
	counter := &customWriter.SizeCounter{}
	teeReader := io.TeeReader(r, counter)

	decoder := json.NewDecoder(teeReader)
	// 读取第一个token判断是数组还是对象
	t, err := decoder.Token()
	if err != nil {
		return 0, err
	}

	switch delim := t.(type) {
	case json.Delim:
		switch delim {
		case '{':
			var obj map[string]any
			if err = decoder.Decode(&obj); err != nil {
				return counter.TotalSize, err
			}
			if err = handler(ctx, obj); err != nil {
				return counter.TotalSize, err
			}
			return counter.TotalSize, nil
		case '[':
			for decoder.More() {
				select {
				case <-ctx.Done():
					return counter.TotalSize, ctx.Err()
				default:
					var obj map[string]any
					if err = decoder.Decode(&obj); err != nil {
						return counter.TotalSize, err
					}
					if err = handler(ctx, obj); err != nil {
						return counter.TotalSize, err
					}
				}
			}
			return counter.TotalSize, nil
		}
	}
	return counter.TotalSize, nil
}
