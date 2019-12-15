package gui

import (
	"sync"

	"github.com/rivo/tview"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

//Config is the gui configuration
type Config struct {
	Dev bool // Enable developer features (recover on run panics)
}

// Gui is a graphical user interface for trello-tui
type Gui struct {
	l    zerolog.Logger
	cfg  *Config
	app  *tview.Application
	view View

	state GetStateFunc

	rw sync.RWMutex
}

// NewGui creates a new instance of `Gui` from the provide configuration
func NewGui(cfg *Config) *Gui {
	app := tview.NewApplication()
	return &Gui{
		l:    log.Logger.With().Str("m", "gui").Logger(),
		cfg:  cfg,
		app:  app,
		view: NewView(app),
	}
}

// Init initializes the gui
func (g *Gui) Init() error {
	g.l.Info().Msg("Initialized")
	return nil
}

// Run executes the gui and event loop
func (g *Gui) Run(getState GetStateFunc) error {
	g.state = getState
	return g.app.Run()
}

// Close handles cleanup of the gui if required
func (g *Gui) Close() {
	g.l.Debug().Msg("Closing")
}

// Sync updates the gui with the latest state
func (g *Gui) Sync() {
	// apply state changes
	g.app.QueueUpdateDraw(func(){g.view.Draw(g.state())})
}
