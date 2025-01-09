package base64Captcha

import (
	"sync"
	"time"
)

// StoreSyncMap use sync.Map as store
type StoreSyncMap struct {
	liveTime time.Duration
	m        *sync.Map
}

// NewStoreSyncMap new a instance
func NewStoreSyncMap(liveTime time.Duration) *StoreSyncMap {
	return &StoreSyncMap{liveTime: liveTime, m: new(sync.Map)}
}

// smv a value type
type syncMap struct {
	time  time.Time
	Value string
}

// newSyncMap create a sync map instance
func newSyncMap(v string) *syncMap {
	return &syncMap{time: time.Now(), Value: v}
}

// delete remove expired items
func (s StoreSyncMap) delete() {
	expireTime := time.Now().Add(-s.liveTime)
	s.m.Range(func(key, value interface{}) bool {
		if sv, ok := value.(*syncMap); ok && sv.time.Before(expireTime) {
			s.m.Delete(key)
		}
		return true
	})
}

// Set a string value
func (s StoreSyncMap) Set(id string, value string) error {
	s.delete()
	s.m.Store(id, newSyncMap(value))
	return nil
}

// Get a string value
func (s StoreSyncMap) Get(id string, clear bool) string {
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
func (s StoreSyncMap) Verify(id, answer string, clear bool) bool {
	return s.Get(id, clear) == answer
}
