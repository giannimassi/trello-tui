package gui

import (
	"reflect"
	"sync"
	"sync/atomic"
)

// StoreStateFunc takes a State interface as input
type StoreStateFunc func(state State)

// GetStateFunc return a State interface
type GetStateFunc func() State

// NewStore returns a store and get state functions given the stateChanged callback provided
func NewStore(stateChanged func()) (StoreStateFunc, GetStateFunc) {
	var (
		s        atomic.Value
		typeLock sync.RWMutex
	)
	storeState := func(state State) {
		previous := s.Load()
		if reflect.TypeOf(previous) != reflect.TypeOf(state) {
			typeLock.Lock()
			s = atomic.Value{}
			typeLock.Unlock()
		}
		s.Store(state)
		stateChanged()
	}

	getState := func() State {
		typeLock.RLock()
		defer typeLock.RUnlock()
		return s.Load().(State)
	}

	// storeState(state)
	return storeState, getState
}
