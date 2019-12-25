// Package app is the core of the application which manages both the back-end and front-end
package app

import (
	"context"

	"github.com/giannimassi/trello-tui/pkg/gui"
	"github.com/giannimassi/trello-tui/pkg/trello"
	"github.com/giannimassi/trello-tui/pkg/store"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Config is the application configuration
type Config struct {
	Trello trello.Config
	Gui   gui.Config
}

// App is the type in charge of handling application lifecycle
type App struct {
	l   zerolog.Logger
	cfg *Config

	gui               *gui.Gui
	updater           *trello.Updater
	cancelStateUpdate context.CancelFunc
}

// NewApp initializes a new instance of app with a logger and the provided configuration
func NewApp(cfg *Config) *App {
	return &App{
		l:   log.With().Str("m", "app").Logger(),
		cfg: cfg,
		gui: gui.NewGui(&cfg.Gui),
	}
}

// Init initializes the app's dependencies
func (a *App) Init() error {
	storeState, getState := store.NewStore(a.gui.Sync)
	a.updater = trello.NewUpdater(&a.cfg.Trello, storeState)
	if err := a.gui.Init(getState); err != nil {
		a.l.Error().Err(err).Msg("Unexpected error while initializing gui")
		return err
	}

	return nil
}

// Run executes the application's backend and frontend
func (a *App) Run() error {
	// Run state updater
	ctx, cancel := context.WithCancel(context.Background())
	go a.updater.Run(ctx)
	a.cancelStateUpdate = cancel

	// run gui
	return a.gui.Run()
}

// Close performs cleanup operations
func (a *App) Close() {
	a.l.Debug().Msg("Closing application")
	a.cancelStateUpdate()
	a.gui.Close()
}
