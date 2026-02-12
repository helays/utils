package tools

// Deprecated: 自 Go 1.26.0 起已弃用，请直接使用 new(T) 获取类型指针
// 示例: new(string) 等价于 ToPtr("")
func ToPtr[T any](v T) *T {
	return &v
}
