package goocache

import (
	"fmt"
	"log"
	"reflect"
	"testing"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func TestGetter(t *testing.T) {
	var f Getter = GetterFunc(func(key string) ([]byte, error) {
		return []byte(key), nil
	})
	expect := []byte("key")
	if v, _ := f.Get("key"); !reflect.DeepEqual(v, expect) {
		t.Error("callback failed")
	}
}

func TestGroupGet(t *testing.T) {
	// to determine whether getter are called multiple times
	loadCounts := make(map[string]int, len(db))
	scoreGroup := NewGroup("scores", 2<<10, GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				if _, ok := loadCounts[v]; !ok {
					loadCounts[v] = 0
				}
				loadCounts[v]++
				return []byte(v), nil
			}

			return nil, fmt.Errorf("%s not exist", key)
		}))

	for k, v := range db {
		// load from callback
		if view, err := scoreGroup.Get(k); err != nil || view.String() != v {
			t.Errorf("failed to get value of %s", k)
		}
		// load from cache
		if _, err := scoreGroup.Get(k); err != nil || loadCounts[k] > 1 {
			t.Errorf("cache %s miss", k)
		}
	}

	if view, err := scoreGroup.Get("unknown"); err == nil {
		t.Errorf("the values of key %s should be empty, but got %s", "unknown", view)
	}
}
