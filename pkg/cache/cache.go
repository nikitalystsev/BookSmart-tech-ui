package cache

import (
	"fmt"
	"reflect"
	"sync"
)

type ICache interface {
	Set(key string, value interface{})
	Get(key string, dest interface{}) error
	Delete(key string)
	Clear()
}

// Cache represents an in-memory key-value store.
type Cache struct {
	data map[string]interface{}
	mu   sync.RWMutex
}

// NewCache creates and initializes a new Cache instance.
func NewCache() *Cache {
	return &Cache{
		data: make(map[string]interface{}),
	}
}

// Set adds or updates a key-value pair in the cache.
func (c *Cache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
}

func (c *Cache) Get(key string, dest interface{}) error {
	c.mu.RLock()
	value, ok := c.data[key]
	c.mu.RUnlock()
	if !ok {
		return fmt.Errorf("key not found: %s", key)
	}

	// Используем reflect для копирования значения
	vDest := reflect.ValueOf(dest)
	if vDest.Kind() != reflect.Ptr || vDest.IsNil() {
		return fmt.Errorf("destination must be a non-nil pointer")
	}

	vValue := reflect.ValueOf(value)
	vDest.Elem().Set(vValue)

	return nil
}

// Delete removes a key-value pair from the cache.
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
}

// Clear removes all key-value pairs from the cache.
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data = make(map[string]interface{})
}
