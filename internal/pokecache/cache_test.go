package pokecache

import (
	// "slices"
	"slices"
	"testing"
	"time"
)

func addLoop(cache *Cache, done <-chan bool, value []byte) {
	for {
		select {
		case <-done:
			return
		default:
			cache.Add("entry", value)
		}
	}
}

func readLoop(cache *Cache, done <-chan bool, responseCh chan<- []byte) {
	for {
		select {
		case <-done:
			return
		default:
			if result, ok := cache.Get("entry"); ok {
				responseCh <- result
			} else {
				responseCh <- []byte{}
			}
		}
	}
}

func TestCache(t *testing.T) {
	interval := 10 * time.Second
	cache := NewCache(interval)
	defer cache.Stop()

	t.Run("adding values to the cache", func(t *testing.T) {
		value := []byte{'h', 'e', 'y'}
		cache.Add("entry", value)
		result, ok := cache.Get("entry")
		if !ok {
			t.Fatalf("the added value for %s is not in the cache", "entry")
		}
		if slices.Compare(value, result) != 0 {
			t.Fatalf("expected resulting value: %v, got: %v", value, result)
		}
	})

	t.Run("pruning values from cache works", func(t *testing.T) {
		value := []byte{}
		cache.Add("entry", value)
		time.Sleep(interval * 2)

		if _, ok := cache.Get("entry"); ok {
			t.Fatalf("the value is still in the cache after %v", interval*2)
		}
	})

	t.Run("concurrency does not affect the contents of the cache", func(t *testing.T) {
		value := []byte{'H', 'E', 'L', 'L', 'O'}
		done := make(chan bool, 1)
		responseCh := make(chan []byte, 1)
		go readLoop(cache, done, responseCh)
		go addLoop(cache, done, value)

		for range 10 {
			time.Sleep(1 * time.Second)
			response := <-responseCh

			if slices.Compare(value, response) != 0 {
				t.Fatalf("the expected response value was not received")
			}
		}
	})
}
