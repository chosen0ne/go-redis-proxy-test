/**
 *
 * @author  chosen0ne(louzhenlin86@126.com)
 * @date    2017-11-17 18:51:41
 */

package store

import (
	"sync"
)

type MemStore struct {
	store map[string]string
	mu    *sync.Mutex
}

func NewMemStore() *MemStore {
	return &MemStore{make(map[string]string), &sync.Mutex{}}
}

func (s *MemStore) Set(key, val string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.store[key] = val
}

func (s *MemStore) Get(key string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	v, ok := s.store[key]
	return v, ok
}
