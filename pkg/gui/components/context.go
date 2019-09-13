package components

import (
	"github.com/fatih/color"
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
	NavPosition() state.NavigationPosition
	IsBoardLoaded() bool
	IsBoardLoading() bool
	IsBoardNotFound() bool
	Errors() []error
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
	if v.IsBoardNotFound() {
		return " Board not found "
	}

	if !v.IsBoardLoaded() {
		return " Loading " + v.Name() + "... "
	}
	return " Board: " + v.Name() + " "
}

func (v *Context) HeaderSubtitle() string {
	if v.IsBoardNotFound() {
		return "  Could not find board \"" + v.NavPosition().SelectedBoard + "\""
	}

	if v.IsBoardLoading() {
		return "  Description loading..."
	}

	if !v.HasDescription() || !v.IsBoardLoaded() {
		return ""
	}

	return "  " + v.Description()
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

func (v *Context) Color(t ElementClass, isSelected bool) *color.Color {
	setting, found := DefaultColorSettings[t]
	if !found {
		setting = DefaultColorSettings[DefaultClass]
	}
	if isSelected {
		return setting.selected
	}
	return setting.normal
}
