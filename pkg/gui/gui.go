package gui

import (
	"runtime/debug"

	"github.com/giannimassi/trello-tui/pkg/gui/components"
	"github.com/giannimassi/trello-tui/pkg/gui/state"
	"github.com/hashicorp/go-multierror"
	"github.com/jroimartin/gocui"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Dev bool // Enable developer features (recover on run panics)
}

type StateFunc func() *state.State

type Gui struct {
	l   zerolog.Logger
	cfg *Config

	stateFunc StateFunc

	gui           *gocui.Gui
	ctx           *components.Context
	viewContainer components.Container
}

func NewGui(log zerolog.Logger, cfg *Config) *Gui {
	return &Gui{
		l:   log,
		cfg: cfg,
	}
}

func (g *Gui) Init(stateFunc StateFunc) error {
	g.stateFunc = stateFunc
	g.ctx = components.NewGuiContext(g.stateFunc())
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
	g.viewContainer = components.NewContainer(g)
	g.l.Info().Msg("Initialized")
	return nil
}

func (g *Gui) layout(gui *gocui.Gui) error {
	g.ctx.Set(g.stateFunc())
	err := g.viewContainer.Draw(gui, g.ctx)
	if err != nil {
		log.Error().Err(err).Msg("while drawing gui")
		return err
	}
	return nil
}

func (g *Gui) Sync() {
	g.gui.Update(func(g *gocui.Gui) error { return nil })
}

func (g *Gui) Run() error {
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

func (g *Gui) setupKeyBindings() error {
	if err := g.gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, g.quit); err != nil {
		return errors.Wrapf(err, "while setting up key binding ctrl + c")
	}

	if err := g.setupNavigationKeybindings(); err != nil {
		return errors.Wrapf(err, "while setting up navigation key bindings")
	}

	return nil
}

func (g *Gui) setupNavigationKeybindings() error {
	var result *multierror.Error
	keyAdapter := func(f func()) func(gui *gocui.Gui, v *gocui.View) error {
		return func(gui *gocui.Gui, v *gocui.View) error {
			f()
			return nil
		}
	}

	if err := g.gui.SetKeybinding("", gocui.KeyArrowLeft, gocui.ModNone, keyAdapter(g.ctx.MoveLeft)); err != nil {
		result = multierror.Append(result, err)
	}

	if err := g.gui.SetKeybinding("", gocui.KeyArrowRight, gocui.ModNone, keyAdapter(g.ctx.MoveRight)); err != nil {
		result = multierror.Append(result, err)
	}

	if err := g.gui.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, keyAdapter(g.ctx.MoveUp)); err != nil {
		result = multierror.Append(result, err)
	}

	if err := g.gui.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, keyAdapter(g.ctx.MoveDown)); err != nil {
		result = multierror.Append(result, err)
	}

	if err := g.gui.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, keyAdapter(g.ctx.OpenCardPopup)); err != nil {
		result = multierror.Append(result, err)
	}

	if err := g.gui.SetKeybinding("", 'q', gocui.ModNone, keyAdapter(g.ctx.CloseCardPopup)); err != nil {
		result = multierror.Append(result, err)
	}

	if err := g.gui.SetKeybinding("", gocui.KeyEsc, gocui.ModNone, keyAdapter(g.ctx.CloseCardPopup)); err != nil {
		result = multierror.Append(result, err)
	}

	if result != nil {
		return result.ErrorOrNil()
	}

	return nil
}

func (g *Gui) quit(gui *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
