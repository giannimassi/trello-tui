package gui

import (
	"github.com/gdamore/tcell"
	"github.com/giannimassi/trello-tui/pkg/store"
	"github.com/rivo/tview"
)

// Header is a gui component displaying the title and description of the current board
type Header struct {
	*tview.Box
	state store.HeaderState
}

// NewHeader returns a new instance of Header
func NewHeader(state store.HeaderState) *Header {
	return &Header{
		Box:   tview.NewBox().SetBorder(true),
		state: state,
	}
}

// SetState updates the Header componer with the HeaderState
func (h *Header) SetState(state store.HeaderState) {
	h.state = state
}

// Draw re-implements the `tview.Primitive` interface Draw function
func (h *Header) Draw(screen tcell.Screen) {
	h.SetTitle(" " + h.state.HeaderTitle() + " ")
	h.Box.Draw(screen)
}
