package components

import (
	"github.com/fatih/color"
	"github.com/giannimassi/trello-tui/pkg/gui/state"
)

// TODO: really bad smell here, why is this interface so huge, need to re-think this context thingy
type View interface {
	Name() string
	Description() string
	ListsLen() int
	ListNameByIndex(idx int) (string, bool)
	ListCardsIds(idx int) []int
	CardNameByID(idx int) (string, bool)
	Errors() []error

	// Navigation
	IsListSelected(idx int) bool
	IsCardSelected(id int) bool
	IsBoardLoaded() bool
	IsBoardLoading() bool
	IsBoardNotFound() bool
	IsCardPopupOpen() bool
	FirstListIndex() int
	FirstCardIndex(idx int) int
}

type Commands interface {
	UpdateFirstListIndex(listsPerPage, totalLists int)
	UpdateFirstCardIndex(cardsPerPage int, cardIDs []int)
	MoveLeft()
	MoveRight()
	MoveUp()
	MoveDown()
	OpenCardPopup()
	CloseCardPopup()
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
		var errStr string
		if errs := v.Errors(); len(errs) != 0 {
			errStr = errs[len(errs)-1].Error()
		}

		return "  Could not find board \"" +
			v.Name() +
			"\" (" + errStr + "). Press ctrl + c to exit application."
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

func (v *Context) CardTitle(id int) string {
	if name, found := v.View.CardNameByID(id); found {
		return name
	}
	return ""
}

func (v *Context) CardPopupTitle() string {
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
