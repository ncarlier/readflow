package boltcache

import (
	"container/list"
	"errors"
	"sync"

	bolt "go.etcd.io/bbolt"
)

var bucketName = []byte("cache")

// BoltLRUCache implement a LRU cache with BoltDB
type BoltLRUCache struct {
	size         int
	evictionList *list.List
	items        map[string]*list.Element
	lock         sync.RWMutex
	db           *bolt.DB
}

// New LRU cache using BoltDB as storage backend
func New(size int, path string) (*BoltLRUCache, error) {
	if size <= 0 {
		return nil, errors.New("invalide cache size")
	}
	db, err := bolt.Open(path, 0640, nil)
	if err != nil {
		return nil, err
	}
	err = db.Update(func(tx *bolt.Tx) (err error) {
		_, err = tx.CreateBucketIfNotExists(bucketName)
		return err
	})
	if err != nil {
		return nil, err
	}
	c := &BoltLRUCache{
		size:         size,
		evictionList: list.New(),
		items:        make(map[string]*list.Element, size),
		db:           db,
	}
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		cur := b.Cursor()
		var key string
		for k, _ := cur.First(); k != nil; k, _ = cur.Next() {
			key = string(k)
			el := c.evictionList.PushFront(key)
			c.items[key] = el
		}
		return nil
	})
	return c, nil
}

// Get item from cache
func (c *BoltLRUCache) Get(key string) (value []byte, err error) {
	err = c.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		value = b.Get([]byte(key))
		return nil
	})
	if value != nil {
		c.lock.Lock()
		if item, ok := c.items[key]; ok {
			c.evictionList.MoveToFront(item)
		}
		c.lock.Unlock()
	}
	return
}

// Put item into the cache
func (c *BoltLRUCache) Put(key string, data []byte) error {
	return c.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		err := b.Put([]byte(key), data)
		if err != nil {
			return err
		}
		c.lock.Lock()
		defer c.lock.Unlock()
		if item, ok := c.items[key]; ok {
			c.evictionList.MoveToFront(item)
		} else {
			el := c.evictionList.PushFront(key)
			c.items[key] = el
			for c.evictionList.Len() > c.size {
				el := c.evictionList.Back()
				if el != nil {
					k := el.Value.(string)
					c.evictionList.Remove(el)
					delete(c.items, k)
					return b.Delete([]byte(k))
				}
			}
		}
		return nil
	})
}

// Close cache
func (c *BoltLRUCache) Close() error {
	return c.db.Close()
}
