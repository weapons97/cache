package filters

import (
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
