package lru

import (
	"container/list"
	"sync"
)

type Cache struct {
	cap int
	l   *list.List
	m   map[interface{}]*list.Element
	mu  sync.Mutex
}

type element struct {
	key   interface{}
	value interface{}
}

func New(size int) *Cache {
	return &Cache{
		cap: size,
		l:   list.New(),
		m:   make(map[interface{}]*list.Element, size),
		mu:  sync.Mutex{},
	}
}

func (c *Cache) Load(key interface{}) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if e, ok := c.m[key]; ok {
		c.l.MoveToBack(e)
		return e.Value.(*element).value, true
	}
	return nil, false
}

func (c *Cache) Store(key, value interface{}) interface{} {
	c.mu.Lock()
	defer c.mu.Unlock()
	if e, ok := c.m[key]; ok {
		c.l.Remove(e)
	}
	c.m[key] = c.l.PushBack(&element{key: key, value: value})
	if c.cap < c.l.Len() {
		e := c.l.Front()
		c.l.Remove(e)
		delete(c.m, e.Value.(*element).key)
	}
	return value
}

func (c *Cache) Delete(key interface{}) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if e, ok := c.m[key]; ok {
		c.l.Remove(e)
		delete(c.m, key)
		return e.Value.(*element).value, ok
	}
	return nil, false
}

func (c *Cache) Len() int {
	return c.l.Len()
}

func (c *Cache) Cap() int {
	return c.cap
}
