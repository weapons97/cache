package filters

import (
	"cmp"
	"reflect"
	"strings"
)

// NoSpace is filter func for strings
func NoSpace(s string) bool {
	return strings.TrimSpace(s) != ""
}

// Range run func in slice
func Range[T any](objs []T, fn func(obj T) bool) {
	for i := range objs {
		if !fn(objs[i]) {
			return
		}
	}
}

// Filter filter one slice
func Filter[T any](objs []T, filter func(obj T) bool) []T {
	res := make([]T, 0, len(objs))
	for i := range objs {
		ok := filter(objs[i])
		if ok {
			res = append(res, objs[i])
		}
	}
	return res
}

// Map one slice
func Map[T any, K any](objs []T, mapper func(obj T) ([]K, bool)) []K {
	res := make([]K, 0, len(objs))
	for i := range objs {
		others, ok := mapper(objs[i])
		if ok {
			res = append(res, others...)
		}
	}
	return res
}

// First make return first for slice
func First[T any](objs []T) (T, bool) {
	if len(objs) > 0 {
		return objs[0], true
	}
	return *new(T), false
}

func isZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Func, reflect.Map, reflect.Slice:
		return v.IsNil()
	case reflect.Array:
		z := true
		for i := 0; i < v.Len(); i++ {
			z = z && isZero(v.Index(i))
		}
		return z
	case reflect.Struct:
		z := true
		for i := 0; i < v.NumField(); i++ {
			z = z && isZero(v.Field(i))
		}
		return z
	}
	// Compare other types directly:
	z := reflect.Zero(v.Type())
	return v.Interface() == z.Interface()
}

// OR return values which not zero value
func OR[T any](vs ...T) T {
	for i := range vs {
		if !IsZero(vs[i]) {
			return vs[i]
		}
	}
	return *new(T)
}

// IsZero true when arg is zero
func IsZero[T any](v T) bool {
	return isZero(reflect.ValueOf(v))
}

// FilterMap filter one map to another map
func FilterMap[K comparable, V any](m map[K]V, f func(K, V) bool) map[K]V {
	ret := make(map[K]V)
	for k, v := range m {
		if f(k, v) {
			ret[k] = v
		}
	}
	return ret
}

func addOrUpdateSliceSingle[T any](slice []T, item T) []T {
	for i, v := range slice {
		switch any(v).(type) {
		case string, int, int32, int64, float32, float64, bool:
			if reflect.DeepEqual(v, item) {
				slice[i] = item
				return slice
			}
			continue
		}

		vn := reflect.ValueOf(v).FieldByName("Name")
		if vn.String() == reflect.ValueOf(item).FieldByName("Name").String() {
			slice[i] = item
			return slice
		}
	}
	return append(slice, item)
}

// AddOrUpdateSlice add or update slice
func AddOrUpdateSlice[T any](slice []T, items ...T) []T {
	if len(slice) == 0 {
		return items
	}
	res := slice
	for _, v := range items {
		res = addOrUpdateSliceSingle(res, v)
	}
	return res
}

// GetDefault godoc
func GetDefault[T any](s, d T) T {
	sv := reflect.ValueOf(s)
	if sv == reflect.ValueOf(nil) {
		return d
	}
	if sv.IsZero() {
		return d
	}
	if sv.Comparable() {
		st := reflect.TypeOf(s)
		if sv.Interface() == reflect.Zero(st).Interface() {
			return d
		}
		return s
	}
	k := sv.Kind()
	switch k {
	case reflect.Chan, reflect.Func, reflect.Map,
		reflect.Pointer, reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		if sv.IsNil() {
			return d
		}
	default:
		return s
	}
	return s
}

func Max[T cmp.Ordered](ts ...T) T {
	if len(ts) == 0 {
		return *new(T)
	}
	m := ts[0]

	for i := 1; i < len(ts); i++ {
		m = max(m, ts[i])
	}
	return m
}

func Min[T cmp.Ordered](ts ...T) T {
	if len(ts) == 0 {
		return *new(T)
	}
	m := ts[0]
	for i := 1; i < len(ts); i++ {
		m = min(m, ts[i])
	}
	return m
}

func Sum[T cmp.Ordered](ts ...T) T {
	if len(ts) == 0 {
		return *new(T)
	}
	s := ts[0]
	for i := 1; i < len(ts); i++ {
		s = s + ts[i]
	}
	return s
}
