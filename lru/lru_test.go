package lru

import (
	"testing"
)

func TestAddContains(t *testing.T) {
	cache := NewCache(2)

	tests := []struct {
		desc    string
		key     string
		value   interface{}
		evicted bool
	}{
		{
			desc:    "add string",
			key:     "key1",
			value:   "string1",
			evicted: false,
		},
		{
			desc:    "add bytes",
			key:     "key2",
			value:   []byte("byte"),
			evicted: false,
		},
		{
			desc:    "duplicated data",
			key:     "key1",
			value:   "string",
			evicted: false,
		},
		{
			desc:    "evicted",
			key:     "key3",
			value:   123,
			evicted: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			if g, w := cache.Add(tt.key, tt.value), tt.evicted; g != w {
				t.Errorf("unexpected added return, get: %t, want: %t", g, w)
			}
			if g, w := cache.Contains(tt.key), true; g != w {
				t.Errorf("unexpected added return, get: %t, want: %t", g, w)
			}
		})
	}
}

func TestGet(t *testing.T) {
	cache := NewCache(2)

	// get string
	key1 := "key1"
	value1 := "value"
	cache.Add(key1, value1)
	v, _ := cache.Get(key1)
	if g, w := v.(string), value1; g != w {
		t.Errorf("unexpected got value, get: %s, want: %s", g, w)
	}

	// get bytes
	key2 := "key2"
	value2 := []byte("byte")
	cache.Add(key2, value2)
	v, _ = cache.Get(key2)
	if g, w := v.([]byte), value2; string(g) != string(w) {
		t.Errorf("unexpected got value, get: %s, want: %s", g, w)
	}

	// get numbet and key1 is removed
	key3 := "key3"
	value3 := 123
	cache.Add(key3, value3)
	v, _ = cache.Get(key3)
	if g, w := v.(int), value3; g != w {
		t.Errorf("unexpected got value, get: %d, want: %d", g, w)
	}
	_, found := cache.Get(key1)
	if g, w := found, false; g != w {
		t.Errorf("expected key1 not found, get: %t, want: %t", g, w)
	}

	// get bool and key3 is removed
	cache.Add(key2, value2) // move key2 front
	key4 := "key4"
	value4 := true
	cache.Add(key4, value4)
	v, _ = cache.Get(key4)
	if g, w := v.(bool), value4; g != w {
		t.Errorf("unexpected got value, get: %t, want: %t", g, w)
	}
	_, found = cache.Get(key3)
	if g, w := found, false; g != w {
		t.Errorf("expected key3 not found, get: %t, want: %t", g, w)
	}
}

func TestRemoveLenCap(t *testing.T) {
	cache := NewCache(2)

	key1 := "key1"
	key2 := "key2"
	cache.Add(key1, "value1")
	cache.Add(key2, "value2")

	// delete key1
	if g, w := cache.Remove(key1), true; g != w {
		t.Errorf("unexpected remove return, get: %t, want: %t", g, w)
	}
	if g, w := cache.Len(), 1; g != w {
		t.Errorf("unexpected len, get: %d, want: %d", g, w)
	}

	// delete key1 again
	if g, w := cache.Remove(key1), false; g != w {
		t.Errorf("unexpected remove return , get: %t, want: %t", g, w)
	}

	// delete key2
	if g, w := cache.Remove(key2), true; g != w {
		t.Errorf("unexpected remove return, get: %t, want: %t", g, w)
	}
	if g, w := cache.Len(), 0; g != w {
		t.Errorf("unexpected len, get: %d, want: %d", g, w)
	}

	// capacity
	if g, w := cache.Cap(), 2; g != w {
		t.Errorf("unexpected capacity, get: %d, want: %d", g, w)
	}
}
