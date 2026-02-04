package query

import (
	"strings"

	"github.com/helays/utils/v2/tools"
)

func (b *Builder) ToES() map[string]any {
	if len(b.Conditions) > 0 {
		return b.buildBoolQuery()
	}
	return b.buildLeafQuery()
}

func (b *Builder) buildBoolQuery() map[string]any {
	// noinspection SpellCheckingInspection
	var exprs = make([]map[string]any, 0)
	for _, cond := range b.Conditions {
		if expr := cond.ToES(); expr != nil {
			exprs = append(exprs, expr)
		}
	}
	if len(exprs) == 0 {
		return nil
	}
	booType := "must"
	if b.Type.ToLower() == Or {
		booType = "should"
	}
	return map[string]any{
		"bool": map[string]any{
			booType: exprs,
		},
	}
}

func (b *Builder) buildLeafQuery() map[string]any {
	field := b.getESField()
	value := b.getESValue()
	if value == nil && !isNullOperator(b.Operator) {
		return nil
	}

	switch strings.ToLower(b.Operator) {
	case "=", "in":
		if tools.IsArray(value) {
			return map[string]any{"terms": map[string]any{field: value}}
		}
		return map[string]any{"term": map[string]any{field: value}}
	case "!=", "<>", "not in":
		if tools.IsArray(value) {
			return map[string]any{
				"bool": map[string]any{
					"must_not": map[string]any{
						"terms": map[string]any{field: value},
					},
				},
			}
		}
		return map[string]any{
			"bool": map[string]any{
				"must_not": map[string]any{
					"term": map[string]any{field: value},
				},
			},
		}
	case ">":
		return map[string]any{"range": map[string]any{field: map[string]any{"gt": value}}}
	case ">=":
		return map[string]any{"range": map[string]any{field: map[string]any{"gte": value}}}
	case "<":
		return map[string]any{"range": map[string]any{field: map[string]any{"lt": value}}}
	case "<=":
		return map[string]any{"range": map[string]any{field: map[string]any{"lte": value}}}
	case "like":
		return map[string]any{"wildcard": map[string]any{field: map[string]any{"value": "*" + tools.Any2string(value) + "*"}}}
	case "not like":
		return map[string]any{
			"bool": map[string]any{
				"must_not": map[string]any{
					"wildcard": map[string]any{field: map[string]any{"value": "*" + tools.Any2string(value) + "*"}},
				},
			},
		}
	case "null":
		return map[string]any{"bool": map[string]any{"must_not": map[string]any{"exists": map[string]any{"field": field}}}}
	case "notnull", "not null":
		return map[string]any{"exists": map[string]any{"field": field}}
	}

	return nil

}

func (b *Builder) getESField() string {
	if b.FieldAdvance != nil {
		return b.FieldAdvance.Name
	}
	return b.Field
}

func (b *Builder) getESValue() any {
	advanceLen := len(b.ValueAdvance)
	if advanceLen > 0 {
		if advanceLen == 1 {
			return b.ValueAdvance[0].Name
		}
		values := make([]string, 0, advanceLen)
		for _, v := range b.ValueAdvance {
			values = append(values, v.Name)
		}
		return values
	}
	return b.Value
}
