package gui

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Data interface {
	Name() string
	Description() string
	ListsNames() []string
}

type Config struct{}

type App struct {
	l   zerolog.Logger
	cfg *Config

	b Data
	g *gocui.Gui
}

func NewApp(log zerolog.Logger, cfg *Config, data Data) *App {
	return &App{
		l:   log,
		cfg: cfg,
		b:   data,
	}
}

func (a *App) Init() error {
	a.l.Debug().Msg("Initializing gui")
	gui := gocui.NewGui()
	if err := gui.Init(); err != nil {
		a.l.Error().Err(err).Msg("Could not initialize gui")
		return err
	}
	a.g = gui
	a.g.SetLayout(a.layout)

	if err := a.g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, a.quit); err != nil {
		a.l.Error().Err(err).Msg("Could not setup key binding for ctrl + c")
		return err
	}
	a.l.Debug().Msg("Gui initialized")
	return nil
}

func (a *App) Run() error {
	if err := a.g.MainLoop(); err != nil && err != gocui.ErrQuit {
		a.l.Error().Err(err).Msg("Unexpected error while running application")
		return err
	}
	return nil
}

func (a *App) Close() {
	a.g.Close()
}

func (a *App) layout(g *gocui.Gui) error {
	var (
		maxX, maxY = a.g.Size()
		W          = maxX - 2
		H          = maxY - 2
		x0         = 1
		y0         = 1
	)

	w, err := a.drawLists(x0, y0+5, W, H-5)
	if err != nil {
		log.Warn().Err(err).Msg("while drawing lists")
		return err
	}

	err = a.drawMainTitle(x0, y0, w, 3)
	if err != nil {
		log.Warn().Err(err).Msg("while setting main view title")
		return err
	}

	return nil
}

func (a *App) quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func (a *App) drawMainTitle(X0, Y0, W, H int) error {
	if v, err := a.g.SetView("main", X0, Y0, X0+W, Y0+H); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		var descStr string
		if desc := a.b.Description(); len(desc) != 0 {
			descStr = fmt.Sprintf(", Description: %s", a.b.Description())
		}
		_, err := fmt.Fprintf(v, " Board: %s"+descStr, a.b.Name())
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *App) drawLists(X0, Y0, W, H int) (int, error) {
	var (
		listNames = a.b.ListsNames()
		w         = W / len(listNames)
	)
	for i, name := range listNames {
		err := a.drawList(X0+i*w, Y0, w, H, name)
		if err != nil {
			return W, err
		}
	}
	return w * len(listNames), nil
}

func (a *App) drawList(X0, Y0, W, H int, name string) error {
	if v, err := a.g.SetView(listViewName(name), X0, Y0, X0+W, Y0+H); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		_, err := fmt.Fprint(v, name)
		if err != nil {
			return err
		}
	}
	return nil
}

func listViewName(listName string) string {
	return fmt.Sprintf("list_%s", strings.ReplaceAll(listName, " ", "_"))
}
