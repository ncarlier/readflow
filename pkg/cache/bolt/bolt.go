package boltcache

import (
	"container/list"
	"fmt"
	"net/url"
	"sync"

	"github.com/ncarlier/readflow/pkg/values"
	bolt "go.etcd.io/bbolt"
)

const (
	defaultMaxEntrySize = 5   // Mb
	defaultMaxSize      = 256 // Mb
	defaultMaxEntries   = 1024
)

var bucketName = []byte("cache")

// BoltLRUCache implement a LRU cache with BoltDB
type BoltLRUCache struct {
	maxEntrySize int
	maxSize      int
	maxEntries   int
	evictionList *list.List
	entries      map[string]*list.Element
	lock         sync.RWMutex
	size         int
	db           *bolt.DB
}

// New LRU cache using BoltDB as storage backend
func New(path string, params url.Values) (*BoltLRUCache, error) {
	db, err := bolt.Open(path, 0o640, nil)
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
	maxSize := values.GetIntOrDefault(params, "maxSize", defaultMaxSize)
	maxEntries := values.GetIntOrDefault(params, "maxEntries", defaultMaxEntries)
	maxEntrySize := values.GetIntOrDefault(params, "maxEntrySize", defaultMaxEntrySize)
	c := &BoltLRUCache{
		maxSize:      maxSize * 1024 * 1024,
		maxEntries:   maxEntries,
		maxEntrySize: maxEntrySize * 1024 * 1024,
		evictionList: list.New(),
		entries:      make(map[string]*list.Element, maxEntries),
		db:           db,
	}
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		cur := b.Cursor()
		var key string
		for k, v := cur.First(); k != nil; k, v = cur.Next() {
			key = string(k)
			el := c.evictionList.PushFront(key)
			c.entries[key] = el
			c.size += len(v)
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
		return err
	})
	if value != nil {
		c.lock.Lock()
		if entry, ok := c.entries[key]; ok {
			c.evictionList.MoveToFront(entry)
		}
		c.lock.Unlock()
	}
	return
}

// Put item into the cache
func (c *BoltLRUCache) Put(key string, value []byte) error {
	if len(value) > c.maxEntrySize {
		return fmt.Errorf("entry size exceeding MaxEntrySize: %d > %d", len(value), c.maxEntrySize)
	}
	return c.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketName)
		c.lock.Lock()
		defer c.lock.Unlock()
		oldEntrySize := 0
		entry, exists := c.entries[key]
		if exists {
			oldEntrySize = len(bucket.Get([]byte(key)))
		}

		// update cache
		if err := bucket.Put([]byte(key), value); err != nil {
			return err
		}

		// update cache size
		c.size = c.size - oldEntrySize + len(value)

		// update eviction list
		if exists {
			c.evictionList.MoveToFront(entry)
		} else {
			el := c.evictionList.PushFront(key)
			c.entries[key] = el
		}

		// run eviction
		return c.runEviction(tx)
	})
}

// Clear cache (for testing usage)
func (c *BoltLRUCache) Clear() error {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.db.Update(func(tx *bolt.Tx) error {
		if err := tx.DeleteBucket(bucketName); err != nil {
			return err
		}
		if _, err := tx.CreateBucketIfNotExists(bucketName); err != nil {
			return err

		}
		c.evictionList = list.New()
		c.entries = make(map[string]*list.Element, c.maxEntries)
		c.size = 0
		return nil
	})
}

func (c *BoltLRUCache) runEviction(tx *bolt.Tx) error {
	bucket := tx.Bucket(bucketName)
	for c.evictionList.Len() > c.maxEntries || c.size > c.maxSize {
		entry := c.evictionList.Back()
		if entry != nil {
			k := entry.Value.(string)
			c.evictionList.Remove(entry)
			delete(c.entries, k)
			c.size -= len(bucket.Get([]byte(k)))
			return bucket.Delete([]byte(k))
		}
	}
	return nil
}

// Close cache
func (c *BoltLRUCache) Close() error {
	return c.db.Close()
}
