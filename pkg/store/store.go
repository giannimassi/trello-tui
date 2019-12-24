package store

import (
	"fmt"
	"reflect"
	"sync"
	"sync/atomic"

	"github.com/rs/zerolog/log"
)

// PutStateFunc takes a gui.State interface as input
type PutStateFunc func(state State)

// GetStateFunc return a State interface
type GetStateFunc func() State

// NewStore returns a put and get state functions given the stateChanged callback provided
func NewStore(stateChanged func()) (PutStateFunc, GetStateFunc) {
	var (
		s        atomic.Value
		typeLock sync.RWMutex
	)
	putState := func(state State) {
		interfaceState, ok := state.(State)
		if !ok {
			log.Error().Interface("in", state).Str("type", fmt.Sprintf("%T", state)).Msg("Unexpected error: interface in is not State")
			return
		}

		previous := s.Load()
		if reflect.TypeOf(previous) != reflect.TypeOf(interfaceState) {
			typeLock.Lock()
			s = atomic.Value{}
			typeLock.Unlock()
		}
		s.Store(interfaceState)
		stateChanged()
	}

	getState := func() State {
		typeLock.RLock()
		defer typeLock.RUnlock()
		state := s.Load()
		stateInterface, ok := state.(State)
		if !ok {
			log.Error().Interface("out", state).Str("type", fmt.Sprintf("%T", state)).Msg("Unexpected error: interface out is not State")
			return nil
		}
		return stateInterface
	}
	return putState, getState
}
