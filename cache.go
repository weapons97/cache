// Package cache
/*=============================================================================
#       Author: peng.wei
#        Email: weapons97@gmail.com
#      Version: 0.0.1
#   LastChange: 20211214
#      History:
=============================================================================*/
package cache

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

var (
	// Forever 永远不超时
	Forever = time.Hour * 24 * 365 * 100
	// defaultTTL 默认超时时间
	defaultTTL = Forever
)

// CacheI 接口, 一个结构实现了CacheI 接口才可以被CacheManager 使用
type CacheI interface {
	Clean()
	Name() string
}

// InterfaceCache 是cache的接口
type InterfaceCache[K comparable, V any] interface {
	Set(req K, values V)
	Remove(key ...K)
	Has(items ...K) bool
	Size() int
	Clear()
	IsEmpty() bool
	Range(fn func(s string, k K) bool)
	List() ([]K, []V)
	ListKey() []K
	ListValue() []V
	Merge(s *Cache[K, V])
}

// Remove 根据key 删除 cache
func (c *Cache[K, V]) Remove(k ...K) {
	for i := range k {
		c.Del(k[i])
	}
}

// Has 是否包含
func (c *Cache[K, V]) Has(k ...K) bool {
	for i := range k {
		_, ok := c.Get(k[i])
		if !ok {
			return false
		}
	}
	return true
}

// HasAny 是否包含任何一个键
func (c *Cache[K, V]) HasAny(k ...K) bool {
	for i := range k {
		_, ok := c.Get(k[i])
		if ok {
			return true
		}
	}
	return false
}

// Size 返回cache 长度
func (c *Cache[K, V]) Size() int {
	return c.Len()
}

// Clear 清空cache
func (c *Cache[K, V]) Clear() {
	c.smap = &sync.Map{}
}

// IsEmpty 是否为空
func (c *Cache[K, V]) IsEmpty() bool {
	return c.Len() == 0
}

// Merge 合并cache
func (c *Cache[K, V]) Merge(s *Cache[K, V]) {
	s.Range(func(k K, v V) bool {
		c.Set(k, v)
		return true
	})
}

// Option cache 的选项
type Option[K comparable, V any] func(*Cache[K, V])

// WithTTL 设置超时
func WithTTL[K comparable, V any](ttl time.Duration) Option[K, V] {
	return func(cache *Cache[K, V]) {
		cache.ttl = ttl
	}
}

// WithName 设置cache名称
func WithName[K comparable, V any](name string) Option[K, V] {
	return func(cache *Cache[K, V]) {
		cache.name = name
	}
}

// WithNoManager 设置cache不受 cacheManager管理
func WithNoManager[K comparable, V any]() Option[K, V] {
	return func(cache *Cache[K, V]) {
		cache.noManager = true
	}
}

// init 根据opts设置cache
func (c *Cache[K, V]) init(opts ...Option[K, V]) {
	c.ttl = defaultTTL
	c.smap = &sync.Map{}
	c.name = uuid.NewString()
	for i := range opts {
		opts[i](c)
	}
	if c.ttl == defaultTTL {
		c.noManager = true
	}
}

// NewCache 创建新cache
func NewCache[K comparable, V any](opts ...Option[K, V]) *Cache[K, V] {
	res := Cache[K, V]{}
	res.opts = opts
	res.init(opts...)
	if res.noManager {
		return &res
	}
	cm := CacheManagerFactory()
	cm.RegisterCache(&res)
	return &res
}

// NewCacheInits 创建新cache
func NewCacheInits[K comparable, V any](inits map[K]V, opts ...Option[K, V]) *Cache[K, V] {
	res := NewCache(opts...)
	for k, v := range inits {
		res.Set(k, v)
	}
	return res
}

// Options 返回cache的选项
func (c *Cache[K, V]) Options() []Option[K, V] {
	return c.opts
}

// Cache 是一个带超时的缓存, 超时的元素会获取不到并删除(默认情况下)
type Cache[K comparable, V any] struct {
	ttl       time.Duration
	smap      *sync.Map
	name      string
	noManager bool
	opts      []Option[K, V]
}

// Name return name of cache
func (c *Cache[K, V]) Name() string {
	return c.name
}

// Set 设置k，v
func (c *Cache[K, V]) Set(req K, values V) {
	if c == nil || c.smap == nil {
		c = NewCache[K, V]()
	}
	c.smap.Store(req, c.wrapTTL(values))
}

// Del 根据key 删除 cache
func (c *Cache[K, V]) Del(k K) {
	c.smap.Delete(k)
}

// Range 遍历cache
func (c *Cache[K, V]) Range(fn func(k K, v V) bool) {
	c.smap.Range(func(k, v any) bool {
		xv, ok := c.unWrapTTL(v)
		if ok {
			vk := k.(K)
			vv := xv.(V)
			return fn(vk, vv)
		}
		return true
	})
}

// List func list k and list v
func (c *Cache[K, V]) List() ([]K, []V) {
	ks := make([]K, 0)
	vs := make([]V, 0)
	c.Range(func(k K, v V) bool {
		ks = append(ks, k)
		vs = append(vs, v)
		return true
	})
	return ks, vs
}

// ListKey func list k
func (c *Cache[K, V]) ListKey() []K {
	ks, _ := c.List()
	return ks
}

// ListValue func list v
func (c *Cache[K, V]) ListValue() []V {
	_, vs := c.List()
	return vs
}

// Clean 会被cache manager 定期调用删除过期的元素
func (c *Cache[K, V]) Clean() {
	c.smap.Range(func(k, v any) bool {
		if _, ok := c.unWrapTTL(v); !ok {
			vk := k.(K)
			c.Del(vk)
		}
		return true
	})
}

// Get 根据key获得value 超时或者空第二个返回值为false，否则返回true
func (c *Cache[K, V]) Get(req K) (V, bool) {
	zeroV := new(V)
	if c == nil || c.smap == nil {
		return *zeroV, false
	}
	wp, ok := c.smap.Load(req)
	if !ok {
		return *zeroV, false
	}
	v, ok := c.unWrapTTL(wp)
	if !ok {
		return *zeroV, false
	}
	vv := v.(V)
	return vv, true
}

// Len 返回cache 长度
func (c *Cache[K, V]) Len() int {
	i := 0
	c.Range(func(k K, v V) bool {
		i++
		return true
	})
	return i
}

// wrap 是cache的元素
type wrap struct {
	timeout time.Time
	v       any
}

// TimeValue 是一个带有时间的值
type TimeValue interface {
	Time() time.Time
}

// TimeoutValue 是一个带有超时时间的值
type TimeoutValue interface {
	Timeout() time.Time
}

func (c *Cache[K, V]) wrapTTL(v any) *wrap {

	switch tv := v.(type) {
	case TimeoutValue:
		return &wrap{
			timeout: tv.Timeout(),
			v:       v,
		}
	case TimeValue:
		return &wrap{
			timeout: tv.Time().Add(c.ttl),
			v:       v,
		}
	default:
		return &wrap{
			timeout: time.Now().Add(c.ttl),
			v:       v,
		}
	}
}

func (c *Cache[K, V]) unWrapTTL(v any) (any, bool) {
	wp, ok := v.(*wrap)
	if !ok {
		return nil, false
	}
	if time.Now().After(wp.timeout) {
		return nil, false
	}
	return wp.v, true
}
