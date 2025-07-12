package operators

import (
	"fmt"
	"github.com/helays/utils/v2/rule-engine/validator/types"
)

func ValidateAdvanced(operator types.Operator, value any, rule []any) (string, bool) {
	switch operator {
	case types.ExprEval:
	case types.ExprGolang:
	case types.ExprCEL:
	case types.ExprJSONLogic:
	default:
		return fmt.Sprintf("未知类型校验操作符：%s", operator), false

	}
	return "", true
}
