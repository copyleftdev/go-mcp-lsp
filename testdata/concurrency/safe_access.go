package concurrency

import (
	"sync"
)

type SafeStore struct {
	mu   sync.RWMutex
	data map[string]string
}

func NewSafeStore() *SafeStore {
	return &SafeStore{
		data: make(map[string]string),
	}
}

func (s *SafeStore) Set(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
}

func (s *SafeStore) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.data[key]
	return val, ok
}

func (s *SafeStore) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, key)
}

func UseSafeStore() {
	store := NewSafeStore()
	
	var wg sync.WaitGroup
	wg.Add(3)
	
	go func() {
		defer wg.Done()
		store.Set("key1", "value1")
	}()
	
	go func() {
		defer wg.Done()
		store.Delete("key2")
	}()
	
	go func() {
		defer wg.Done()
		_, _ = store.Get("key1")
	}()
	
	wg.Wait()
}
