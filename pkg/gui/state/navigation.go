package state

import (
	"github.com/jroimartin/gocui"
	"github.com/rs/zerolog/log"
)

type NavigationPosition struct {
	SelectedBoard     string
	SelectedListIndex int
	SelectedCardID    int
	SelectedCardState CardState
	FirstListIdx      int
	FirstCardIdxs     []int
}

// View
func (n *NavigationPosition) IsListSelected(idx int) bool {
	return n.SelectedListIndex == idx
}

func (n *NavigationPosition) IsCardSelected(id int) bool {
	return n.SelectedCardID == id
}

func (n *NavigationPosition) isInitialized() bool {
	return n.SelectedListIndex != -1 && n.SelectedCardID != -1
}

// Commands
func (n *NavigationPosition) selectFirstCardAvailable(s *State) {
	if len(s.Lists) != 0 {
		for i := 0; i < len(s.Lists); i++ {
			cardIDs := s.ListCardsIds(i)
			if len(cardIDs) != 0 {
				s.NavigationPosition.SelectedListIndex = i
				s.NavigationPosition.SelectedCardID = cardIDs[0]
				break
			}
		}
	}
}

func (n *NavigationPosition) update(s *State, k gocui.Key) {
	if len(s.Lists) == 0 || len(s.Cards) < 2 {
		return
	}
	switch k {
	case gocui.KeyArrowLeft:
		// move to first in previous list (first in current list if its the first on)
		for i := n.SelectedListIndex - 1; i >= 0; i-- {
			cardIDs := s.ListCardsIds(i)
			if len(cardIDs) != 0 {
				n.SelectedListIndex = i
				n.SelectedCardID = cardIDs[0]
				break
			}
		}

	case gocui.KeyArrowRight:
		// move to first in next list (first in current list if its the last in board)
		for i := n.SelectedListIndex + 1; i < len(s.Lists); i++ {
			cardIDs := s.ListCardsIds(i)
			log.Print(cardIDs, i, n)
			if len(cardIDs) != 0 {
				n.SelectedListIndex = i
				n.SelectedCardID = cardIDs[0]
				break
			}
		}

	case gocui.KeyArrowUp:
		// Move to previous card in list or stop
		cardIDs := s.ListCardsIds(n.SelectedListIndex)
		if i := cardIndexInListFromID(cardIDs, n.SelectedCardID); i > 0 {
			n.SelectedCardID = cardIDs[i-1]
		}

	case gocui.KeyArrowDown:
		// move to next card or stop
		cardIDs := s.ListCardsIds(n.SelectedListIndex)
		if i := cardIndexInListFromID(cardIDs, n.SelectedCardID); i+1 < len(cardIDs) {
			n.SelectedCardID = cardIDs[i+1]
		}
	}
}

func (n *NavigationPosition) FirstListIndex() int {
	return n.FirstListIdx
}

func (n *NavigationPosition) FirstCardIndex(idx int) int {
	return n.FirstCardIdxs[idx]
}

func (n *NavigationPosition) UpdateFirstListIndex(listsPerPage, totalLists int) {
	min := minIndex(n.SelectedListIndex, listsPerPage)
	max := maxIndex(n.SelectedListIndex, listsPerPage, totalLists)
	if n.FirstListIdx < min {
		n.FirstListIdx = min
	} else if n.FirstListIdx > max {
		n.FirstListIdx = max
	}
}

func (n *NavigationPosition) UpdateFirstCardIndex(cardsPerPage int, cardIDs []int) {
	currentIdx := n.FirstCardIdxs[n.SelectedListIndex]
	min := minIndex(cardIndexInListFromID(cardIDs, n.SelectedCardID), cardsPerPage)
	max := maxIndex(cardIndexInListFromID(cardIDs, n.SelectedCardID), cardsPerPage, len(cardIDs))
	if currentIdx < min {
		n.FirstCardIdxs[n.SelectedListIndex] = min
	} else if currentIdx > max {
		n.FirstCardIdxs[n.SelectedListIndex] = max
	}
}

func minIndex(selected, perPage int) int {
	if selected-perPage+1 > 0 {
		return selected - perPage + 1
	}
	return 0
}

func maxIndex(selected, perPage, total int) int {
	if selected+perPage-1 < total {
		return selected
	}
	if total-perPage < 0 {
		return 0
	}
	return total - perPage
}

func cardIndexInListFromID(cardIds []int, id int) int {
	for i, cId := range cardIds {
		if cId == id {
			return i
		}
	}
	return 0
}
