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

// Option cache 的选项
type Option func(*Cache)

// WithTTL 设置超时
func WithTTL(ttl time.Duration) Option {
	return func(cache *Cache) {
		cache.ttl = ttl
	}
}

// WithName 设置cache名称
func WithName(name string) Option {
	return func(cache *Cache) {
		cache.name = name
	}
}

// WithNoManager 设置cache不受 cacheManager管理
func WithNoManager() Option {
	return func(cache *Cache) {
		cache.noManager = true
	}
}

// init 根据opts设置cache
func (c *Cache) init(opts ...Option) {
	c.ttl = defaultTTL
	c.smap = &sync.Map{}
	c.name = uuid.NewString()
	for i := range opts {
		opts[i](c)
	}
}

// NewCache 创建新cache
func NewCache(opts ...Option) *Cache {
	res := new(Cache)
	res.init(opts...)
	if res.noManager {
		return res
	}
	cm := CacheManagerFactory()
	cm.RegisterCache(res)
	return res
}

// Cache 是一个带超时的缓存, 超时的元素会获取不到并删除(默认情况下)
type Cache struct {
	ttl       time.Duration
	smap      *sync.Map
	name      string
	noManager bool
}

// Set 设置k，v
func (c *Cache) Set(req interface{}, values interface{}) {
	if c == nil || c.smap == nil {
		c = NewCache()
	}
	c.smap.Store(req, c.wrapTTL(values))
}

// Del 根据key 删除 cache
func (c *Cache) Del(k interface{}) {
	c.smap.Delete(k)
}

// Range 遍历cache
func (c *Cache) Range(fn func(k, v interface{}) bool) {
	c.smap.Range(func(k, v interface{}) bool {
		xv, ok := c.unWrapTTL(v)
		if ok {
			return fn(k, xv)
		}
		return true
	})
}

// Clean 会被cache manager 定期调用删除过期的元素
func (c *Cache) Clean() {
	c.smap.Range(func(k, v interface{}) bool {
		if _, ok := c.unWrapTTL(v); !ok {
			c.Del(k)
		}
		return true
	})
}

// Get 根据key获得value 超时或者空第二个返回值为false，否则返回true
func (c *Cache) Get(req interface{}) (interface{}, bool) {
	if c == nil || c.smap == nil {
		return nil, false
	}
	wp, ok := c.smap.Load(req)
	if !ok {
		return nil, false
	}
	v, ok := c.unWrapTTL(wp)
	if !ok {
		return nil, false
	}
	return v, true
}

// wrap 是cache的元素
type wrap struct {
	timeout time.Time
	v       interface{}
}

func (c *Cache) wrapTTL(v interface{}) *wrap {
	return &wrap{
		timeout: time.Now().Add(c.ttl),
		v:       v,
	}
}
func (c *Cache) unWrapTTL(v interface{}) (interface{}, bool) {
	wp, ok := v.(*wrap)
	if !ok {
		return nil, false
	}
	if time.Now().After(wp.timeout) {
		return nil, false
	}
	return wp.v, true
}
