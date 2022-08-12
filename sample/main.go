package main

import (
	"fmt"

	"github.com/tak1827/go-cache/lru"
)

func main() {
	size := 2
	ttl := 60 * 60 // 1h
	cache := lru.NewCache(size, ttl)

	// add key
	cache.Add("key1", "value")

	// get key
	v, ok := cache.Get("key1")
	if !ok {
		panic("key not found")
	}
	if v.(string) != "value" {
		panic("unexpected value")
	}
	fmt.Printf("get: %v\n", v)

	// remove key
	cache.Remove("key1")
	if cache.Len() != 0 {
		panic("unexpected length")
	}
}
