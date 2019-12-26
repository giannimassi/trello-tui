package store

import (
	"sync"
)

// Store makes State safe for concurrent use
type Store struct {
	State
	stateChangedFunc func()
	m                sync.RWMutex
}

// NewStore returns a new instate of Store
func NewStore(state State) *Store {
	return &Store{
		State:            state,
		stateChangedFunc: func() {},
	}
}

// BeginWrite acquires a write lock
func (s *Store) BeginWrite() { s.m.Lock() }

// EndWrite releases a write lock and calls `stateChangedFunc` if necessary
func (s *Store) EndWrite(changed bool) {
	s.m.Unlock()
	if changed {
		s.stateChangedFunc()
	}
}

// BeginRead acquires a read lock
func (s *Store) BeginRead() { s.m.RLock() }

// EndRead releases a read lock
func (s *Store) EndRead() { s.m.RUnlock() }

// SetStateChangedFunc configures a callback on state changes, called after unlocking
// write lock. Used to hook gui re=-draw to state changes.
func (s *Store) SetStateChangedFunc(f func()) {
	s.stateChangedFunc = f
}
