package pokecache

import (
	"testing"
	"time"
)

func TestAddCache(t *testing.T) {
	cache := NewCache(time.Minute)
	key, value := "key", []byte("value")
	cache.Add(key, value)

	actual, ok := cache.Get(key)
	if !ok {
		t.Errorf("cache.Get(key) = %v, want = true", ok)
	} else if string(actual) != "value" {
		t.Errorf("cache.Get(key) = %s, want = \"value\"", actual)
	}
}

func TestReapLoop(t *testing.T) {
	const reapKey string = "reapKey"
	cache := NewCache(5 * time.Millisecond)
	cache.Add(reapKey, []byte("i should not exist"))

	<-time.After(20 * time.Millisecond)

	data, ok := cache.Get(reapKey)
	if ok {
		t.Errorf("cache.Get(reapKey) = (%s, %v); want !ok", data, ok)
	}
}