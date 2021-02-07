package main

import (
	"fmt"
	"log"
	"net/http"

	"kentliuqiao.com/goo-cache/day3-http-server/goocache"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func main() {
	goocache.NewGroup("scores", 2<<10, goocache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exists", key)
		}))

	addr := "localhost:9999"
	pool := goocache.NewHTTPPool(addr)
	log.Println("goocache is running at " + addr)
	log.Fatal(http.ListenAndServe(addr, pool))
}
