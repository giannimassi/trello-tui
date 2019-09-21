package state

type NavigationPosition struct {
	SelectedBoard     string
	SelectedListIndex int
	SelectedCardID    int
	FirstListIdx      int
	FirstCardIdxs     []int // TODO[gianni] this should not be state, or it should be on a different layer of state (visualization depends on navigation, but can be implemented  independently in gui if not required to be persisted)
	SelectedCardState
	BoardState
}

// View
func (n *NavigationPosition) IsListSelected(idx int) bool {
	return n.SelectedListIndex == idx
}

func (n *NavigationPosition) IsCardSelected(id int) bool {
	return n.SelectedCardID == id
}

func (n *NavigationPosition) IsBoardLoading() bool {
	return n.BoardState == BoardLoading
}

func (n *NavigationPosition) IsBoardLoaded() bool {
	return n.BoardState == BoardLoaded
}

func (n *NavigationPosition) IsBoardNotFound() bool {
	return n.BoardState == BoardNotFound
}

func (n *NavigationPosition) IsCardPopupOpen() bool {
	return n.SelectedCardState == CardPopupOpen
}

func (n *NavigationPosition) FirstListIndex() int {
	return n.FirstListIdx
}

func (n *NavigationPosition) FirstCardIndex(idx int) int {
	return n.FirstCardIdxs[idx]
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
