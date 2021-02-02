package lru

import "container/list"

// Cache is a LRU cache. It's not safe for concurrent access.
type Cache struct {
	// the maximum memory space Cache can use
	maxBytes int64
	// the space Cache has already consumed
	nBytes int64
	// the value of the map is a pointer to a element in the doubly linked list
	cache map[string]*list.Element
	// doubly linked list
	ll *list.List
	// optional and executed when an entry is purged if it's set
	OnEvicted func(key string, value Value)
}

// entry is the data type stored in the list.
// It's handy to delete the key in the cache when the front element of the list
// is poped from the list.
type entry struct {
	key   string
	value Value
}

// Value uses Len to count bytes it takes
type Value interface {
	Len() int
}

// NewCache constructs a new Cache.
func NewCache(maxBytes int64, onEvicted func(key string, value Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		cache:     make(map[string]*list.Element),
		ll:        list.New(),
		OnEvicted: onEvicted,
	}
}

// Get gets a value by its key.
func (c *Cache) Get(key string) (Value, bool) {
	if ele, ok := c.cache[key]; ok {
		// when the key exists, move the element to the front of the list
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, ok
	}

	return nil, false
}

// RemoveOldest removes the oldest element in the list
// which is at the back of the list
func (c *Cache) RemoveOldest() {
	// get the oldest element
	ele := c.ll.Back()
	if ele != nil {
		// remove it from the list
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		// delete it from the map
		delete(c.cache, kv.key)
		// update the bytes Cache consumed
		c.nBytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// Add adds a value to the cache
func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		// if the key has already existed, update its value and push it to the front of the list
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nBytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		// if the key does not exist, push a new element to the front of the list
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nBytes += int64(value.Len())
	}
	// if cache consumes bytes more than the maxBytes, remove old element from the list till less than maxBytes
	// if maxBytes is 0, there is no constriction of the space cache can use and it's dangerous!
	for c.maxBytes > 0 && c.maxBytes < c.nBytes {
		c.RemoveOldest()
	}
}

// Len gets the number of cached entries
func (c *Cache) Len() int {
	return c.ll.Len()
}
