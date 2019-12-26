package trello

import (
	"github.com/VojtechVitek/go-trello"
	"github.com/giannimassi/trello-tui/pkg/store"
)

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

func (b *boardLoadingOffline) online(newBoard *trello.Board, lists []trello.List, cards []trello.Card) store.State {
	online := &boardOnline{
		boardLoading: b.boardLoading,
	}
	return online.online(newBoard, lists, cards)
}

func (b *boardLoadingOffline) offline(err error) *boardLoadingOffline {
	b.err = err
	return b
}
