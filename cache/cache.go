package cache

import (
	"sync"
	"time"
)

type Cache interface {
	ContainsOrAdd(val uint64) bool
}

type nodeValue struct {
	Value     uint64
	Timestamp uint32
}

type node struct {
	Data nodeValue
	Next *node
}

type SimpleCache struct {
	head   *node
	tail   *node
	values map[uint64]struct{}
	ttl    uint32
	lock   sync.Mutex
}

func (c *SimpleCache) ContainsOrAdd(val uint64) bool {
	c.lock.Lock()
	defer c.lock.Unlock()

	if _, ok := c.values[val]; ok {
		return true
	}

	c.values[val] = struct{}{}

	n := &node{
		Data: nodeValue{
			Value:     val,
			Timestamp: uint32(time.Now().Unix()),
		},
		Next: nil,
	}

	if c.tail == nil {
		c.head = n
		c.tail = n
	} else {
		c.tail.Next = n
		c.tail = n
	}

	return false
}

func (c *SimpleCache) Clean() {
	currTime := uint32(time.Now().Unix())

	c.lock.Lock()

	if c.head == nil {
		c.lock.Unlock()
		return
	}

	c.lock.Unlock()

	for c.head != nil {
		if c.head.Data.Timestamp+c.ttl > currTime {
			break
		}

		c.lock.Lock()

		delete(c.values, c.head.Data.Value)

		c.head = c.head.Next

		if c.head == nil {
			c.tail = nil
			c.values = make(map[uint64]struct{})
		}

		c.lock.Unlock()
	}
}

func (c *SimpleCache) run() {
	ticker := time.NewTicker(time.Duration(c.ttl) * time.Second)

	for range ticker.C {
		c.Clean()
	}
}

func NewSimpleCache(ttl uint32) *SimpleCache {
	c := &SimpleCache{
		head:   nil,
		tail:   nil,
		values: make(map[uint64]struct{}),
		ttl:    ttl,
		lock:   sync.Mutex{},
	}

	go c.run()

	return c
}
