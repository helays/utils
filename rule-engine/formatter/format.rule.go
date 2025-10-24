package formatter

import "fmt"

type FormatRule[T any] struct {
	FormatType string   `yaml:"format_type" json:"format_type"` // 格式化类型
	InputRules []string `yaml:"input_rules" json:"input_rules"` // 识别格式化规则
	OutputRule string   `yaml:"output_rule" json:"output_rule"` //  输出格式规则
}

// Format 格式化
func (f FormatRule[T]) Format(src any) (T, error) {
	var (
		zero   T // 创建一个 T 类型的零值
		result any
		err    error
	)
	switch f.FormatType {
	case "date_format":
		result, err = f.dateFormat(src)
	case "output_date":
		result, err = f.dateObjectFormat(src)
	default:
		return zero, fmt.Errorf("不支持的格式化类型：%s", f.FormatType)
	}
	if err != nil {
		return zero, err
	}
	return result.(T), nil
}
