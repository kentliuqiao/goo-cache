package lru

import (
	"reflect"
	"testing"
)

type String string

func (s String) Len() int {
	return len(s)
}

func TestGet(t *testing.T) {
	cache := NewCache(int64(0), nil)
	cache.Add("key1", String("123"))
	if val, ok := cache.Get("key1"); !ok || string(val.(String)) != "123" {
		t.Fatalf("cache hit key1=123 failed")
	}

	if _, ok := cache.Get("key2"); ok {
		t.Fatalf("cache miss key2 failed")
	}
}

func TestRemoveOldest(t *testing.T) {
	k1, k2, k3 := "key1", "key2", "k3"
	v1, v2, v3 := "value1", "value2", "v3"
	cap := len(v1) + len(v2)
	cache := NewCache(int64(cap), nil)
	cache.Add(k1, String(v1))
	cache.Add(k2, String(v2))
	cache.Add(k3, String(v3))

	if _, ok := cache.Get(k1); ok || cache.Len() != 2 {
		t.Fatalf("Removeoldest key1 failed")
	}
}

func TestOnEvicted(t *testing.T) {
	keys := make([]string, 0)
	callBack := func(key string, val Value) {
		keys = append(keys, key)
	}

	cache := NewCache(int64(10), callBack)
	cache.Add("k1", String("1234"))
	cache.Add("k2", String("56785123123"))
	cache.Add("k3", String("91"))
	cache.Add("k4", String("1221"))

	expect := []string{"k1", "k2"}

	if !reflect.DeepEqual(expect, keys) {
		t.Fatalf("Call OnEvicted failed, expect keys equals to %s, got: %s", expect, keys)
	}
}
