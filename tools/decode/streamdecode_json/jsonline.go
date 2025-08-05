package streamdecode_json

import (
	"bufio"
	"encoding/json"
	"github.com/helays/utils/v2/dataType/customWriter"
	"io"
)

func (i *Import) processLine(handler JSONHandler) (totalSize int64, err error) {
	// 按行读取内容
	counter := &customWriter.SizeCounter{}
	teeReader := io.TeeReader(i.rd, counter)
	scanner := bufio.NewScanner(teeReader)
	for scanner.Scan() {
		select {
		case <-i.ctx.Done():
			return 0, i.ctx.Err()
		default:
			line := scanner.Text()
			var dst map[string]any
			if err = json.Unmarshal([]byte(line), &dst); err != nil {
				return counter.TotalSize, err
			}
			if err = handler(i.ctx, dst); err != nil {
				return counter.TotalSize, err
			}
		}
	}
	return counter.TotalSize, nil
}
