package store

import (
	"sync"
	"time"
)

// NewSyncMapStore new a instance
func NewSyncMapStore(duration time.Duration) *syncMapStore {
	return &syncMapStore{d: duration, m: new(sync.Map)}
}

type syncMap struct {
	time  time.Time
	Value string
}

// syncMapStore use sync.Map as store
type syncMapStore struct {
	d time.Duration
	m *sync.Map
}

// Set a string value
func (s syncMapStore) Set(id string, value string) error {
	s.delete()
	s.m.Store(id, &syncMap{time: time.Now(), Value: value})
	return nil
}

// Get a string value
func (s syncMapStore) Get(id string, clear bool) string {
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
func (s syncMapStore) Verify(id, answer string, clear bool) bool {
	return s.Get(id, clear) == answer
}

// delete remove expired items
func (s syncMapStore) delete() {
	expireTime := time.Now().Add(-s.d)
	s.m.Range(func(key, value interface{}) bool {
		if sv, ok := value.(*syncMap); ok && sv.time.Before(expireTime) {
			s.m.Delete(key)
		}
		return true
	})
}
