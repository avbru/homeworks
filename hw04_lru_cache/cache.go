package hw04_lru_cache //nolint:golint,stylecheck
import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	sync.Mutex
	capacity int
	queue    List
	items    map[Key]*item
}

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(cap int) Cache {
	return &lruCache{
		capacity: cap,
		queue:    NewList(),
		items:    make(map[Key]*item, cap),
	}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.Lock()
	defer c.Unlock()
	item, ok := c.items[key]

	if ok {
		c.queue.MoveToFront(item)
		item.value = cacheItem{key, value}
		return true
	}

	if c.queue.Len() >= c.capacity {
		backItem := c.queue.Back()
		c.queue.Remove(backItem)
		delete(c.items, backItem.value.(cacheItem).key)
	}

	c.items[key] = c.queue.PushFront(cacheItem{key, value})
	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.Lock()
	defer c.Unlock()
	item, ok := c.items[key]
	if ok {
		c.queue.MoveToFront(item)
		return item.value.(cacheItem).value, true
	}
	return nil, false
}

func (c *lruCache) Clear() {
	c.Lock()
	defer c.Unlock()

	c.queue = NewList()
	c.items = make(map[Key]*item, c.capacity)
}
