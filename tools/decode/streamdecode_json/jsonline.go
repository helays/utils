package streamdecode_json

import (
	"bufio"
	"encoding/json"
	"io"

	"helay.net/go/utils/v3/dataType/customWriter"
)

func (i *Import) processLine(handler JSONHandler) (totalSize int64, err error) {
	// 按行读取内容
	counter := &customWriter.SizeCounter{}
	teeReader := io.TeeReader(i.rd, counter)
	scanner := bufio.NewScanner(teeReader)
	for scanner.Scan() {
		select {
		case <-i.ctx.Done():
			return counter.TotalSize, i.ctx.Err()
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

// processLineEnhanced 优化处理超长行
func (i *Import) processLineEnhanced(handler JSONHandler) (totalSize int64, err error) {
	counter := &customWriter.SizeCounter{}
	teeReader := io.TeeReader(i.rd, counter)
	reader := bufio.NewReader(teeReader)
	for {
		select {
		case <-i.ctx.Done():
			return counter.TotalSize, i.ctx.Err()
		default:
			line, _err := i.readFullLine(reader)
			if _err != nil {
				if _err == io.EOF {
					return counter.TotalSize, nil
				}
				return counter.TotalSize, _err
			}
			if len(line) > 0 {
				var dst map[string]any
				if err = json.Unmarshal(line, &dst); err != nil {
					return counter.TotalSize, err
				}
				if err = handler(i.ctx, dst); err != nil {
					return counter.TotalSize, err
				}
			}
		}
	}
}

// readFullLine 读取完整的一行（处理超长行）
func (i *Import) readFullLine(reader *bufio.Reader) ([]byte, error) {
	var line []byte
	var isPrefix bool
	var err error

	for {
		var partial []byte
		partial, isPrefix, err = reader.ReadLine()
		if err != nil {
			break
		}

		line = append(line, partial...)
		if !isPrefix {
			break
		}
	}

	// 处理EOF时可能有部分数据的情况
	if err == io.EOF && len(line) > 0 {
		err = nil // 重置错误，让外层处理数据
	}

	return line, err
}
