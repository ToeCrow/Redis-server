package server

import "sync"

// KVStore holds string key-value pairs safe for concurrent use.
type KVStore struct {
	mu   sync.RWMutex
	data map[string]string
}

// NewKVStore returns an empty key-value store.
func NewKVStore() *KVStore {
	return &KVStore{
		data: make(map[string]string),
	}
}

// Get returns the value for key and whether it exists.
func (s *KVStore) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.data[key]
	return v, ok
}

// Set stores value under key.
func (s *KVStore) Set(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
}
