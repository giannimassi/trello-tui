package state

import (
	"github.com/giannimassi/trello-tui/pkg/domain"
)

type selected struct {
	ListIndex             int
	CardId                int
	CardPopupOpen         bool
	firstVisibleListIndex int
	firstVisibleCardIndex []int
}

type boardLoading struct {
	selected
	boardName string
}

var _ board = &boardLoading{}

func (b *boardLoading) online(newBoard *domain.Board) board {
	onlineBoard := &boardOnline{
		boardLoading: *b,
	}
	return onlineBoard.online(newBoard)
}

func (b *boardLoading) offline(err error) board {
	offline := &boardLoadingOffline{
		boardLoading: *b,
		err:          err,
	}
	return offline
}

func (b *boardLoading) LeftPressed()                                    { return }
func (b *boardLoading) RightPressed()                                   { return }
func (b *boardLoading) UpPressed()                                      { return }
func (b *boardLoading) DownPressed()                                    { return }
func (b *boardLoading) EnterPressed()                                   { return }
func (b *boardLoading) BackPressed()                                    { return }
func (b *boardLoading) HeaderTitle() string                             { return b.boardName + " - loading" }
func (b *boardLoading) HeaderSubtitle() string                          { return "..." }
func (b *boardLoading) IsListSelected(idx int) bool                     { return false }
func (b *boardLoading) ListName(idx int) string                         { return "Blank" }
func (b *boardLoading) ListCardsIds(idx int) []int                      { return nil }
func (b *boardLoading) FirstVisibleCardIndex(idx, cardsPerPage int) int { return 0 }
func (b *boardLoading) CardName(id int) string                          { return "Blank" }
func (b *boardLoading) IsCardSelected(id int) bool                      { return false }
func (b *boardLoading) SelectedCardName() string                        { return "Blank" }
func (b *boardLoading) SelectedCardDescription() string                 { return "Blank" }
func (b *boardLoading) IsCardPopupOpen() bool                           { return false }
func (b *boardLoading) FirstVisibleListIndex(listsPerPage int) int      { return 0 }
func (b *boardLoading) ListsLen() int                                   { return 0 }
