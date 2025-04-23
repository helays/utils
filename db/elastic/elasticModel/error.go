package elasticModel

type ESError struct {
	Status int `json:"status"` // HTTP状态码
	// 错误信息
	Error *ErrorDetail `json:"error"` // 错误详情
}
