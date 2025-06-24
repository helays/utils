package types

// ValidationError 定义验证错误详情
type ValidationError struct {
	Field      string             `json:"field"`       // 验证失败的字段名
	RealField  string             `json:"real_field"`  // 实际验证的字段名
	InputValue interface{}        `json:"input_value"` // 输入的原始值
	RuleValue  interface{}        `json:"rule_value"`  // 规则要求的比较值
	Operator   Operator           `json:"operator"`    // 使用的操作符
	Category   ValidationCategory `json:"category"`    // 验证类型分类
	Message    string             `json:"message"`     // 错误描述信息
}

// Error 实现error接口
func (e ValidationError) Error() string {
	return e.Message
}
