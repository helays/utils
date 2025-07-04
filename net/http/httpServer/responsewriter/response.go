package responsewriter

import "net/http"

type StatusRecorder struct {
	http.ResponseWriter
	status       int   // 状态码（如200、404）
	bytesWritten int64 // 响应体字节数
}

func New(w http.ResponseWriter) *StatusRecorder {
	return &StatusRecorder{
		ResponseWriter: w,
		status:         http.StatusOK,
	}
}

// WriteHeader 捕获状态码
func (r *StatusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

// Write 统计写入的字节数
func (r *StatusRecorder) Write(b []byte) (int, error) {
	n, err := r.ResponseWriter.Write(b)
	r.bytesWritten += int64(n) // 累加字节数
	return n, err
}

func (r *StatusRecorder) GetStatus() int {
	return r.status
}

func (r *StatusRecorder) GetBytesWritten() int64 {
	return r.bytesWritten
}
