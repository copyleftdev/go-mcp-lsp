package concurrency

import (
	"sync"
)

// DataStore demonstrates a race condition due to unsynchronized access
type DataStore struct {
	data map[string]string
}

func NewDataStore() *DataStore {
	return &DataStore{
		data: make(map[string]string),
	}
}

// Race condition: concurrent map access without synchronization
func (ds *DataStore) Set(key, value string) {
	ds.data[key] = value
}

func (ds *DataStore) Get(key string) string {
	return ds.data[key]
}

func (ds *DataStore) Delete(key string) {
	delete(ds.data, key)
}

// UseDataStore shows concurrent access that would cause race conditions
func UseDataStore() {
	store := NewDataStore()
	
	var wg sync.WaitGroup
	wg.Add(3)
	
	go func() {
		defer wg.Done()
		store.Set("key1", "value1")
	}()
	
	go func() {
		defer wg.Done()
		store.Delete("key1")
	}()
	
	go func() {
		defer wg.Done()
		_ = store.Get("key1")
	}()
	
	wg.Wait()
}
