package gui

import (
	"github.com/gdamore/tcell"
	"github.com/giannimassi/trello-tui/pkg/store"
	"github.com/rivo/tview"
)

const (
	paddingRelSize     = 1
	titleFixedSize     = 3
	labelsFixedSize    = 3
	descriptionRelSize = 3
)

type cardInputHandler interface {
	switchToListContainerView()
}

// CardView is a gui component in charge of displaying an open card and the list it belongs to
type CardView struct {
	*tview.Flex
	title       *tview.TextView
	labels      *tview.TextView
	description *tview.TextView

	id      int
	handler cardInputHandler
	state   store.CardState
}

// NewCardView returns an new instance of CardView
func NewCardView(state store.CardState, handler cardInputHandler) *CardView {
	c := CardView{
		id:      -1,
		state:   state,
		handler: handler,
	}
	root := tview.NewFlex()
	root.SetInputCapture(c.captureInput)
	innerF := tview.NewFlex().SetDirection(tview.FlexRow)
	// left padding
	root.AddItem(nil, 0, paddingRelSize, false)
	root.AddItem(innerF, 0, 4, false)
	// right padding
	root.AddItem(nil, 0, paddingRelSize, false)

	title := tview.NewTextView()
	title.SetBorder(true)

	labels := tview.NewTextView()
	labels.SetBorder(true)
	labels.SetDynamicColors(true)

	description := tview.NewTextView()
	description.SetBorder(true)
	description.SetInputCapture(c.captureInput)

	// top padding
	innerF.AddItem(nil, 0, paddingRelSize, false)
	innerF.AddItem(title, titleFixedSize, 1, false)
	innerF.AddItem(labels, labelsFixedSize, 1, false)
	innerF.AddItem(description, 0, descriptionRelSize, true)
	// bottom padding
	innerF.AddItem(nil, 0, paddingRelSize, false)
	c.Flex = root
	c.title = title
	c.labels = labels
	c.description = description
	return &c
}

// FocusedItem returns the gui component currently in focus
func (c *CardView) FocusedItem() tview.Primitive {
	return c.description
}

// Draw re-implements the `tview.Primitive` interface Draw function
func (c *CardView) Draw(screen tcell.Screen) {
	// TODO: CardView should be in charge of returning navigation
	// to list container if the id points to unexistent card id
	c.title.SetText(c.state.CardName(c.id))
	c.labels.SetText(c.state.CardLabelsStr(c.id))
	c.description.SetText(c.state.Description(c.id))
	c.Flex.Draw(screen)
}

func (c *CardView) captureInput(event *tcell.EventKey) *tcell.EventKey {
	if event.Rune() == 'q' {
		c.handler.switchToListContainerView()
		return event
	}

	switch event.Key() {
	case tcell.KeyEsc:
		c.handler.switchToListContainerView()
		return event
	case tcell.KeyEnter:
		return nil
	}
	return event
}
