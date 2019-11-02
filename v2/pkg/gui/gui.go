package gui

import (
	"reflect"
	"runtime/debug"
	"sync"
	"sync/atomic"

	"github.com/jroimartin/gocui"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/giannimassi/trello-tui/v2/pkg/gui/components"
	"github.com/giannimassi/trello-tui/v2/pkg/gui/theme"
)

type Config struct {
	Dev bool // Enable developer features (recover on run panics)
}

type UserActionHandler interface {
	LeftPressed()
	RightPressed()
	UpPressed()
	DownPressed()
	EnterPressed()
	BackPressed()
}

type State interface {
	UserActionHandler
	components.ViewState
}

type StoreStateFunc func(state State)
type GetStateFunc func() State

func NewStateStore(state State, stateChanged func()) (StoreStateFunc, GetStateFunc) {
	var (
		s        atomic.Value
		typeLock sync.RWMutex
	)
	storeState := func(state State) {
		previous := s.Load()
		if reflect.TypeOf(previous) != reflect.TypeOf(state) {
			typeLock.Lock()
			s = atomic.Value{}
			typeLock.Unlock()
		}
		s.Store(state)
		stateChanged()
	}

	getState := func() State {
		typeLock.RLock()
		defer typeLock.RUnlock()
		return s.Load().(State)
	}

	storeState(state)
	return storeState, getState
}

type Gui struct {
	l   zerolog.Logger
	cfg *Config

	gui   *gocui.Gui
	theme *theme.Theme
	view  *components.View
	state GetStateFunc

	rw sync.RWMutex
}

func NewGui(cfg *Config) *Gui {
	return &Gui{
		l:     log.Logger.With().Str("m", "gui").Logger(),
		cfg:   cfg,
		theme: theme.DefaultTheme(),
	}
}

func (g *Gui) Init() error {
	gui, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		g.l.Error().Err(err).Msg("Could not initialize gui")
		return err
	}
	g.gui = gui
	g.gui.SetManagerFunc(g.layout)
	if err = g.setupKeyBindings(); err != nil {
		g.l.Error().Err(err).Msg("Could not setup key bindings")
		return err
	}
	g.view = components.NewView(g, g.theme)
	g.l.Info().Msg("Initialized")
	return nil
}

func (g *Gui) layout(gui *gocui.Gui) error {
	g.rw.RLock()
	defer g.rw.RUnlock()
	err := g.view.Draw(gui, g.state())
	if err != nil {
		log.Error().Err(err).Msg("while drawing gui")
		return err
	}
	return nil
}

func (g *Gui) setupKeyBindings() error {
	if err := g.gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, g.quit); err != nil {
		return err
	}

	keyAdapter := func(f func(s State)) func(gui *gocui.Gui, v *gocui.View) error {
		return func(gui *gocui.Gui, v *gocui.View) error {
			g.rw.Lock()
			defer g.rw.Unlock()
			f(g.state())
			return nil
		}
	}

	if err := g.gui.SetKeybinding("", gocui.KeyArrowLeft, gocui.ModNone, keyAdapter(func(s State) { s.LeftPressed() })); err != nil {
		return err
	}

	if err := g.gui.SetKeybinding("", gocui.KeyArrowRight, gocui.ModNone, keyAdapter(func(s State) { s.RightPressed() })); err != nil {
		return err
	}

	if err := g.gui.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, keyAdapter(func(s State) { s.UpPressed() })); err != nil {
		return err
	}

	if err := g.gui.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, keyAdapter(func(s State) { s.DownPressed() })); err != nil {
		return err
	}

	if err := g.gui.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, keyAdapter(func(s State) { s.EnterPressed() })); err != nil {
		return err
	}

	if err := g.gui.SetKeybinding("", 'q', gocui.ModNone, keyAdapter(func(s State) { s.BackPressed() })); err != nil {
		return err
	}

	if err := g.gui.SetKeybinding("", gocui.KeyDelete, gocui.ModNone, keyAdapter(func(s State) { s.BackPressed() })); err != nil {
		return err
	}

	if err := g.gui.SetKeybinding("", gocui.KeyEsc, gocui.ModNone, keyAdapter(func(s State) { s.BackPressed() })); err != nil {
		return err
	}

	return nil
}

func (g *Gui) quit(gui *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func (g *Gui) Sync() {
	g.gui.Update(func(g *gocui.Gui) error { return nil })
}

func (g *Gui) Run(s GetStateFunc) error {
	g.state = s
	defer func() {
		if r := recover(); r != nil {
			if g.cfg.Dev {
				g.l.Error().Interface("recovered", r).Msg("Unexpected panic while running application")
				_, _ = g.l.Write(debug.Stack())
			}
		}
	}()

	g.l.Debug().Msg("Running")
	if err := g.gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		g.l.Error().Err(err).Msg("Unexpected error while running application")
		return err
	}
	return nil
}

func (g *Gui) Close() {
	g.l.Debug().Msg("Closing")
	g.gui.Close()
}

func (g *Gui) Size() (int, int) {
	w, h := g.gui.Size()
	return w - 1, h - 1
}

func (g *Gui) Position() (float64, float64) {
	return 0, 0
}
