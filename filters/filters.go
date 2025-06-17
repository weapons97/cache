package filters

import (
	"cmp"
	"reflect"
	"strings"
)

// NoSpace 检查字符串是否不为空（去除空格后）
func NoSpace(s string) bool {
	return strings.TrimSpace(s) != ""
}

// Range 遍历切片，当函数返回false时停止
func Range[T any](objs []T, fn func(obj T) bool) {
	for _, obj := range objs {
		if !fn(obj) {
			return
		}
	}
}

// Filter 过滤切片，返回满足条件的元素
func Filter[T any](objs []T, filter func(obj T) bool) []T {
	if len(objs) == 0 {
		return nil
	}

	res := make([]T, 0, len(objs))
	for _, obj := range objs {
		if filter(obj) {
			res = append(res, obj)
		}
	}
	return res
}

// Map 对切片进行映射转换
func Map[T any, K any](objs []T, mapper func(obj T) ([]K, bool)) []K {
	if len(objs) == 0 {
		return nil
	}

	res := make([]K, 0, len(objs))
	for _, obj := range objs {
		if others, ok := mapper(obj); ok {
			res = append(res, others...)
		}
	}
	return res
}

// First 返回切片的第一个元素
func First[T any](objs []T) (T, bool) {
	if len(objs) > 0 {
		return objs[0], true
	}
	var zero T
	return zero, false
}

// isZero 检查值是否为零值（内部函数）
func isZero(v reflect.Value) bool {
	// 检查无效的reflect.Value
	if !v.IsValid() {
		return true
	}

	switch v.Kind() {
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Complex64, reflect.Complex128:
		return v.Complex() == 0
	case reflect.String:
		return v.String() == ""
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Pointer, reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		return v.IsNil()
	case reflect.Array:
		for i := 0; i < v.Len(); i++ {
			if !isZero(v.Index(i)) {
				return false
			}
		}
		return true
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			if !isZero(v.Field(i)) {
				return false
			}
		}
		return true
	default:
		// 对于其他类型，使用零值比较
		z := reflect.Zero(v.Type())
		return v.Interface() == z.Interface()
	}
}

// OR 返回第一个非零值
func OR[T any](vs ...T) T {
	for _, v := range vs {
		if !IsZero(v) {
			return v
		}
	}
	var zero T
	return zero
}

// IsZero 检查值是否为零值
func IsZero[T any](v T) bool {
	return isZero(reflect.ValueOf(v))
}

// FilterMap 过滤map
func FilterMap[K comparable, V any](m map[K]V, f func(K, V) bool) map[K]V {
	if len(m) == 0 {
		return make(map[K]V)
	}

	ret := make(map[K]V, len(m))
	for k, v := range m {
		if f(k, v) {
			ret[k] = v
		}
	}
	return ret
}

// addOrUpdateSliceSingle 添加或更新切片中的单个元素
func addOrUpdateSliceSingle[T any](slice []T, item T) []T {
	itemValue := reflect.ValueOf(item)
	itemType := reflect.TypeOf(item)

	for i, v := range slice {
		vValue := reflect.ValueOf(v)

		// 类型不匹配，跳过
		if vValue.Type() != itemType {
			continue
		}

		// 使用反射进行类型安全的比较
		switch vValue.Kind() {
		case reflect.String:
			if vValue.String() == itemValue.String() {
				slice[i] = item
				return slice
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if vValue.Int() == itemValue.Int() {
				slice[i] = item
				return slice
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			if vValue.Uint() == itemValue.Uint() {
				slice[i] = item
				return slice
			}
		case reflect.Float32, reflect.Float64:
			if vValue.Float() == itemValue.Float() {
				slice[i] = item
				return slice
			}
		case reflect.Bool:
			if vValue.Bool() == itemValue.Bool() {
				slice[i] = item
				return slice
			}
		default:
			// 对于复杂类型，使用DeepEqual
			if reflect.DeepEqual(v, item) {
				slice[i] = item
				return slice
			}
		}

		// 对于结构体，检查Name字段（如果存在）
		if vValue.Kind() == reflect.Struct && itemValue.Kind() == reflect.Struct {
			vName := vValue.FieldByName("Name")
			itemName := itemValue.FieldByName("Name")

			if vName.IsValid() && itemName.IsValid() &&
				vName.String() == itemName.String() {
				slice[i] = item
				return slice
			}
		}
	}
	return append(slice, item)
}

// AddOrUpdateSlice 添加或更新切片中的元素
func AddOrUpdateSlice[T any](slice []T, items ...T) []T {
	if len(items) == 0 {
		return slice
	}
	if len(slice) == 0 {
		return items
	}

	res := make([]T, len(slice))
	copy(res, slice)

	for _, item := range items {
		res = addOrUpdateSliceSingle(res, item)
	}
	return res
}

// GetDefault 获取默认值，如果输入为零值则返回默认值
func GetDefault[T any](s, d T) T {
	if IsZero(s) {
		return d
	}
	return s
}

// Max 返回最大值
func Max[T cmp.Ordered](ts ...T) T {
	if len(ts) == 0 {
		var zero T
		return zero
	}

	m := ts[0]
	for i := 1; i < len(ts); i++ {
		if ts[i] > m {
			m = ts[i]
		}
	}
	return m
}

// Min 返回最小值
func Min[T cmp.Ordered](ts ...T) T {
	if len(ts) == 0 {
		var zero T
		return zero
	}

	m := ts[0]
	for i := 1; i < len(ts); i++ {
		if ts[i] < m {
			m = ts[i]
		}
	}
	return m
}

// Sum 计算总和
func Sum[T cmp.Ordered](ts ...T) T {
	if len(ts) == 0 {
		var zero T
		return zero
	}

	s := ts[0]
	for i := 1; i < len(ts); i++ {
		s += ts[i]
	}
	return s
}

// Contains 检查切片是否包含指定元素
func Contains[T comparable](slice []T, item T) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

// Unique 返回去重后的切片
func Unique[T comparable](slice []T) []T {
	if len(slice) == 0 {
		return make([]T, 0)
	}

	seen := make(map[T]struct{}, len(slice))
	result := make([]T, 0, len(slice))

	for _, item := range slice {
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

// Reverse 反转切片
func Reverse[T any](slice []T) []T {
	if len(slice) <= 1 {
		return slice
	}

	result := make([]T, len(slice))
	for i, j := 0, len(slice)-1; i < len(slice); i, j = i+1, j-1 {
		result[i] = slice[j]
	}
	return result
}

// Chunk 将切片分割成指定大小的块
func Chunk[T any](slice []T, size int) [][]T {
	if size <= 0 || len(slice) == 0 {
		return make([][]T, 0)
	}

	chunks := make([][]T, 0, (len(slice)+size-1)/size)
	for i := 0; i < len(slice); i += size {
		end := i + size
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slice[i:end])
	}
	return chunks
}
