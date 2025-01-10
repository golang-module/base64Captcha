package store

import (
	"container/list"
	"strings"
	"sync"
	"time"
)

type itemMemory struct {
	time time.Time
	id   string
}

// storeMemory is an internal store for captcha ids and their values.
type storeMemory struct {
	rw sync.RWMutex
	m  map[string]string
	l  *list.List
	// Number of items stored since last collection.
	s int
	// Number of saved items that triggers collection.
	c int
	// Expiration time of captchas.
	d time.Duration
}

// NewStoreMemory returns a new standard memory store for captchas with the
// given collection threshold and expiration time (duration). The returned
// store must be registered with SetCustomStore to replace the default one.
func NewStoreMemory(collectNum int, expiration time.Duration) Store {
	s := new(storeMemory)
	s.m = make(map[string]string)
	s.l = list.New()
	s.c = collectNum
	s.d = expiration
	return s
}

func (s *storeMemory) Set(id string, value string) error {
	s.rw.Lock()
	s.m[id] = value
	s.l.PushBack(itemMemory{time.Now(), id})
	s.s++
	s.rw.Unlock()
	if s.s > s.c {
		go s.collect()
	}
	return nil
}

func (s *storeMemory) Verify(id, answer string, clear bool) bool {
	if id == "" || answer == "" {
		return false
	}
	v := s.Get(id, clear)
	return strings.EqualFold(v, answer)
}

func (s *storeMemory) Get(id string, clear bool) string {
	if !clear {
		// When we don't need to clear captcha, acquire read lock.
		s.rw.RLock()
		defer s.rw.RUnlock()
	} else {
		s.rw.Lock()
		defer s.rw.Unlock()
	}
	value, ok := s.m[id]
	if !ok {
		return ""
	}
	if clear {
		delete(s.m, id)
	}
	return value
}

func (s *storeMemory) collect() {
	now := time.Now()
	s.rw.Lock()
	defer s.rw.Unlock()
	for ele := s.l.Front(); ele != nil; {
		mv, ok := ele.Value.(itemMemory)
		if !ok {
			ele = nil
			continue
		}
		if mv.time.Add(s.d).Before(now) {
			delete(s.m, mv.id)
			next := ele.Next()
			s.l.Remove(ele)
			s.s--
			ele = next
		} else {
			ele = ele.Next()
		}
	}
}
