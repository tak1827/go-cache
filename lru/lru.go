package lru

import (
	"container/list"
	"fmt"
	"log"
	"sync"
	"time"
)

const (
	MIN_TTL = 60 // 1min
)

type LRUCache struct {
	sync.Mutex

	size     int
	ttl      int64 // sec
	elements map[string]*list.Element
	list     *list.List
}

type entry struct {
	key       string
	value     interface{}
	valueWith valueWithTTL
}

type valueWithTTL struct {
	ttl   int64
	value interface{}
}

func NewCache(size int, ttl int) LRUCache {
	if size == 0 {
		panic(fmt.Sprintf("cache size should be positive"))
	}

	// just output waring log
	if 0 < ttl && ttl <= MIN_TTL {
		log.Printf("too smal ttl setting. recomend to be equal or longer than %d sec\n", MIN_TTL)
	}

	return LRUCache{
		size:     size,
		ttl:      int64(ttl),
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

	if c.hasTTL() {
		now := time.Now().Unix()
		// if expired, delete
		if en.valueWith.ttl < now {
			c.list.Remove(elm)
			delete(c.elements, elm.Value.(*entry).key)
			ok = false
			return
		}
		val = en.valueWith.value
		// update ttl
		elm.Value.(*entry).valueWith.ttl = now + c.ttl
	} else {
		val = en.value
	}

	c.list.MoveToFront(elm)
	return
}

func (c *LRUCache) Add(key string, value interface{}) (evicted bool) {
	c.Lock()
	defer c.Unlock()

	if elm, ok := c.elements[key]; ok {
		// everytime overwrite the value
		c.list.Remove(elm)
	} else {
		for len(c.elements) >= c.size {
			lastElm := c.list.Back()
			c.list.Remove(lastElm)
			lastEntry := lastElm.Value.(*entry)
			delete(c.elements, lastEntry.key)
			evicted = true
		}
	}

	var e entry
	if c.hasTTL() {
		e = entry{
			key:       key,
			valueWith: valueWithTTL{time.Now().Unix() + c.ttl, value},
		}
	} else {
		e = entry{
			key:   key,
			value: value,
		}
	}

	elm := c.list.PushFront(&e)
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
	delete(c.elements, key)
	return
}

func (c *LRUCache) Clear() {
	c.Lock()
	defer c.Unlock()

	for k := range c.elements {
		c.list.Remove(c.elements[k])
		delete(c.elements, k)
	}

	return
}

func (c *LRUCache) Contains(key string) bool {
	c.Lock()
	defer c.Unlock()

	if !c.hasTTL() {
		_, ok := c.elements[key]
		return ok
	}

	// check ttl
	elm, ok := c.elements[key]
	if !ok {
		return false
	}
	return time.Now().Unix() <= elm.Value.(*entry).valueWith.ttl
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

func (c *LRUCache) hasTTL() bool {
	return 0 < c.ttl
}
