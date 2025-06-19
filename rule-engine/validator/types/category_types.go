package types

type ValidationCategory string // 条件运算分类

const (
	CategoryDataType ValidationCategory = "data_type" // 数据类型校验（原基础格式的子类）
	CategoryLength   ValidationCategory = "length"    // 长度校验（原基础格式的子类）
	CategoryFormat   ValidationCategory = "format"    // 格式校验（原基础格式的子类）

	CategoryContent  ValidationCategory = "content"  // 内容校验
	CategoryAdvanced ValidationCategory = "advanced" // 高级校验
)

type CategoryAttribute struct {
	Category ValidationCategory     `json:"category"`
	Zh       string                 `json:"zh"`
	List     map[Operator]Attribute `json:"list"`
}

var CategoryChineseNames = []CategoryAttribute{
	{Category: CategoryDataType, Zh: "数据类型校验", List: DataTypeAttributes},
	{Category: CategoryLength, Zh: "长度校验", List: LengthAttributes},
	{Category: CategoryFormat, Zh: "格式校验", List: FormatAttributes},
	{Category: CategoryContent, Zh: "内容校验", List: ContentAttributes},
	{Category: CategoryAdvanced, Zh: "安全校验", List: AdvancedAttributes},
}
