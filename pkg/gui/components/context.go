package components

import (
	"github.com/giannimassi/trello-tui/pkg/gui/state"
	"github.com/jroimartin/gocui"
)

type View interface {
	Name() string
	Description() string
	ListsLen() int
	ListNameByIndex(idx int) (string, bool)
	ListCardsIds(idx int) []int
	CardNameByID(idx int) (string, bool)
	Errors() []error
	Loading() bool
}

type Commands interface {
	KeyPressed(k gocui.Key, m gocui.Modifier)
}

type Context struct {
	View
	Commands
}

func NewGuiContext(s *state.State) *Context {
	return &Context{
		View:     s,
		Commands: s,
	}
}

func (v *Context) Set(s *state.State) {
	v.View = s
	v.Commands = s
}

func (v *Context) HasDescription() bool {
	return len(v.View.Description()) != 0
}

func (v *Context) HeaderTitle() string {
	if v.Loading() {
		return " Loading " + v.Name() + "... "
	}
	return " Board: " + v.Name() + " "
}

func (v *Context) HeaderSubtitle() string {
	if v.Loading() {
		return ""
	}
	return " Description: " + v.Description() + " "
}

func (v *Context) ListTitle(idx int) string {
	if name, found := v.View.ListNameByIndex(idx); found {
		return " " + name + " "
	}
	return ""
}

func (v *Context) CardTitle(idx int) string {
	if name, found := v.View.CardNameByID(idx); found {
		return name
	}
	return ""
}
