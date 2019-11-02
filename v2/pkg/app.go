package pkg

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/giannimassi/trello-tui/v2/pkg/gui"
	"github.com/giannimassi/trello-tui/v2/pkg/state"
)

type Config struct {
	Gui   gui.Config
	State state.Config
}

type App struct {
	l                 zerolog.Logger
	cfg               *Config
	gui               *gui.Gui
	state             *state.State
	cancelStateUpdate context.CancelFunc
}

func NewApp(cfg *Config) *App {
	return &App{
		l:   log.With().Str("m", "app").Logger(),
		cfg: cfg,
	}
}

func (a *App) Init() error {
	s := state.NewState(&a.cfg.State)
	a.state = s

	g := gui.NewGui(&a.cfg.Gui)
	if err := g.Init(); err != nil {
		a.l.Error().Err(err).Msg("Unexpected error while initializing gui")
		return err
	}
	a.gui = g

	return nil
}

func (a *App) Run() error {
	storeState, getState := gui.NewStateStore(a.state.State().(gui.State), a.gui.Sync)
	ctx, cancel := context.WithCancel(context.Background())
	go a.refreshState(ctx, storeState)
	a.cancelStateUpdate = cancel

	return a.gui.Run(getState)
}

func (a *App) refreshState(ctx context.Context, storeState gui.StoreStateFunc) {
	var (
		t   = time.NewTimer(0)
		err error
	)
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

func (a *App) Close() {
	a.l.Debug().Msg("Closing application")
	a.cancelStateUpdate()
	a.gui.Close()
}
