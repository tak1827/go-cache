package cache

import (
	"container/list"
	"sync"
)

type LRUCache struct {
	sync.Mutex

	size     int
	elements map[string]*list.Element
	list     *list.List
}

type entry struct {
	key   string
	value interface{}
}

func NewCache(size int) LRUCache {
	return LRUCache{
		size:     size,
		elements: make(map[string]*list.Element, size),
		list:     list.New(),
	}
}

func (c *LRUCache) Get(key string) (val interface{}, ok bool) {
	c.Lock()
	defer c.Unlock()

	elm, ok := c.elements[key]
	if !ok {
		return
	}
	en := elm.Value.(*entry)
	val = en.value
	return
}

func (c *LRUCache) Add(key string, value interface{}) (evicted bool) {
	c.Lock()
	defer c.Unlock()

	if elm, ok := c.elements[key]; ok {
		c.list.MoveToFront(elm)
		return
	}

	for len(c.elements) >= c.size {
		lastElm := c.list.Back()
		c.list.Remove(lastElm)
		lastEntry := lastElm.Value.(*entry)
		delete(c.elements, lastEntry.key)
		evicted = true
	}

	elm := c.list.PushFront(&entry{key, value})
	c.elements[key] = elm
	return
}

func (c *LRUCache) Remove(key string) (present bool) {
	c.Lock()
	defer c.Unlock()

	elm, present := c.elements[key]
	if !present {
		return
	}
	c.list.Remove(elm)
	en := elm.Value.(*entry)
	delete(c.elements, en.key)
	return
}

func (c *LRUCache) Contains(key string) bool {
	c.Lock()
	defer c.Unlock()

	_, ok := c.elements[key]
	return ok
}

func (c *LRUCache) Len() int {
	c.Lock()
	defer c.Unlock()

	return len(c.elements)
}

func (c *LRUCache) Cap() int {
	c.Lock()
	defer c.Unlock()

	return c.size
}
