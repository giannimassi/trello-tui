package state

import (
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/giannimassi/trello-tui/pkg/domain"

	"github.com/giannimassi/trello-tui/pkg/trello"
)

type Config struct {
	Trello trello.Config

	SelectedBoard        string
	BoardRefreshInterval time.Duration
}

type board interface {
	online(*domain.Board) board
	offline(err error) board
}

type State struct {
	l      zerolog.Logger
	cfg    *Config
	client *trello.Client
	board
}

func NewState(cfg *Config) *State {
	return &State{
		l:   log.Logger.With().Str("m", "state").Logger(),
		cfg: cfg,
		board: &boardLoading{
			boardName: cfg.SelectedBoard,
		},
	}
}

func (s *State) ensureClientInitialized() error {
	if s.client == nil {
		client := trello.NewClient(&s.cfg.Trello)
		if err := client.Init(); err != nil {
			return err
		}
		s.client = client
	}
	return nil
}

// Update returns a copy of the current state with updated board / err
func (s *State) Update() (*State, error) {
	err := s.ensureClientInitialized()
	if err != nil {
		log.Error().Err(err).Msg("attaching err")
		s.board = s.board.offline(err)
		return s, err
	}
	b, err := s.client.Board(s.cfg.SelectedBoard)
	if err != nil {
		s.board = s.board.offline(err)
		return s, err
	}
	s.board = s.online(b)
	return s, nil
}

func (s *State) State() interface{} {
	return s.board
}
