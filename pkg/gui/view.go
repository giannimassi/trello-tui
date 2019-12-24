package gui

import (
	"github.com/giannimassi/trello-tui/pkg/store"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

// View is the root gui component
type View struct {
	*tview.Flex
	header        *Header
	listContainer *ListContainer
	card          *CardView

	focuser     focuser
	cardFocused bool
}

type focuser interface {
	SetFocus(p tview.Primitive)
}

const (
	headerHeight = 1
	bodyHeight   = 6
)

// NewView returns a new instance of View
func NewView(state store.ViewState, f focuser) *View {
	var (
		v = View{
			focuser: f,
		}
		header        = NewHeader(state)
		listContainer = NewListContainer(3, state, f, &v)
		card          = NewCardView(state, &v)
		flex          = tview.NewFlex().
				SetFullScreen(true).
				SetDirection(tview.FlexRow).
				AddItem(header, 0, headerHeight, false).
				AddItem(listContainer, 0, bodyHeight, false)
	)
	v.Flex = flex
	v.header = header
	v.listContainer = listContainer
	v.card = card
	return &v
}

// SetState updates the View with the ViewState
func (v *View) SetState(s store.ViewState) {
	v.header.SetState(s)
	v.listContainer.SetState(s)
	v.card.SetState(s)
}

// FocusedItem returns the gui component currently in focus
// TODO: this is mainly necessary to ensure the focus is on the correct
// item at startup, this is not nice.
func (v *View) FocusedItem() tview.Primitive {
	if v.cardFocused {
		return v.card.FocusedItem()
	}
	return v.listContainer.FocusedItem()
}

func (v *View) switchToCardView(id int) {
	log.Debug().Int("id", id).Msg("switching to card view")
	// remove list container
	v.RemoveItem(v.listContainer)
	// add card view
	v.card.id = id
	v.AddItem(v.card, 0, bodyHeight, true)
	v.cardFocused = true
	v.focuser.SetFocus(v.FocusedItem())
}

func (v *View) switchToListContainerView() {
	// remove list container
	v.RemoveItem(v.card)
	v.AddItem(v.listContainer, 0, bodyHeight, true)
	v.cardFocused = false
	v.focuser.SetFocus(v.FocusedItem())
}
