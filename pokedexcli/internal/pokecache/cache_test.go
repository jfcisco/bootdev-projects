package pokecache

import (
	"fmt"
	"testing"
	"time"
)

func TestAddCache(t *testing.T) {
	const interval = time.Minute
	cases := []struct {
		key   string
		value []byte
	}{
		{
			key:   "http://pokeapi.dev/1",
			value: []byte("value-for-testing-1"),
		},
		{
			key:   "http://pokeapi.dev/2",
			value: []byte("value-for-testing-2"),
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			cache := NewCache(interval)
			cache.Add(c.key, c.value)
			val, ok := cache.Get(c.key)
			if !ok {
				t.Errorf("cache.Get(key) = %v, want = true", ok)
			} else if string(val) != string(c.value) {
				t.Errorf("cache.Get(key) = %s, want = %s", val, c.value)
			}
		})
	}
}

func TestReapLoop(t *testing.T) {
	const reapKey string = "reapKey"
	cache := NewCache(5 * time.Millisecond)
	cache.Add(reapKey, []byte("i should not exist"))

	_, ok := cache.Get(reapKey)
	if !ok {
		t.Errorf("cache.Get(reapKey); expected to find key")
	}

	<-time.After(20 * time.Millisecond)

	data, ok := cache.Get(reapKey)
	if ok {
		t.Errorf("cache.Get(reapKey) = (%s, %v); want !ok", data, ok)
	}
}
