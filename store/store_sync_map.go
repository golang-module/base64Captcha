package store

import (
	"strings"
	"sync"
	"time"
)

type itemSyncMap struct {
	time  time.Time
	Value string
}

// storeSyncMap use sync.Map as store
type storeSyncMap struct {
	m *sync.Map
	d time.Duration
}

// NewStoreSyncMap create a new sync map store
func NewStoreSyncMap(duration time.Duration) *storeSyncMap {
	return &storeSyncMap{d: duration, m: new(sync.Map)}
}

// Set a string value
func (s storeSyncMap) Set(id string, value string) error {
	s.delete()
	s.m.Store(id, &itemSyncMap{time: time.Now(), Value: value})
	return nil
}

// Get a string value
func (s storeSyncMap) Get(id string, clear bool) string {
	v, ok := s.m.Load(id)
	if !ok {
		return ""
	}
	s.m.Delete(id)
	if sv, ok := v.(*itemSyncMap); ok {
		return sv.Value
	}
	return ""
}

// Verify check a string value
func (s storeSyncMap) Verify(id, answer string, clear bool) bool {
	vv := s.Get(id, clear)
	return strings.EqualFold(vv, answer)
}

// delete remove expired items
func (s storeSyncMap) delete() {
	expireTime := time.Now().Add(-s.d)
	s.m.Range(func(key, value interface{}) bool {
		if sv, ok := value.(*itemSyncMap); ok && sv.time.Before(expireTime) {
			s.m.Delete(key)
		}
		return true
	})
}
