package cache

import (
	"crypto/md5"
	"sync"
	"time"
)

// This package provides a very simple cache. It's designed to hide the values of the keys because
// they will be used for storing authentication information, so keys are hashed before being used.
//
// 	c := Cache[int]()
// 	k := c.Key("secret number")
// 	c.Put(k, 42)
// 	...
// 	if v, ok := c.Get(k); ok {
//		...
// 	}

type Key [16]byte

type Entry[Value any] struct {
	value  *Value
	expiry int64
}

type Cache[Value any] struct {
	entries    *sync.Map
	last_flush int64
}

func New[Value any]() *Cache[Value] {
	return &Cache[Value]{
		entries: &sync.Map{},
	}
}

func (c *Cache[V]) Key(k string) Key {
	return md5.Sum([]byte(k))
}

func (c *Cache[Value]) Get(k Key) (*Value, bool) {
	if value, ok := c.entries.Load(k); ok {
		if entry, ok := value.(Entry[Value]); ok && entry.expiry > time.Now().Unix() {
			return entry.value, true
		}
	}
	return nil, false
}

func (c *Cache[Value]) Put(k Key, v *Value) {
	c.entries.Store(k, Entry[Value]{
		value:  v,
		expiry: time.Now().Unix() + 60,
	})
	if time.Now().Unix() > c.last_flush+5000 {
		c.last_flush = time.Now().Unix()
		go c.Flush()
	}
}

func (c *Cache[Value]) Flush() {
	c.entries.Range(func(key, value interface{}) bool {
		if entry, ok := value.(Entry[Value]); ok {
			if entry.expiry > time.Now().Unix() {
				c.entries.Delete(key)
			}
		}
		return true
	})
}
