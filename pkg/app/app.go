package app

import (
	"context"
	"errors"
	"time"

	"github.com/giannimassi/trello-tui/pkg/gui"
	"github.com/giannimassi/trello-tui/pkg/store"
	"github.com/giannimassi/trello-tui/pkg/trello"
	"github.com/rs/zerolog"
)

type Config struct {
	Trello          *trello.Config
	Gui             *gui.Config
	StateFile       string
	RefreshInterval time.Duration
	SelectBoard     string
}

type App struct {
	l   zerolog.Logger
	cfg *Config

	client *trello.Client
	store  *store.JSONFileStore

	gui *gui.Gui

	cancelUpdate context.CancelFunc
}

func NewApp(l zerolog.Logger, cfg *Config) *App {
	if cfg.StateFile == "" {
		cfg.StateFile = "./.trello-tui"
	}

	return &App{
		l:   l,
		cfg: cfg,

		client: trello.NewClient(logger(l, "trello"), cfg.Trello),
		store:  store.NewJSONFileStore(logger(l, "store"), cfg.StateFile),
		gui:    gui.NewGui(logger(l, "gui"), cfg.Gui),
	}
}

func logger(l zerolog.Logger, module string) zerolog.Logger {
	return l.With().Str("m", module).Logger()
}

func (a *App) Init() error {
	// Init client
	if err := a.client.Init(); err != nil {
		a.l.Error().Err(err).Msg("Unexpected error while initializing trello client")
		return err
	}

	// Init store
	if err := a.store.Init(); err != nil {
		a.l.Error().Err(err).Msg("Unexpected error while initializing file store")
		return err
	}

	// Init gui
	if err := a.gui.Init(a.store.State); err != nil {
		a.l.Error().Err(err).Msg("Unexpected error while initializing gui")
		return err
	}

	s := a.store.State()
	if s.SelectedBoard == "" && a.cfg.SelectBoard == "" {
		err := errors.New("no board name provided")
		a.l.Error().Err(err).Msg("Unexpected error while initializing gui")
		return err
	}

	if s.SelectedBoard != a.cfg.SelectBoard {
		s.SelectedBoard = a.cfg.SelectBoard
		s.SetLoading(true)
	}

	return nil
}

// updateState retrieves fresh information via trello api and stores in the file cache.
// if boardName
func (a *App) updateState() error {
	s := a.store.State()
	board, lists, cards, err := a.client.BoardInfo(s.SelectedBoard)
	if err != nil {
		s.AppendErr(err)
		return nil
	}

	s.Updated = time.Now()
	s.Board = board
	s.Lists = lists
	s.Cards = cards
	s.SetLoading(false)

	// write to store also (will block state reads)
	err = a.store.Write(s)
	if err != nil {
		s.ErrorList = append(s.ErrorList, err)
		return nil
	}

	return nil
}

func (a *App) updateStateLoop(ctx context.Context) {
	var timer = time.NewTimer(0)
	for {
		select {
		case <-ctx.Done():
			if !timer.Stop() {
				<-timer.C
			}
			return
		case now := <-timer.C:
			if err := a.updateState(); err != nil {
				a.l.Error().Err(err).Msg("Unexpected error while updating state")
				continue
			}
			a.l.Info().Dur("d", time.Since(now)).Msg("State updated")
			a.gui.Sync()
			timer.Reset(a.cfg.RefreshInterval)
		}
	}
}

func (a *App) Run() error {
	ctx, cancelUpdate := context.WithCancel(context.Background())
	go a.updateStateLoop(ctx)
	a.cancelUpdate = cancelUpdate
	return a.gui.Run()
}

func (a *App) Close() {
	a.cancelUpdate()
	a.gui.Close()
}