package state

import "github.com/giannimassi/trello-tui/pkg/domain"

type boardLoadingOffline struct {
	boardLoading
	err error
}

func (b *boardLoadingOffline) HeaderTitle() string {
	return b.boardName + " - offline"
}

func (b *boardLoadingOffline) HeaderSubtitle() string {
	return "Could not load board.\nDetails:\n" + b.err.Error()
}

func (b *boardLoadingOffline) online(newBoard *domain.Board) board {
	online := &boardOnline{
		boardLoading: b.boardLoading,
	}
	return online.online(newBoard)
}

func (b *boardLoadingOffline) offline(err error) board {
	b.err = err
	return b
}
