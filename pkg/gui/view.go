package gui

import (
	"github.com/giannimassi/trello-tui/pkg/store"
	"github.com/rivo/tview"
)

// View is the root gui component
type View struct {
	*tview.Flex
	header        *Header
	listContainer *ListContainer
	card          *CardView

	focuser     focuser
	cardFocused bool

	state store.State
}

type focuser interface {
	SetFocus(p tview.Primitive)
}

const (
	headerHeight = 1
	bodyHeight   = 6
)

// NewView returns a new instance of View
func NewView(s store.State, f focuser) *View {
	var (
		v = View{
			focuser: f,
		}
		header        = NewHeader(s)
		listContainer = NewListContainer(3, s, f, &v)
		card          = NewCardView(s, &v)
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
	v.state = s
	return &v
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
