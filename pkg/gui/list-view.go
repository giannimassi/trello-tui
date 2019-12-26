package gui

import (
	"github.com/gdamore/tcell"
	"github.com/giannimassi/trello-tui/pkg/store"
	"github.com/rivo/tview"
)

type listInputHandler interface {
	handleSelectPreviousList()
	handleSelectNextList()
	handleCardSelected(selectedID int)
}

// ListView is a gui component in charge of displaying a single board list
type ListView struct {
	parent listInputHandler
	*tview.Frame
	list     *tview.List
	index    int
	state    store.SingleListState
	hasFocus bool
}

// NewListView returns a new instance of ListView
func NewListView(parent listInputHandler, state store.SingleListState) *ListView {
	listView := ListView{
		parent: parent,
		state:  state,
	}
	ls := tview.NewList()
	ls.SetSelectedFocusOnly(true)
	ls.SetShortcutColor(tcell.ColorBlack)
	ls.SetInputCapture(listView.captureInput)
	ls.SetHighlightFullLine(true)
	ls.SetSelectedFunc(listView.handleSelected)
	listView.list = ls
	f := tview.NewFrame(ls)
	f.SetBorder(true)
	listView.Frame = f
	return &listView
}

// Draw re-implements the `tview.Primitive` interface Draw function
func (l *ListView) Draw(screen tcell.Screen) {
	l.SetTitle(" " + l.state.ListName(l.index) + " ")
	l.updateListItems(l.state.ListCardsIds(l.index))
	l.Frame.Draw(screen)
}

func (l *ListView) updateListItems(cardIds []int) {
	for i, id := range cardIds {
		cardName := l.state.CardName(id)
		cardLabels := l.state.CardLabelsStr(id) + "\n\n"
		// Add new list items
		if i >= l.list.GetItemCount() {
			l.list.AddItem(cardName, " ", ' ', nil)
		}
		// Update existing list items
		if oldTitle, oldLbls := l.list.GetItemText(i); oldTitle != cardName || oldLbls != cardLabels {
			l.list.SetItemText(i, cardName, cardLabels)
		}
		// Remove deleted list items
		for i := l.list.GetItemCount() - 1; i >= len(cardIds); i-- {
			l.list.RemoveItem(i)
		}
	}
}

func (l *ListView) captureInput(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	// intercept and pass to parent (list container):
	// - left/right keys: navigate to next or previous list
	case tcell.KeyLeft:
		l.parent.handleSelectPreviousList()
		return nil
	case tcell.KeyRight:
		l.parent.handleSelectNextList()
		return nil
	}
	// let default handler of the handle all other keys as well for now
	return event
}

func (l *ListView) handleSelected(index int, _, _ string, _ rune) {
	l.parent.handleCardSelected(l.selectedIndexToID(index))
}

func (l *ListView) selectedID() int {
	current := l.list.GetCurrentItem()
	cardIDs := l.state.ListCardsIds(l.index)
	if current >= len(cardIDs) {
		return -1
	}
	return cardIDs[current]
}

func (l *ListView) selectedIndexToID(index int) int {
	cardIDs := l.state.ListCardsIds(l.index)
	if index >= len(cardIDs) {
		return -1
	}
	return cardIDs[index]
}
