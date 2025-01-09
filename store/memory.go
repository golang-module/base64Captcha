package store

import (
	"container/list"
	"sync"
	"time"
)

type memory struct {
	time time.Time
	id   string
}

// memoryStore is an internal store for captcha ids and their values.
type memoryStore struct {
	sync.RWMutex
	idMap  map[string]string
	idList *list.List
	// Number of items stored since last collection.
	storedNum int
	// Number of saved items that triggers collection.
	collectNum int
	// Expiration time of captchas.
	expiration time.Duration
}

// NewMemoryStore returns a new standard memory store for captchas with the
// given collection threshold and expiration time (duration). The returned
// store must be registered with SetCustomStore to replace the default one.
func NewMemoryStore(collectNum int, expiration time.Duration) Store {
	s := new(memoryStore)
	s.idMap = make(map[string]string)
	s.idList = list.New()
	s.collectNum = collectNum
	s.expiration = expiration
	return s
}

func (s *memoryStore) Set(id string, value string) error {
	s.Lock()
	s.idMap[id] = value
	s.idList.PushBack(memory{time.Now(), id})
	s.storedNum++
	s.Unlock()
	if s.storedNum > s.collectNum {
		go s.collect()
	}
	return nil
}

func (s *memoryStore) Verify(id, answer string, clear bool) bool {
	if id == "" || answer == "" {
		return false
	}
	v := s.Get(id, clear)
	return v != "" && v == answer
}

func (s *memoryStore) Get(id string, clear bool) string {
	if !clear {
		// When we don't need to clear captcha, acquire read lock.
		s.RLock()
		defer s.RUnlock()
	} else {
		s.Lock()
		defer s.Unlock()
	}
	value, ok := s.idMap[id]
	if !ok {
		return ""
	}
	if clear {
		delete(s.idMap, id)
	}
	return value
}

func (s *memoryStore) collect() {
	now := time.Now()
	s.Lock()
	defer s.Unlock()
	for ele := s.idList.Front(); ele != nil; {
		mv, ok := ele.Value.(memory)
		if !ok {
			ele = nil
			continue
		}
		if mv.time.Add(s.expiration).Before(now) {
			delete(s.idMap, mv.id)
			next := ele.Next()
			s.idList.Remove(ele)
			s.storedNum--
			ele = next
		} else {
			ele = ele.Next()
		}
	}
}
