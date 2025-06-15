package core

import (
	"sync"
)

func CreateCache() *KVStore {
	store := requestStorePool.Get().(*KVStore)
	store.Reset()
	return store
}

// 1. 定义一个线程安全的 K/V 存储结构
type KVStore struct {
	mu   sync.RWMutex
	data map[string]interface{}
}

func (s *KVStore) Set(key string, value interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
}

func (s *KVStore) Get(key string) (map[string]any, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, ok := s.data[key]
	if !ok {
		return nil, false
	}
	return value.(map[string]any), true
}

func (s *KVStore) Recover() {
	s.Reset()
	requestStorePool.Put(s)
}

// Reset 清空 map 以便被复用
func (s *KVStore) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	// 关键：不是创建新 map，而是清空旧 map
	// 这样可以保留已分配的 map 容量，提高效率
	for k := range s.data {
		delete(s.data, k)
	}
}

// 创建一个 sync.Pool 来管理 RequestStore 实例
var requestStorePool = sync.Pool{
	New: func() interface{} {
		return &KVStore{
			data: make(map[string]any),
		}
	},
}
