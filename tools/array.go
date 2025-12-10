package tools

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/helays/utils/v2/config"
	"golang.org/x/exp/constraints"
)

// ArrayChunk 高性能泛型切片分块函数
func ArrayChunk[T any](slice []T, size int) [][]T {
	if size <= 0 {
		return [][]T{slice}
	}

	length := len(slice)
	if length == 0 {
		return nil
	}

	chunkCount := (length + size - 1) / size
	chunks := make([][]T, chunkCount)

	for i := 0; i < chunkCount; i++ {
		start := i * size
		end := start + size
		if end > length {
			end = length
		}
		// 直接引用原切片的底层数组，避免额外内存分配
		chunks[i] = slice[start:end:end] // 使用完整切片表达式限制容量
	}

	return chunks
}

// ArrayChunkCopy 高性能且安全的版本，会复制数据而不是引用原切片
func ArrayChunkCopy[T any](slice []T, size int) [][]T {
	if size <= 0 {
		copied := make([]T, len(slice))
		copy(copied, slice)
		return [][]T{copied}
	}

	length := len(slice)
	if length == 0 {
		return nil
	}

	chunkCount := (length + size - 1) / size
	chunks := make([][]T, chunkCount)

	for i := 0; i < chunkCount; i++ {
		start := i * size
		end := start + size
		if end > length {
			end = length
		}
		chunk := make([]T, end-start)
		copy(chunk, slice[start:end])
		chunks[i] = chunk
	}

	return chunks
}

// Slice2Map 更清晰的参数命名
func Slice2Map[Key comparable, Elem any](slice []Elem, keyFunc func(Elem) Key) map[Key]Elem {
	result := make(map[Key]Elem)
	for _, item := range slice {
		result[keyFunc(item)] = item
	}
	return result
}

// SliceToMultiMap 将切片转换为映射，允许键重复，相同键的值会保存在切片中
// slice: 要转换的切片
// keyFunc: 从元素中提取键的函数
func SliceToMultiMap[Key comparable, Elem any](slice []Elem, keyFunc func(Elem) Key) map[Key][]Elem {
	result := make(map[Key][]Elem)
	for _, item := range slice {
		key := keyFunc(item)
		result[key] = append(result[key], item)
	}
	return result
}

// RemoveDuplicates 对slice去重
func RemoveDuplicates[T comparable](slice []T) []T {
	encountered := map[T]struct{}{}
	result := make([]T, 0, len(slice))

	for _, v := range slice {
		if _, ok := encountered[v]; !ok {
			encountered[v] = struct{}{}
			result = append(result, v)
		}
	}

	return result
}

// RemoveDuplicatesWithKeyFunc 对slice去重，使用keyFunc生成比较键
func RemoveDuplicatesWithKeyFunc[T any, K comparable](slice []T, keyFunc func(T) K) []T {
	encountered := map[K]struct{}{}
	result := make([]T, 0, len(slice))

	for _, v := range slice {
		key := keyFunc(v)
		if _, exists := encountered[key]; !exists {
			encountered[key] = struct{}{}
			result = append(result, v)
		}
	}

	return result
}

func ContainsFilterHelper[T any](v T, filters ...func(T) T) T {
	if len(filters) < 1 {
		return v
	}
	for i := range filters {
		v = filters[i](v)
	}
	return v
}

// Contains 检查某个值是否在切片中（泛型实现）
func Contains[T comparable](slice []T, target T, filters ...func(T) T) bool {
	target = ContainsFilterHelper(target, filters...)
	for _, v := range slice {
		v = ContainsFilterHelper(v, filters...)
		if v == target {
			return true
		}
	}
	return false
}

// ContainsFunc 支持自定义比较逻辑（适用于不可比较类型，如 `struct` 带非 `comparable` 字段）
func ContainsFunc[T any](slice []T, target T, equal func(a, b T) bool) bool {
	for _, v := range slice {
		if equal(v, target) {
			return true
		}
	}
	return false
}

// ContainsByField 可用于[]struct中，检查某个字段的值是否存在
func ContainsByField[T any, F comparable](slice []T, target F, fieldExtractor func(T) F) bool {
	for _, item := range slice {
		if fieldExtractor(item) == target {
			return true
		}
	}
	return false
}

// ContainsAny 检查 `elems` 中是否有任意一个元素在 `targets` 里
func ContainsAny[T comparable](elems []T, targets []T) bool {
	if len(targets) == 0 || len(elems) == 0 {
		return false
	}
	// 选择较小的集合缓存
	var staticSet map[T]struct{}
	var dynamic []T
	if len(targets) < len(elems) {
		staticSet = make(map[T]struct{}, len(targets))
		for _, t := range targets {
			staticSet[t] = struct{}{}
		}
		dynamic = elems
	} else {
		staticSet = make(map[T]struct{}, len(elems))
		for _, e := range elems {
			staticSet[e] = struct{}{}
		}
		dynamic = targets
	}
	// 遍历动态集合
	for _, v := range dynamic {
		if _, ok := staticSet[v]; ok {
			return true
		}
	}
	return false
}

// ContainsAnyHashBest 检查 `elems` 中是否有任意一个元素在 `targets` 里，自定义比较函数
// elems 检查列表
// targets 被搜索列表
// hashFunc  计算hash
// equal 精确比较
func ContainsAnyHashBest[T any, H comparable](elems []T, targets []T, hashFunc func(T) H, equal func(a, b T) bool) bool {
	// 选择较小的集合作为缓存
	var staticSet map[H][]T
	var dynamic []T
	if len(targets) < len(elems) {
		staticSet = make(map[H][]T, len(targets))
		for _, t := range targets {
			h := hashFunc(t)
			staticSet[h] = append(staticSet[h], t)
		}
		dynamic = elems
	} else {
		staticSet = make(map[H][]T, len(elems))
		for _, e := range elems {
			h := hashFunc(e)
			staticSet[h] = append(staticSet[h], e)
		}
		dynamic = targets
	}

	// 遍历动态集合
	for _, v := range dynamic {
		h := hashFunc(v)
		for _, s := range staticSet[h] {
			if equal(v, s) {
				return true
			}
		}
	}
	return false
}

// Ordered 约束，表示可排序的类型
type Ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 |
		~string
}

// SortSlice 对 Ordered 类型的切片进行排序
func SortSlice[T Ordered](slice []T, order config.SortType) {
	sort.Slice(slice, func(i, j int) bool {
		if order == config.SortAsc {
			return slice[i] < slice[j]
		}
		return slice[i] > slice[j]
	})
}

func CompareArray[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// CompareArraySorted 比较两个数组在排序后是否相同（适用于可排序类型）
func CompareArraySorted[T Ordered](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}

	// 创建副本以避免修改原数组
	aCopy := make([]T, len(a))
	bCopy := make([]T, len(b))
	copy(aCopy, a)
	copy(bCopy, b)

	// 排序
	SortSlice(aCopy, config.SortAsc)
	SortSlice(bCopy, config.SortAsc)

	// 比较排序后的数组
	return CompareArray(aCopy, bCopy)
}

// CutStrSlice2Slice 获取切片的子切片
func CutStrSlice2Slice(s []string, key string, direct int) []string {
	for idx, v := range s {
		if v == key {
			if idx+direct < len(s) {
				return s[idx+direct:]
			} else {
				return []string{} // 索引越界时返回空切片
			}
		}
	}
	return []string{}
}

// IsArray 检查值是否为数组或切片
func IsArray(v any) bool {
	if v == nil {
		return false
	}

	val := reflect.ValueOf(v)
	kind := val.Kind()

	// 检查是否是数组或切片
	return kind == reflect.Array || kind == reflect.Slice
}

func SumSlice[T Ordered](slice []T) T {
	var sum T
	for _, v := range slice {
		sum += v
	}
	return sum
}

// Number 使用 constraints 包简化类型约束
type Number interface {
	constraints.Integer | constraints.Float
}

// AvgSlice 返回原类型的平均数（浮点结果）
func AvgSlice[T Number](slice []T) T {
	if len(slice) == 0 {
		return 0
	}

	var sum T
	for _, v := range slice {
		sum += v
	}

	return sum / T(len(slice))
}

// DeleteStrarr 删除字符串切片的某一个元素
// noinspection all
func DeleteStrarr(arr []string, val string) []string {
	for index, _id := range arr {
		if _id == val {
			arr = append(arr[:index], arr[index+1:]...)
			break
		}
	}
	return arr
}

// Str2StrSlice 字符串转切片
func Str2StrSlice(values string) ([]string, error) {
	values = strings.TrimSpace(values)
	if values == "" {
		return nil, nil
	}
	var slice []string
	if err := json.Unmarshal([]byte(values), &slice); err != nil {
		return nil, fmt.Errorf("解析数据 [%s] 失败：%v", values, err)
	}
	return slice, nil
}

func StrSlice2AnySlice(inp []string) []any {
	var out []any
	for _, v := range inp {
		out = append(out, v)
	}
	return out
}

func AnySlice2StrSlice(slice []any) []string {
	var out []string
	for _, v := range slice {
		out = append(out, Any2string(v))
	}
	return out
}

// AnySlice2Str 将任意切片转成字符串
func AnySlice2Str(slice []any, _sep ...string) string {
	var builder strings.Builder
	l := len(slice)
	sep := ","
	if len(_sep) > 0 {
		sep = _sep[0]
	}
	for index, elem := range slice {
		// 使用 fmt.Sprint 将任何类型转换为字符串形式
		strElem := fmt.Sprint(elem)
		if strElem == "" {
			continue
		}
		builder.WriteString(strElem)
		// 可以选择在此处添加分隔符，如空格、逗号等
		if index < (l - 1) {
			builder.WriteString(sep)
		}
	}
	return builder.String()
}

func AnySlice2StrWithEmpty(slice []any, _sep ...string) string {
	var builder strings.Builder
	l := len(slice)
	sep := ","
	if len(_sep) > 0 {
		sep = _sep[0]
	}
	for index, elem := range slice {
		// 使用 fmt.Sprint 将任何类型转换为字符串形式
		var strElem string
		if elem == nil {
			strElem = ""
		} else {
			strElem = Any2string(elem)
		}

		builder.WriteString(strElem)
		// 可以选择在此处添加分隔符，如空格、逗号等
		if index < (l - 1) {
			builder.WriteString(sep)
		}
	}
	return builder.String()
}

// Slice2MapWithHeader 主要是将excel 或者 csv的每一行转为map，键为header，值为cell
func Slice2MapWithHeader(rows any, header []string) map[string]any {
	// 获取 rows 的反射值
	rowsValue := reflect.ValueOf(rows)
	// 检查 rows 是否为切片类型
	if rowsValue.Kind() != reflect.Slice && rowsValue.Kind() != reflect.Ptr {
		return nil
	}
	// 如果 rows 是切片的指针，则获取指向的切片
	if rowsValue.Kind() == reflect.Ptr {
		if rowsValue.IsNil() {
			return nil
		}
		rowsValue = rowsValue.Elem()
	}
	fieldLen := len(header)
	var tmp = make(map[string]any)
	//判断rows是切片，或者是切片的指针，如果是就遍历，不是就返回nil
	// 遍历 rows 切片
	for i := 0; i < rowsValue.Len(); i++ {
		if i >= fieldLen {
			continue
		}
		tmp[header[i]] = rowsValue.Index(i).Interface()
	}
	return tmp
}
