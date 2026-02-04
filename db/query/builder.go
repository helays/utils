package query

import (
	"strings"

	"gorm.io/gorm/clause"
)

// LogicOperator 逻辑操作符类型
type LogicOperator string

func (l LogicOperator) String() string {
	return string(l)
}

func (l LogicOperator) ToLower() LogicOperator {
	return LogicOperator(strings.ToLower(l.String()))
}

const (
	AND LogicOperator = "and"
	Or  LogicOperator = "or"
)

type Builder struct {
	Type         LogicOperator   `json:"type"`
	Field        string          `json:"field"`         // 普通字段
	FieldAdvance *clause.Column  `json:"field_advance"` // 高级字段配置
	OperatorType string          `json:"operator_type"` // 操作符类型
	Operator     string          `json:"operator"`
	Value        any             `json:"value"`
	ValueAdvance []clause.Column `json:"value_advance"` // 高级值配置
	Conditions   []Builder       `json:"conditions"`
}
