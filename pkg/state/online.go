package state

import (
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/giannimassi/trello-tui/pkg/domain"
)

type boardOnline struct {
	boardLoading
	Board *domain.Board
}

func (b *boardOnline) online(newBoard *domain.Board) board {
	b.Board = newBoard
	if !b.isSelectedCardValid() {
		b.selectFirstCardAvailable()
	}
	return b
}

func (b *boardOnline) offline(err error) board {
	offline := &boardOffline{
		boardOnline: *b,
		err:         err,
	}
	return offline
}

func (b *boardOnline) HeaderTitle() string {
	return b.boardName + " - online"
}

func (b *boardOnline) HeaderSubtitle() string {
	return b.Board.Description
}

func (b *boardOnline) ListName(idx int) string {
	if idx >= len(b.Board.Lists) {
		return ""
	}
	return b.Board.Lists[idx].Name
}

func (b *boardOnline) ListCardsIds(idx int) []int {
	if idx >= len(b.Board.Lists) {
		return []int{}
	}
	return b.Board.Lists[idx].CartIds
}

func (b *boardOnline) CardName(id int) string {
	c, found := b.Board.CardById(id)
	if !found {
		return ""
	}
	return fmt.Sprintf("%v - %s", c.Pos, c.Name)
}

func (b *boardOnline) SelectedCardName() string {
	card, found := b.Board.CardById(b.CardId)
	if !found {
		return ""
	}
	return card.Name
}

func (b *boardOnline) SelectedCardDescription() string {
	card, found := b.Board.CardById(b.CardId)
	if !found {
		return ""
	}
	return card.Description
}

func (b *boardOnline) IsCardPopupOpen() bool {
	return b.CardPopupOpen
}

func (b *boardOnline) ListsLen() int {
	return len(b.Board.Lists)
}

func (b *boardOnline) IsListSelected(idx int) bool {
	return b.ListIndex == idx
}

func (b *boardOnline) IsCardSelected(id int) bool {
	return b.CardId == id
}

func (b *boardOnline) isSelectedCardValid() bool {
	if _, found := b.Board.CardById(b.CardId); !found {
		return false
	}
	return true
}

func (b *boardOnline) selectFirstCardAvailable() {
	for listIndex := range b.Board.Lists {
		for _, cardID := range b.Board.Lists[listIndex].CartIds {
			b.ListIndex = listIndex
			b.CardId = cardID
			return
		}
	}
}

func (b *boardOnline) LeftPressed() {
	if b.Board.IsEmpty {
		return
	}
	// move to first in previous list (first in current list if its the first on)
	for i := b.ListIndex - 1; i >= 0; i-- {
		cardIDs := b.Board.Lists[i].CartIds
		if len(cardIDs) != 0 {
			b.ListIndex = i
			b.CardId = cardIDs[0]
			return
		}
	}
}

func (b *boardOnline) RightPressed() {
	if b.Board.IsEmpty {
		return
	}
	// move to first in next list (first in current list if its the last in board)
	for i := b.ListIndex + 1; i < len(b.Board.Lists); i++ {
		cardIDs := b.Board.Lists[i].CartIds
		if len(cardIDs) != 0 {
			b.ListIndex = i
			b.CardId = cardIDs[0]
			return
		}
	}
}

func (b *boardOnline) UpPressed() {
	if b.Board.IsEmpty {
		return
	}
	// Move to previous card in list or stop
	cardIDs := b.Board.Lists[b.ListIndex].CartIds
	if i, found := cardIndexInListFromID(cardIDs, b.CardId); i > 0 && found {
		b.CardId = cardIDs[i-1]
	}
	return
}

func (b *boardOnline) DownPressed() {
	if b.Board.IsEmpty {
		return
	}
	// move to next card or stop
	cardIDs := b.ListCardsIds(b.ListIndex)
	if i, found := cardIndexInListFromID(cardIDs, b.CardId); i+1 < len(cardIDs) && found {
		b.CardId = cardIDs[i+1]
	}

	return
}

func (b *boardOnline) EnterPressed() {
	b.CardPopupOpen = true
}

func (b *boardOnline) BackPressed() {
	log.Info().Msg("back ")
	b.CardPopupOpen = false
}

func (b *boardOnline) FirstVisibleListIndex(listsPerPage int) int {
	min := minIndex(b.ListIndex, listsPerPage)
	max := maxIndex(b.ListIndex, listsPerPage, len(b.Board.Lists))
	if b.firstVisibleListIndex < min {
		log.Info().Msg("changing first visible index <")
		b.firstVisibleListIndex = min
	} else if b.firstVisibleListIndex > max {
		log.Info().Msg("changing first visible index >")
		b.firstVisibleListIndex = max
	}
	return b.firstVisibleListIndex
}

func (b *boardOnline) FirstVisibleCardIndex(listIndex, cardsPerPage int) int {
	if len(b.firstVisibleCardIndex) == 0 {
		b.firstVisibleCardIndex = make([]int, len(b.Board.Lists))
	}
	cardIds := b.Board.Lists[listIndex].CartIds
	currentIndex := b.firstVisibleCardIndex[listIndex]
	selectedCardIndex, found := cardIndexInListFromID(cardIds, b.CardId)
	min := minIndex(selectedCardIndex, cardsPerPage)
	max := maxIndex(selectedCardIndex, cardsPerPage, len(cardIds))
	if !found {
		return b.firstVisibleCardIndex[listIndex]
	}
	if currentIndex < min {
		log.Info().Msg("changing first visible index <")
		b.firstVisibleCardIndex[listIndex] = min
	} else if currentIndex > max {
		log.Info().Msg("changing first visible index >")
		b.firstVisibleCardIndex[listIndex] = max
	}
	return b.firstVisibleCardIndex[listIndex]
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

func cardIndexInListFromID(cardIds []int, id int) (int, bool) {
	for i, cId := range cardIds {
		if cId == id {
			return i, true
		}
	}
	return 0, false
}
