package gui

import (
	"github.com/gdamore/tcell"
	"github.com/giannimassi/trello-tui/pkg/store"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

const (
	paddingSize     = 1
	titleSize       = 3
	labelsSize      = 3
	descriptionSize = 3
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
	root.AddItem(nil, 0, paddingSize, false)
	root.AddItem(innerF, 0, 4, false)
	// right padding
	root.AddItem(nil, 0, paddingSize, false)

	title := tview.NewTextView()
	title.SetBorder(true)

	labels := tview.NewTextView()
	labels.SetBorder(true)
	labels.SetDynamicColors(true)

	description := tview.NewTextView()
	description.SetBorder(true)
	description.SetInputCapture(c.captureInput)

	// top padding
	innerF.AddItem(nil, 0, paddingSize, false)
	innerF.AddItem(title, titleSize, 1, false)
	innerF.AddItem(labels, labelsSize, 1, false)
	innerF.AddItem(description, 0, descriptionSize, true)
	// bottom padding
	innerF.AddItem(nil, 0, paddingSize, false)
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

// CardView is in charge of returning navigation to list container
// if the id points to unexistent card id

// Draw re-implements the `tview.Primitive` interface Draw function
func (c *CardView) Draw(screen tcell.Screen) {
	c.title.SetText(c.state.CardName(c.id))
	c.labels.SetText(c.state.CardLabelsStr(c.id))
	c.description.SetText(c.state.Description(c.id))
	c.Flex.Draw(screen)
}

func (c *CardView) captureInput(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyEsc:
		c.handler.switchToListContainerView()
		return event
	case tcell.KeyEnter:
		return nil
	}

	log.Debug().Msg("captured input")
	return event
}
