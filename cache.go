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
	Forever = time.Hour * 24 * 365 * 10
	// defaultTTL 默认超时时间
	defaultTTL = Forever
)

// CacheI 接口, 一个结构实现了CacheI 接口才可以被CacheManager 使用
type CacheI interface {
	Clean()
	Name() string
}

// Option cache 的选项
type Option[K any, V any] func(*Cache[K, V])

// WithTTL 设置超时
func WithTTL[K any, V any](ttl time.Duration) Option[K, V] {
	return func(cache *Cache[K, V]) {
		cache.ttl = ttl
	}
}

// WithName 设置cache名称
func WithName[K any, V any](name string) Option[K, V] {
	return func(cache *Cache[K, V]) {
		cache.name = name
	}
}

// WithNoManager 设置cache不受 cacheManager管理
func WithNoManager[K any, V any]() Option[K, V] {
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
}

// NewCache 创建新cache
func NewCache[K any, V any](opts ...Option[K, V]) *Cache[K, V] {
	res := Cache[K, V]{}
	res.init(opts...)
	if res.noManager {
		return &res
	}
	cm := CacheManagerFactory()
	cm.RegisterCache(&res)
	return &res
}

// Cache 是一个带超时的缓存, 超时的元素会获取不到并删除(默认情况下)
type Cache[K any, V any] struct {
	ttl       time.Duration
	smap      *sync.Map
	name      string
	noManager bool
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

// ListSet func list k and list v with set
func (c *Cache[K, V]) ListSet() (*Set[K], *Set[V]) {
	ks := make([]K, 0)
	vs := make([]V, 0)
	c.Range(func(k K, v V) bool {
		ks = append(ks, k)
		vs = append(vs, v)
		return true
	})

	return NewSetInits(ks), NewSetInits(vs)
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

func (c *Cache[K, V]) wrapTTL(v any) *wrap {
	return &wrap{
		timeout: time.Now().Add(c.ttl),
		v:       v,
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
