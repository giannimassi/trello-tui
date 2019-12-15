package v2

import (
	"context"
	"time"

	"github.com/giannimassi/trello-tui/pkg/state"
	"github.com/giannimassi/trello-tui/pkg/v2/gui"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Config is the application configuration
type Config struct {
	State state.Config
	Gui   gui.Config	
}

// App is the type in charge of handling application lifecycle
type App struct {
	l   zerolog.Logger
	cfg *Config

	gui               *gui.Gui
	state             *state.State
	cancelStateUpdate context.CancelFunc
}

// NewApp initializes a new instance of app with a logger and the provided configuration
func NewApp(cfg *Config) *App {
	return &App{
		l:     log.With().Str("m", "app").Logger(),
		cfg:   cfg,
		gui:   gui.NewGui(&cfg.Gui),
		state: state.NewState(&cfg.State),
	}
}

// Init initializes the app's dependencies
func (a *App) Init() error {
	if err := a.gui.Init(); err != nil {
		a.l.Error().Err(err).Msg("Unexpected error while initializing gui")
		return err
	}

	return nil
}

// Run executes the application's backend and frontend
func (a *App) Run() error {
	storeState, getState := gui.NewStore(a.gui.Sync)
	ctx, cancel := context.WithCancel(context.Background())
	go a.pollState(ctx, storeState)
	a.cancelStateUpdate = cancel

	return a.gui.Run(getState)
}

// pollState is the go routine in charge of regularly updating the gui with the current state
func (a *App) pollState(ctx context.Context, storeState gui.StoreStateFunc) {
	var (
		t   = time.NewTimer(0)
		err error
	)
	a.l.Debug().Msg("Refresh state started")
	storeState(a.state.State().(gui.State))
	for {
		select {
		case <-ctx.Done():
			t.Stop()
			return

		case now := <-t.C:
			a.l.Debug().Msg("Refreshing state")
			if a.state, err = a.state.Update(); err != nil {
				a.l.Error().Err(err).Msg("Could not update board")
			}
			storeState(a.state.State().(gui.State))
			t.Reset(a.cfg.State.BoardRefreshInterval - time.Since(now))
		}
	}
}

// Close performs cleanup operations
func (a *App) Close() {
	a.l.Debug().Msg("Closing application")
	a.cancelStateUpdate()
	a.gui.Close()
}
