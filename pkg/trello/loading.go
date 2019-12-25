package trello

import (
	"github.com/giannimassi/trello-tui/pkg/domain"
)

type boardLoading struct {
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

func (b *boardLoading) HeaderTitle() string          { return b.boardName + " - loading" }
func (b *boardLoading) HeaderSubtitle() string       { return "..." }
func (b *boardLoading) ListName(idx int) string      { return "Loading..." }
func (b *boardLoading) ListCardsIds(idx int) []int   { return nil }
func (b *boardLoading) CardName(id int) string       { return "" }
func (b *boardLoading) CardLabelsStr(id int) string  { return "" }
func (b *boardLoading) Description(id int) string    { return "" }
func (b *boardLoading) CardComments(id int) []string { return []string{} }
func (b *boardLoading) ListsLen() int                { return 0 }
