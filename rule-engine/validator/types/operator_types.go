package types

type Operator string //条件运算类型

type Attribute struct {
	Title       string `json:"title"`
	ValueNum    int    `json:"value_num"`
	DataType    string `json:"data_type"`
	Description string `json:"description"`
}

const (
	TypeString  Operator = "is_string"  // 字符串类型
	TypeNumber  Operator = "is_number"  // 数字类型（含整型/浮点）
	TypeInteger Operator = "is_integer" // 整型
	TypeBoolean Operator = "is_boolean" // 布尔类型
	TypeTime    Operator = "is_time"    // 时间类型
	TypeArray   Operator = "is_array"   // 数组类型
	TypeObject  Operator = "is_object"  // 对象类型
	TypeNull    Operator = "is_null"    // 空值
)

var DataTypeChineseMap = map[Operator]string{
	TypeString:  "必须是文本类型",
	TypeNumber:  "必须是数字",
	TypeInteger: "必须是整数",
	TypeBoolean: "必须是布尔值(true/false)",
	TypeTime:    "必须是有效日期时间",
	TypeArray:   "必须是数组",
	TypeObject:  "必须是对象",
	TypeNull:    "必须为空值",
}

var DataTypeAttributes = map[Operator]Attribute{
	TypeString:  {Title: "文本", ValueNum: 0, DataType: ""},
	TypeNumber:  {Title: "数字", ValueNum: 0, DataType: ""},
	TypeInteger: {Title: "整数", ValueNum: 0, DataType: ""},
	TypeBoolean: {Title: "布尔值", ValueNum: 0, DataType: ""},
	TypeTime:    {Title: "时间", ValueNum: -1, DataType: "", Description: "可手动指定模板，模板语法[语言:模板]，比如[php:Y-m-d H:i:s]"},
	TypeArray:   {Title: "数组", ValueNum: 0, DataType: ""},
	TypeObject:  {Title: "对象", ValueNum: 0, DataType: ""},
	TypeNull:    {Title: "空值", ValueNum: 0, DataType: ""},
}

// ================================================================================

const (
	LenMin    Operator = "min_length"   // 最小长度
	LenMax    Operator = "max_length"   // 最大长度
	LenEquals Operator = "exact_length" // 固定长度
	LenRange  Operator = "length_range" // 长度范围
)

var LengthChineseMap = map[Operator]string{
	LenMin:    "长度不能少于%d个字符",
	LenMax:    "长度不能超过%d个字符",
	LenEquals: "长度必须等于%d个字符",
	LenRange:  "长度必须在%d到%d个字符之间",
}

var LengthAttributes = map[Operator]Attribute{
	LenMin:    {Title: "最小长度", ValueNum: 1, DataType: "number"},
	LenMax:    {Title: "最大长度", ValueNum: 1, DataType: "number"},
	LenEquals: {Title: "固定长度", ValueNum: 1, DataType: "number"},
	LenRange:  {Title: "长度范围", ValueNum: 2, DataType: "number"},
}

// ================================================================================

const (
	FormatEmail      Operator = "email"
	FormatURL        Operator = "url"
	FormatIP         Operator = "ip"
	FormatPhone      Operator = "phone"
	FormatIDCard     Operator = "id_card"     // 身份证
	FormatCreditCard Operator = "credit_card" // 信用卡
	FormatHexColor   Operator = "hex_color"   // 十六进制颜色
	FormatJSON       Operator = "json"        // JSON格式
	FormatXML        Operator = "xml"         // XML格式
	FormatBase64     Operator = "base64"
	FormatUUID       Operator = "uuid"
)

var FormatChineseMap = map[Operator]string{
	FormatEmail:      "请输入有效的邮箱地址",
	FormatURL:        "请输入有效的URL",
	FormatIP:         "请输入有效的IP地址",
	FormatPhone:      "请输入有效的手机号",
	FormatIDCard:     "请输入有效的身份证号",
	FormatCreditCard: "请输入有效的信用卡号",
	FormatHexColor:   "请输入有效的十六进制颜色值（如#FFFFFF）",
	FormatJSON:       "必须是有效的JSON格式",
	FormatXML:        "必须是有效的XML格式",
	FormatBase64:     "必须是有效的Base64编码",
	FormatUUID:       "必须是有效的UUID",
}

var FormatAttributes = map[Operator]Attribute{
	FormatEmail:      {Title: "邮箱地址", ValueNum: 0, DataType: ""},
	FormatURL:        {Title: "URL", ValueNum: 0, DataType: ""},
	FormatIP:         {Title: "IP地址", ValueNum: 0, DataType: ""},
	FormatPhone:      {Title: "手机号", ValueNum: 0, DataType: ""},
	FormatIDCard:     {Title: "身份证号", ValueNum: 0, DataType: ""},
	FormatCreditCard: {Title: "信用卡号", ValueNum: 0, DataType: ""},
	FormatHexColor:   {Title: "十六进制颜色值", ValueNum: 0, DataType: ""},
	FormatJSON:       {Title: "JSON格式", ValueNum: 0, DataType: ""},
	FormatXML:        {Title: "XML格式", ValueNum: 0, DataType: ""},
	FormatBase64:     {Title: "Base64编码", ValueNum: 0, DataType: ""},
	FormatUUID:       {Title: "UUID", ValueNum: 0, DataType: ""},
}

// ================================================================================

const (
	Required     Operator = "required"    // 必填
	NotBlank     Operator = "not_blank"   // 非空白字符
	InEnum       Operator = "in_enum"     // 枚举值
	NotInEnum    Operator = "not_in_enum" // 不在枚举中
	GreaterThan  Operator = "gt"          // 大于
	GreaterEqual Operator = "ge"          // 大于等于
	LessThan     Operator = "lt"          // 小于
	LessEqual    Operator = "le"          // 小于等于
	Equal        Operator = "eq"          // 等于
	NotEqual     Operator = "ne"          // 不等于
	RegexMatch   Operator = "regex"       // 正则匹配
)

var ContentChineseMap = map[Operator]string{
	Required:     "为必填项",
	NotBlank:     "不能全是空白字符",
	InEnum:       "值必须在预选范围内：%v",
	NotInEnum:    "值不能在范围内：%v",
	GreaterThan:  "值必须大于%v",
	GreaterEqual: "值必须大于等于%v",
	LessThan:     "值必须小于%v",
	LessEqual:    "值必须小于等于%v",
	Equal:        "值必须等于%v",
	NotEqual:     "值不能等于%v",
	RegexMatch:   "格式不符合规则：%s",
}

var ContentAttributes = map[Operator]Attribute{
	Required:     {Title: "必填", ValueNum: 0, DataType: ""},
	NotBlank:     {Title: "非空白字符", ValueNum: 0, DataType: ""},
	InEnum:       {Title: "枚举值", ValueNum: -1, DataType: "array"},
	NotInEnum:    {Title: "不在枚举中", ValueNum: -1, DataType: "array"},
	GreaterThan:  {Title: "大于", ValueNum: 1, DataType: "number"},
	GreaterEqual: {Title: "大于等于", ValueNum: 1, DataType: "number"},
	LessThan:     {Title: "小于", ValueNum: 1, DataType: "number"},
	LessEqual:    {Title: "小于等于", ValueNum: 1, DataType: "number"},
	Equal:        {Title: "等于", ValueNum: 1, DataType: "number"},
	NotEqual:     {Title: "不等于", ValueNum: 1, DataType: "number"},
	RegexMatch:   {Title: "正则匹配", ValueNum: 1, DataType: "string"},
}

// ================================================================================

const (
	ExprEval      Operator = "expr_eval"      // 逻辑表达式求值
	ExprGolang    Operator = "expr_golang"    // Go语法表达式
	ExprCEL       Operator = "expr_cel"       // Google CEL表达式
	ExprJSONLogic Operator = "expr_jsonlogic" // JSONLogic规则
)

var AdvancedAttributes = map[Operator]Attribute{
	ExprEval:      {Title: "逻辑表达式", ValueNum: 1, DataType: "string"},
	ExprGolang:    {Title: "Go表达式", ValueNum: 1, DataType: "string"},
	ExprCEL:       {Title: "CEL表达式", ValueNum: 1, DataType: "string"},
	ExprJSONLogic: {Title: "JSONLogic规则", ValueNum: 1, DataType: "string"},
}
