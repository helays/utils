package customWriter

// SizeCounter 自定义的 Writer，用于统计读取的字节数
// 用途：
// 服务端转发r.Body的时候，不需要io.ReadAll 就可以获取到body的size
type SizeCounter struct {
	TotalSize int64
}

func (s *SizeCounter) Write(p []byte) (int, error) {
	n := len(p)
	s.TotalSize += int64(n)

	return n, nil
}
