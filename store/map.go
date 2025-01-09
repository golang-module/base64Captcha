package store

import (
	"sync"
	"time"
)

// NewMapStore new a instance
func NewMapStore(duration time.Duration) *MapStore {
	return &MapStore{d: duration, m: new(sync.Map)}
}

type syncMap struct {
	time  time.Time
	Value string
}

// MapStore use sync.Map as store
type MapStore struct {
	d time.Duration
	m *sync.Map
}

// Set a string value
func (s MapStore) Set(id string, value string) error {
	s.delete()
	s.m.Store(id, &syncMap{time: time.Now(), Value: value})
	return nil
}

// Get a string value
func (s MapStore) Get(id string, clear bool) string {
	v, ok := s.m.Load(id)
	if !ok {
		return ""
	}
	s.m.Delete(id)
	if sv, ok := v.(*syncMap); ok {
		return sv.Value
	}
	return ""
}

// Verify check a string value
func (s MapStore) Verify(id, answer string, clear bool) bool {
	return s.Get(id, clear) == answer
}

// delete remove expired items
func (s MapStore) delete() {
	expireTime := time.Now().Add(-s.d)
	s.m.Range(func(key, value interface{}) bool {
		if sv, ok := value.(*syncMap); ok && sv.time.Before(expireTime) {
			s.m.Delete(key)
		}
		return true
	})
}
