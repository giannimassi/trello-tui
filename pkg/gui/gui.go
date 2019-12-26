package gui

import (
	"runtime/debug"

	"github.com/gdamore/tcell"
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
}

// NewGui creates a new instance of `Gui` from the provide configuration
func NewGui(cfg *Config) *Gui {
	return &Gui{
		l:   log.Logger.With().Str("m", "gui").Logger(),
		cfg: cfg,
		app: tview.NewApplication(),
	}
}

// Init initializes the gui
func (g *Gui) Init(store *store.Store) error {
	g.view = NewView(store, g)
	// Hook state changed to re-draw func
	store.SetStateChangedFunc(g.ReDraw)
	// Lock state before drawing
	g.app.SetBeforeDrawFunc(func(screen tcell.Screen) bool {
		store.BeginRead()
		return false
	})
	// Unlock state after drawing
	g.app.SetAfterDrawFunc(func(screen tcell.Screen) {
		store.EndRead()
	})
	g.l.Info().Msg("Initialized")
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
	defer func() {
		if err := recover(); err != nil {
			log.Error().Interface("recover", err).Msg("panic in run: \n" + string(debug.Stack()))
		}
	}()

	g.app.SetRoot(g.view, true)
	g.app.SetFocus(g.view.FocusedItem())
	return g.app.Run()
}

// Close handles cleanup of the gui if required
func (g *Gui) Close() {
	g.l.Debug().Msg("Closing")
}

// ReDraw updates the gui with the latest state
func (g *Gui) ReDraw() {
	// trigger re-draw
	g.app.QueueUpdateDraw(func() {})
}
