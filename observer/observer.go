package observer

import (
	"github.com/weapons97/cache/filters"
)

// Observer 观察者模式，封装切片并提供链式操作
type Observer[T any] struct {
	data []T
}

// New 创建新的Observer实例
func New[T any](data []T) *Observer[T] {
	return &Observer[T]{data: data}
}

// From 从切片创建Observer
func From[T any](data []T) *Observer[T] {
	return New(data)
}

// Data 获取底层数据
func (o *Observer[T]) Data() []T {
	return o.data
}

// Len 返回数据长度
func (o *Observer[T]) Len() int {
	return len(o.data)
}

// IsEmpty 检查是否为空
func (o *Observer[T]) IsEmpty() bool {
	return len(o.data) == 0
}

// Range 对切片执行函数，返回第一个非零值
func (o *Observer[T]) Range(fn func(obj T) bool) {
	filters.Range(o.data, fn)
}

// Filter 过滤切片，返回满足条件的元素
func (o *Observer[T]) Filter(filter func(obj T) bool) *Observer[T] {
	return New(filters.Filter(o.data, filter))
}

// Map 对切片进行映射转换
func (o *Observer[T]) Map(mapper func(obj T) ([]T, bool)) *Observer[T] {
	return New(filters.Map(o.data, mapper))
}

// First 返回切片的第一个元素
func (o *Observer[T]) First() (T, bool) {
	return filters.First(o.data)
}

// OR 返回第一个非零值
func (o *Observer[T]) OR() T {
	return filters.OR(o.data...)
}

// IsZero 检查值是否为零值
func (o *Observer[T]) IsZero() bool {
	return filters.IsZero(o.data)
}

// AddOrUpdate 添加或更新切片中的元素
func (o *Observer[T]) AddOrUpdate(items ...T) *Observer[T] {
	return New(filters.AddOrUpdateSlice(o.data, items...))
}

// GetDefault 获取默认值，如果输入为零值则返回默认值
func (o *Observer[T]) GetDefault(d T) T {
	if len(o.data) == 0 {
		return d
	}
	return filters.GetDefault(o.data[0], d)
}
