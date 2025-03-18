package formatRuleEngine

import "fmt"

type FormatRule[T any] struct {
	FormatType string   `yaml:"format_type" json:"format_type"` // 格式化类型
	InputRules []string `yaml:"input_rules" json:"input_rules"` // 识别格式化规则
	OutputRule string   `yaml:"output_rule" json:"output_rule"` //  输出格式规则
}

// Format 格式化
func (this FormatRule[T]) Format(src any) (T, error) {
	var (
		zero   T // 创建一个 T 类型的零值
		result any
		err    error
	)
	switch this.FormatType {
	case "date_format":
		result, err = this.dateFormat(src)
	case "output_date":
		result, err = this.dateObjectFormat(src)
	default:
		return zero, fmt.Errorf("不支持的格式化类型：%s", this.FormatType)
	}
	if err != nil {
		return zero, err
	}
	return result.(T), nil
}
