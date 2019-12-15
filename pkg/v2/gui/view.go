package gui

import (
	"github.com/rivo/tview"
)

// View is the root gui component
type View struct {
	app  *tview.Application // TODO: replace with interface only exposing required methods
	flex *tview.Flex

	header        *Header
	listContainer *ListContainer
}

// NewView returns a new instance of View
func NewView(a *tview.Application) View {
	return View{
		app: a,
	}
}

// Draw initializes the View with the ViewState
func (v *View) Draw(s ViewState) {
	if v.flex == nil {
		v.flex = tview.NewFlex().SetFullScreen(true)
		v.app.SetRoot(v.flex, true)
	}

	if v.header == nil {
		v.header = NewHeader(v.app)
		v.flex.AddItem(v.header.box, 0, 1, false)
	}
	v.header.Draw(s)

	if v.listContainer == nil {
		v.listContainer = NewListContainer(v.app)
		v.flex.AddItem(v.listContainer.grid, 1, 4, true)
	}
	v.listContainer.Draw(s)
}

// Header is a gui component displaying the title and description of the current board
type Header struct {
	app *tview.Application // TODO: replace with interface only exposing required methods
	box *tview.Box
}

// NewHeader returns a new instance of Header
func NewHeader(a *tview.Application) *Header {
	return &Header{
		app: a,
		box: tview.NewBox().SetBorder(true),
	}
}

// Draw initializes the Header componer with the HeaderState
func (h *Header) Draw(state HeaderState) {
	h.box.SetTitle(state.HeaderTitle())
}

// ListContainer is a gui component in charge of displaying the board's lists
type ListContainer struct {
	app  *tview.Application // TODO: replace with interface only exposing required methods
	grid *tview.Grid
}

// NewListContainer returns a new instance of ListContainer
func NewListContainer(a *tview.Application) *ListContainer {
	return &ListContainer{
		app:  a,
		grid: tview.NewGrid(),
	}
}

// Draw initializes the ListContainer component with the ListState
func (l *ListContainer) Draw(state ListState) {

}
