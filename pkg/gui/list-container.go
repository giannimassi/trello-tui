package gui

import (
	"github.com/gdamore/tcell"
	"github.com/giannimassi/trello-tui/pkg/store"
	"github.com/rivo/tview"
)

type switcher interface {
	switchToCardView(id int)
}

// ListContainer is a gui component in charge of displaying the board's lists
type ListContainer struct {
	*tview.Flex

	focuser  focuser
	switcher switcher

	listV        []*ListView // listV is a collection of List gui components
	state        store.ListsState
	firstV, maxV int
	focusedV     int
	listALen     int
}

// NewListContainer returns a new instance of ListContainer
func NewListContainer(maxVLists int, state store.ListsState, f focuser, s switcher) *ListContainer {
	var (
		flex = tview.NewFlex().SetDirection(tview.FlexColumn)
		ls   = ListContainer{
			Flex:     flex,
			focuser:  f,
			switcher: s,
			maxV:     maxVLists,
			state:    state,
		}
	)
	for i := 0; i < ls.maxV; i++ {
		l := NewListView(&ls, state)
		flex.AddItem(l, 0, 1, i == 0)
		ls.listV = append(ls.listV, l)
	}
	return &ls
}

// Draw re-implements the `tview.Primitive` interface Draw function
func (l *ListContainer) Draw(screen tcell.Screen) {
	// log.Debug().Int("listsLen", l.state.ListsLen()).Int("focus", l.focusedV).Msg("ListContainer.Draw")
	l.listALen = l.state.ListsLen()
	// Check if there are less lists since last time state was set
	if l.focusedV >= l.listALen {
		if l.listALen == 0 {
			l.focusedV = 0
		} else {
			l.focusedV = l.listALen - 1
		}
	}
	// Point to right list based on navigation
	for i := range l.listV {
		l.listV[i].index = l.firstV + i
	}
	l.Flex.Draw(screen)
}

func (l *ListContainer) handleSelectPreviousList() {
	if l.focusedV != 0 {
		l.focusedV--
		l.focuser.SetFocus(l.listV[l.focusedV])
	} else if l.firstV > 0 {
		l.firstV--
	}
}

func (l *ListContainer) handleSelectNextList() {
	if l.focusedV < l.maxV-1 {
		l.focusedV++
		l.focuser.SetFocus(l.listV[l.focusedV])
	} else if l.firstV+l.maxV+1 < l.listALen {
		l.firstV++
	}
}

func (l *ListContainer) handleCardSelected(id int) {
	l.switcher.switchToCardView(id)
}

// FocusedItem returns the currently focused list
func (l *ListContainer) FocusedItem() tview.Primitive {
	if len(l.listV) == 0 {
		return nil
	}
	return l.listV[l.focusedV]
}
