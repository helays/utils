package query

import (
	"fmt"
	"strings"

	"github.com/helays/utils/v2/tools"
	"gorm.io/gorm/clause"
)

// 辅助函数
func isNullOperator(op string) bool {
	return tools.Contains([]string{"null", "not null", "notnull"}, op)
}

func ParseFieldToColumn(field string) (clause.Column, error) {
	field = strings.TrimSpace(field)
	if field == "" {
		return clause.Column{}, fmt.Errorf("字段名不能为空")
	}
	if isComplexExpression(field) {
		return clause.Column{Name: field, Raw: true}, nil
	}
	// 去除所有引号
	cleanField := removeAllQuotes(field)
	if cleanField == "" {
		return clause.Column{}, fmt.Errorf("字段名不能为空")
	}
	col := clause.Column{}
	parts := strings.Split(cleanField, ".")
	partsLen := len(parts)
	if partsLen == 1 {
		col.Name = cleanField
	} else if partsLen == 2 {
		col.Table = parts[0]
		col.Name = parts[1]
	} else {
		return clause.Column{}, fmt.Errorf("字段名格式错误")
	}
	return col, nil
}

// isComplexExpression 判断是否是复杂表达式
func isComplexExpression(field string) bool {
	return strings.Contains(field, "(") || // 函数调用
		strings.Contains(field, ")") ||
		strings.Contains(field, "->") || // JSON 操作
		strings.Contains(field, "->>") ||
		strings.Contains(field, " ") || // 包含空格
		strings.Contains(field, "*") || // 通配符
		strings.Contains(field, "?") || // 参数
		strings.Contains(field, "@") // 变量
}

// removeAllQuotes 去除所有引号
func removeAllQuotes(s string) string {
	s = strings.ReplaceAll(s, "`", "")
	s = strings.ReplaceAll(s, `"`, "")
	s = strings.ReplaceAll(s, "'", "")
	return s
}
