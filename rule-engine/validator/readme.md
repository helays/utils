# 规则引擎验证库使用说明

## 概述

这是一个功能强大的Go语言规则引擎验证库，用于对数据结构进行复杂的条件验证。它支持多种验证类型、逻辑组合和自定义规则，适用于各种数据验证场景。

## 主要特性

1. **多类型验证**：支持数据类型、长度、格式、内容和高级验证
2. **逻辑组合**：支持AND/OR逻辑组合，可构建复杂的验证规则树
3. **通配符支持**：支持字段名通配符匹配
4. **高效实现**：使用对象池和缓存优化性能
5. **国际化支持**：内置中文错误提示
6. **可扩展架构**：易于添加新的验证类型和操作符

## 验证类型分类

### 1. 数据类型验证 (CategoryDataType)
- 验证字段的数据类型是否符合预期
- 支持的类型：字符串、数字、整数、布尔值、时间、数组、对象、空值

### 2. 长度验证 (CategoryLength)
- 验证字符串或数组的长度
- 支持的操作：最小长度、最大长度、固定长度、长度范围

### 3. 格式验证 (CategoryFormat)
- 验证特定格式的数据
- 支持的格式：邮箱、URL、IP、手机号、身份证、信用卡、颜色值、JSON、XML、Base64、UUID

### 4. 内容验证 (CategoryContent)
- 验证字段内容的具体值
- 支持的操作：必填、非空白、枚举值、比较运算(大于/小于/等于等)、正则匹配

### 5. 高级验证 (CategoryAdvanced)
- 支持表达式验证
- 支持的表达式类型：逻辑表达式、Go表达式、CEL表达式、JSONLogic规则

## 使用说明

### 基本用法

```go
import (
	"helay.net/go/utils/v3/rule-engine/validator"
	"helay.net/go/utils/v3/rule-engine/validator/types"
)

func main() {
	// 1. 定义验证规则
	rule := &validator.Rule{
		Field:         "username",
		FieldDataType: "string",
		Category:      types.CategoryContent,
		Operator:      types.Required,
	}

	// 2. 准备待验证数据
	data := map[string]any{
		"username": "testuser",
	}

	// 3. 执行验证
	msg, ok := rule.Validate(data)
	if !ok {
		fmt.Println("验证失败:", msg)
	} else {
		fmt.Println("验证通过")
	}
}
```

### 复杂规则示例

```go
// AND逻辑组合示例
rule := &validator.Rule{
	Logic: validator.LogicAnd,
	Conditions: []*validator.Rule{
		{
			Field:         "age",
			FieldDataType: "int",
			Category:      types.CategoryContent,
			Operator:      types.GreaterEqual,
			Value:         []any{18},
		},
		{
			Field:         "email",
			FieldDataType: "string",
			Category:      types.CategoryFormat,
			Operator:      types.FormatEmail,
		},
	},
}

// OR逻辑组合示例
rule := &validator.Rule{
	Logic: validator.LogicOr,
	Conditions: []*validator.Rule{
		{
			Field:         "phone",
			FieldDataType: "string",
			Category:      types.CategoryFormat,
			Operator:      types.FormatPhone,
		},
		{
			Field:         "email",
			FieldDataType: "string",
			Category:      types.CategoryFormat,
			Operator:      types.FormatEmail,
		},
	},
}

// 嵌套逻辑示例
rule := &validator.Rule{
	Logic: validator.LogicAnd,
	Conditions: []*validator.Rule{
		{
			Logic: validator.LogicOr,
			Conditions: []*validator.Rule{
				{
					Field:         "phone",
					FieldDataType: "string",
					Category:      types.CategoryFormat,
					Operator:      types.FormatPhone,
				},
				{
					Field:         "email",
					FieldDataType: "string",
					Category:      types.CategoryFormat,
					Operator:      types.FormatEmail,
				},
			},
		},
		{
			Field:         "password",
			FieldDataType: "string",
			Category:      types.CategoryLength,
			Operator:      types.LenMin,
			Value:         []any{8},
		},
	},
}
```

### 通配符验证示例

```go
rule := &validator.Rule{
	Field:         "items.*.price",
	FieldDataType: "float",
	Category:      types.CategoryContent,
	Operator:      types.GreaterThan,
	Value:         []any{0},
}

data := map[string]any{
	"items": map[string]any{
		"item1": map[string]any{
			"price": 10.5,
		},
		"item2": map[string]any{
			"price": 8.0,
		},
	},
}
```

## 操作符参考

### 数据类型验证操作符
- `is_string` - 字符串类型
- `is_number` - 数字类型
- `is_integer` - 整数类型
- `is_boolean` - 布尔类型
- `is_time` - 时间类型
- `is_array` - 数组类型
- `is_object` - 对象类型
- `is_null` - 空值

### 长度验证操作符
- `min_length` - 最小长度
- `max_length` - 最大长度
- `exact_length` - 固定长度
- `length_range` - 长度范围

### 格式验证操作符
- `email` - 邮箱
- `url` - URL
- `ip` - IP地址
- `phone` - 手机号
- `id_card` - 身份证
- `credit_card` - 信用卡
- `hex_color` - 十六进制颜色
- `json` - JSON格式
- `xml` - XML格式
- `base64` - Base64编码
- `uuid` - UUID

### 内容验证操作符
- `required` - 必填
- `not_blank` - 非空白字符
- `in_enum` - 在枚举中
- `not_in_enum` - 不在枚举中
- `gt` - 大于
- `ge` - 大于等于
- `lt` - 小于
- `le` - 小于等于
- `eq` - 等于
- `ne` - 不等于
- `regex` - 正则匹配

### 高级验证操作符
- `expr_eval` - 逻辑表达式
- `expr_golang` - Go表达式
- `expr_cel` - CEL表达式
- `expr_jsonlogic` - JSONLogic规则

## 性能优化建议

1. **复用规则对象**：尽可能复用已创建的规则对象，避免重复解析
2. **使用对象池**：库内部已使用对象池优化栈操作，无需额外配置
3. **缓存正则表达式**：通配符模式会自动缓存，无需手动处理
4. **短路优化**：AND/OR逻辑会自动短路，减少不必要的验证

## 扩展开发

如需添加新的验证类型或操作符：

1. 在`types.Operator`中添加新的操作符常量
2. 在对应的`*Attributes`映射中添加操作符属性
3. 在对应的验证器文件中实现验证逻辑
4. 在`ContentChineseMap`中添加中文错误提示

## 注意事项

1. 日期类型验证需要特殊处理，建议明确指定日期格式
2. 比较操作符对数据类型敏感，确保比较的数据类型一致
3. 通配符验证性能低于直接字段访问，在性能敏感场景慎用
4. 错误消息默认使用中文，如需其他语言需自行扩展

这个验证库设计灵活且高效，适合各种复杂业务规则的验证场景。通过合理的规则组合，可以构建出强大的验证逻辑来满足业务需求。