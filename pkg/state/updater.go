package state

import (
	"context"
	"time"

	"github.com/giannimassi/trello-tui/pkg/store"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Updater ensures the state is updated as required
type Updater struct {
	l   zerolog.Logger
	cfg *Config

	*state
	put store.PutStateFunc
}

// NewUpdater returns a new instance of Updater
func NewUpdater(cfg *Config, put store.PutStateFunc) *Updater {
	u := Updater{
		l:     log.Logger.With().Str("m", "state-update").Logger(),
		cfg:   cfg,
		state: newState(cfg),
		put:   put,
	}
	log.Info().Msg("updating state init")
	put(u.storable())
	return &u
}

// Run executes the loop exuting the requests via the trello client
func (u *Updater) Run(ctx context.Context) {
	var (
		t   = time.NewTimer(0)
		err error
	)
	u.l.Debug().Msg("Refresh state started")
	for {
		select {
		case <-ctx.Done():
			t.Stop()
			return

		case now := <-t.C:
			u.l.Debug().Msg("Refreshing state")
			if _, err = u.update(); err != nil {
				u.l.Error().Err(err).Msg("Could not update board")
			}
			u.put(u.storable())
			t.Reset(u.cfg.BoardRefreshInterval - time.Since(now))
		}
	}
}
