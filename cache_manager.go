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

	"github.com/weapons97/cache/wait"
)

var (
	// defaultManager 默认的cachemanager
	defaultManager = newCacheManager()
)

// CacheManager 持有所有的cache 可以定时执行cache 的tasks
type CacheManager interface {
	Tasks()                  // CacheManager 的tasks 会定期执行
	RegisterCache(c *Cache)  // 注册cache
	Interval() time.Duration // 返回执行tasks的间隔
}

// RegisterCache 向CacheManager注册cache
func RegisterCache(c *Cache) {
	cm := CacheManagerFactory()
	cm.RegisterCache(c)
}

func background() {
	cm := CacheManagerFactory()
	go wait.Forever(cm.Tasks, cm.Interval(), true)
}

func newCacheManager() CacheManager {
	cm := &cacheManager{caches: make(map[string]*Cache, 10)}
	cm.cleanInterval = time.Second * 60
	return cm
}

func init() {
	background()
}

// cacheManager 持有所有的cache 可以定时执行cache 的tasks
type cacheManager struct {
	caches        map[string]*Cache
	cachesl       sync.Mutex
	cleanInterval time.Duration
}

// Interval 返回执行tasks的间隔
func (cm *cacheManager) Interval() time.Duration {
	return cm.cleanInterval
}

// Tasks cacheManager 的tasks 会定期执行
func (cm *cacheManager) Tasks() {
	cm.clean()
}

func (cm *cacheManager) clean() {
	cm.cachesl.Lock()
	defer cm.cachesl.Unlock()
	for _, c := range cm.caches {
		c.Clean()
	}
}

// RegisterCache 注册cache
func (cm *cacheManager) RegisterCache(c *Cache) {
	cm.cachesl.Lock()
	defer cm.cachesl.Unlock()
	cm.caches[c.name] = c
}

// CacheManagerFactory 返回 CacheManager 的方法
var CacheManagerFactory = func() CacheManager {
	return defaultManager
}
