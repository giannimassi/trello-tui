package trello

import (
	"sync"
	"time"

	"github.com/giannimassi/trello-tui/pkg/domain"
	"github.com/giannimassi/trello-tui/pkg/store"
)

// Config is the state configuration, including trello authentication details
type Config struct {
	User, Key, Token     string
	Timeout              time.Duration
	SelectedBoard        string
	BoardRefreshInterval time.Duration
}

// board is the interface used for managing state behaviour changes
type board interface {
	online(*domain.Board) board
	offline(err error) board
}

// state implements github.com/giannimassi/trello-tui/pkg/gui `gui.State`
type state struct {
	cfg    *Config
	client *Client
	board
	m sync.RWMutex
}

// newState returns a new instance of state
func newState(cfg *Config) *state {
	return &state{
		cfg: cfg,
		board: &boardLoading{
			boardName: cfg.SelectedBoard,
		},
	}
}

// ensureClientInitialized checks that the client has been initialized
func (s *state) ensureClientInitialized() error {
	if s.client == nil {
		client := NewClient(s.cfg)
		if err := client.Init(); err != nil {
			return err
		}
		s.client = client
	}
	return nil
}

// update returns a copy of the current state with updated board / err
func (s *state) update() (*state, error) {
	err := s.ensureClientInitialized()
	if err != nil {
		s.setBoardOffline(err)
		return s, err
	}
	b, err := s.client.Board(s.cfg.SelectedBoard)
	if err != nil {
		s.setBoardOffline(err)
		return s, err
	}
	s.setBoardOnline(b)
	return s, nil
}

func (s *state) updateCardComments(id int) {

}

func (s *state) setBoardOnline(b *domain.Board) {
	s.BeginWrite()
	s.board = s.online(b)
	s.EndWrite()
}

func (s *state) setBoardOffline(err error) {
	s.BeginWrite()
	s.board = s.board.offline(err)
	s.EndWrite()
}

func (s *state) storable() store.State {
	return s.board.(store.State)
}

func (s *state) BeginWrite() {
	s.m.Lock()
}

func (s *state) EndWrite() {
	s.m.Unlock()
}

func (s *state) BeginRead() {
	s.m.RLock()
}

func (s *state) EndRead() {
	s.m.RUnlock()
}
