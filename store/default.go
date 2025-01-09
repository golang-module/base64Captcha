package store

import "time"

var (
	// GCLimitNumber The number of captchas created that triggers garbage collection used by default store.
	GCLimitNumber = 10240
	// Expiration time of captchas used by default store.
	Expiration = 10 * time.Minute

	DefaultMemStore = NewMemoryStore(GCLimitNumber, Expiration)
	DefaultMapStore = NewMapStore(Expiration)
)
