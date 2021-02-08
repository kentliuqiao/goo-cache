package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// Hash maps bytes to uint32
// For IoC/DI and test convenience
type Hash func([]byte) uint32

// Map contains all hashed keys
type Map struct {
	hash Hash
	// virtual hosts multiple
	replicas int
	// hash circle for consistent hash
	keys []int
	// maps the vistual host's hash to its real host name
	hashMap map[int]string
}

// New constructs a new Map
func New(replicas int, fn func([]byte) uint32) *Map {
	m := &Map{
		hash:     fn,
		replicas: replicas,
		hashMap:  make(map[int]string),
	}

	if m.hash == nil {
		// use crc32.ChecksumIEEE for default hash function
		m.hash = crc32.ChecksumIEEE
	}

	return m
}

// Add adds some keys to the hash
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}
	}
	sort.Ints(m.keys)
}

// Get gets the closest item in the hash to the provided key
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}

	hash := int(m.hash([]byte(key)))
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})

	return m.hashMap[m.keys[idx%len(m.keys)]]
}
