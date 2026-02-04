package query

import (
	"strings"

	"gorm.io/gorm/clause"
)

func (b *Builder) ToGORM() clause.Expression {
	if len(b.Conditions) > 0 {
		// noinspection SpellCheckingInspection
		var exprs = make([]clause.Expression, 0)
		for _, item := range b.Conditions {
			expr := item.ToGORM()
			if expr == nil {
				continue
			}
			exprs = append(exprs, expr)
		}
		if len(exprs) == 0 {
			return nil
		} else if len(exprs) == 1 {
			return exprs[0]
		}
		logic := b.Type.ToLower()
		if logic == AND {
			return clause.And(exprs...)
		} else {
			return clause.Or(exprs...)
		}
	}
	val := b.getGormValue()
	if val == nil && !isNullOperator(b.Operator) {
		return nil
	}
	field, err := b.getGormField()
	if err != nil {
		return nil
	}
	switch strings.ToLower(b.Operator) {
	case "=", "in":
		return clause.Eq{Column: field, Value: val}
	case "!=", "<>", "not in":
		return clause.Neq{Column: field, Value: val}
	case ">":
		return clause.Gt{Column: field, Value: val}
	case ">=":
		return clause.Gte{Column: field, Value: val}
	case "<":
		return clause.Lt{Column: field, Value: val}
	case "<=":
		return clause.Lte{Column: field, Value: val}
	case "like":
		return clause.Like{Column: field, Value: val}
	case "not like", "notlike":
		return clause.Not(clause.Like{Column: field, Value: val})
	case "null":
		return clause.Eq{Column: field, Value: nil}
	case "notnull", "not null":
		return clause.Neq{Column: field, Value: nil}
	}
	return nil

}

func (b *Builder) getGormValue() any {
	advanceLen := len(b.ValueAdvance)
	if advanceLen > 0 {
		if advanceLen == 1 {
			return b.ValueAdvance[0]
		}
		return b.ValueAdvance
	}
	return b.Value
}

func (b *Builder) getGormField() (clause.Column, error) {
	if b.FieldAdvance != nil {
		return *b.FieldAdvance, nil
	}

	return ParseFieldToColumn(b.Field)
}
