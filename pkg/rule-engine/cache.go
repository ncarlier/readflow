package ruleengine

import (
	"container/list"
	"sync"
)

// Cache is a cache to strore rules processors
// LRU strategy: evicts least recently used items
type Cache struct {
	capacity int
	items    map[uint]*CacheItem
	list     *list.List
	lock     *sync.RWMutex
}

// CacheItem is an item of the cache
type CacheItem struct {
	key         uint
	value       *ProcessorPipeline
	listElement *list.Element
}

// NewRuleEngineCache creates a new cache for the rule engine
func NewRuleEngineCache(capacity int) *Cache {
	return &Cache{
		capacity: capacity,
		items:    make(map[uint]*CacheItem, capacity),
		list:     list.New(),
		lock:     new(sync.RWMutex),
	}
}

// Get retrieve a list of rule processor from the cache
func (c *Cache) Get(key uint) *ProcessorPipeline {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if item, exists := c.items[key]; exists {
		c.promote(item)
		return item.value
	}
	return nil
}

// Set put a list of rule processor into the cache
func (c *Cache) Set(key uint, value *ProcessorPipeline) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.capacity <= c.list.Len() {
		c.prune()
	}

	if item, exists := c.items[key]; exists {
		item.value = value
		c.promote(item)
	} else {
		item = &CacheItem{key: key, value: value}
		item.listElement = c.list.PushFront(item)
		c.items[key] = item
	}
}

// Evict remove a list of rule processor from the cache
func (c *Cache) Evict(key uint) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if item, exists := c.items[key]; exists {
		c.list.Remove(item.listElement)
		delete(c.items, item.key)
	}
}

func (c *Cache) promote(item *CacheItem) {
	c.list.MoveToFront(item.listElement)
}

func (c *Cache) prune() {
	for i := 0; i < 50; i++ {
		tail := c.list.Back()
		if tail == nil {
			return
		}
		item := c.list.Remove(tail).(*CacheItem)
		delete(c.items, item.key)
	}
}
