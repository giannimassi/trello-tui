package gui

import (
	"sync"

	"github.com/giannimassi/trello-tui/pkg/store"
	"github.com/rivo/tview"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Config is the gui configuration
type Config struct {
	Dev bool // Enable developer features (recover on run panics)
}

// Gui is a graphical user interface for trello-tui
type Gui struct {
	l    zerolog.Logger
	cfg  *Config
	app  *tview.Application
	view *View

	state store.GetStateFunc

	rw sync.RWMutex
}

// NewGui creates a new instance of `Gui` from the provide configuration
func NewGui(cfg *Config) *Gui {
	return &Gui{
		l:   log.Logger.With().Str("m", "gui").Logger(),
		cfg: cfg,
	}
}

// Init initializes the gui
func (g *Gui) Init(getState store.GetStateFunc) error {
	g.l.Info().Msg("Initialized")
	g.app = tview.NewApplication()
	g.state = getState
	g.view = NewView(g.state(), g)
	return nil
}

// SetFocus implements the focuser interface, allowing to set focus on the provided gui component
func (g *Gui) SetFocus(p tview.Primitive) {
	g.app.QueueUpdateDraw(func() {
		g.app.SetFocus(p)
	})
}

// Run executes the gui and event loop
func (g *Gui) Run() error {
	g.app.SetRoot(g.view, true)
	g.app.SetFocus(g.view.FocusedItem())
	return g.app.Run()
}

// Close handles cleanup of the gui if required
func (g *Gui) Close() {
	g.l.Debug().Msg("Closing")
}

// Sync updates the gui with the latest state
func (g *Gui) Sync() {
	if g.view == nil {
		log.Error().Msg("nil view")
		return
	}

	s := g.state()
	if s == nil {
		log.Error().Msg("nil state")
		return
	}
	// apply state changes and draw
	g.app.QueueUpdateDraw(func() {
		defer func() {
			if err := recover(); err != nil {
				g.l.Error().Interface("recover", err).Msg("recovered panic while setting state")
			}
		}()
		g.view.SetState(s)
	})
}
