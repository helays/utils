package streamline

import (
	"bufio"
	"io"
	"os"
	"strings"

	"helay.net/go/utils/v3/close/vclose"
)

type StreamingLineRemover struct {
	file     *os.File
	filename string
}

func New(filename string) (*StreamingLineRemover, error) {
	file, err := os.OpenFile(filename, os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	return &StreamingLineRemover{file: file, filename: filename}, nil
}

type Callback func(line string) error

func (s *StreamingLineRemover) ProcessLines(callback Callback) error {
	defer vclose.Close(s.file)
	for {
		var fileSize int64
		if fileInfo, err := s.file.Stat(); err != nil {
			return err
		} else if fileInfo.Size() == 0 {
			break // 文件已空
		} else {
			fileSize = fileInfo.Size()
		}
		if _, err := s.file.Seek(0, io.SeekStart); err != nil {
			return err
		}
		reader := bufio.NewReader(s.file)
		firstLine, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return err
		}
		if err == io.EOF && firstLine == "" {
			break
		}

		// 获取reader已读取的字节数（包括缓冲）
		bytesRead := int64(len(firstLine))
		if !strings.HasSuffix(firstLine, "\n") {
			// 如果没有换行符，说明是最后一行
			bytesRead = fileSize
		}
		firstLine = strings.TrimSpace(firstLine)
		if err = callback(firstLine); err != nil {
			return err
		}

		if bytesRead < fileSize {
			// 移动到第一行之后的位置
			if _, err = s.file.Seek(bytesRead, io.SeekStart); err != nil {
				return err
			}
			remaining := make([]byte, fileSize-bytesRead)
			_, err = s.file.Read(remaining)
			if err != nil {
				return err
			}
			_ = s.file.Truncate(0)
			_, _ = s.file.Seek(0, io.SeekStart)
			_, err = s.file.Write(remaining)
			if err != nil {
				return err
			}
		} else {
			// 没有剩余内容，清空文件
			_ = s.file.Truncate(0)
			break
		}
	}
	return nil
}

func (s *StreamingLineRemover) Close() {
	vclose.Close(s.file)
}
