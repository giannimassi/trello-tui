package state

import (
	"time"

	"github.com/VojtechVitek/go-trello"
)

type CardState int

const (
	_ = iota
	CardHighlighted
)

type State struct {
	Updated time.Time

	Board trello.Board
	Lists []trello.List
	Cards []trello.Card

	ErrorList         []error
	SelectedBoard     string
	SelectedList      int
	SelectedCard      int
	SelectedCardState CardState

	loading bool
}

func NewState() *State {
	return &State{
		Updated:      time.Now(),
		SelectedCard: -1,
	}
}

func (s *State) Name() string {
	return s.Board.Name
}

func (s *State) Description() string {
	return s.Board.Desc
}

func (s *State) ListsLen() int {
	return len(s.Lists)
}

func (s *State) ListName(idx int) (string, bool) {
	if idx >= len(s.Lists) {
		return "", false
	}

	return s.Lists[idx].Name, true
}

func (s *State) ListCardsIds(idx int) []int {
	var ids []int
	if idx >= len(s.Lists) {
		return ids
	}

	for _, c := range s.Cards {
		if c.IdList == s.Lists[idx].Id {
			ids = append(ids, c.IdShort)
		}
	}

	return ids
}

func (s *State) CardName(cardID int) (string, bool) {
	for _, c := range s.Cards {
		if c.IdShort == cardID {
			return c.Name, true
		}
	}

	return "", false
}

func (s *State) Errors() []error {
	return s.ErrorList
}

func (s *State) Loading() bool {
	return s.loading
}

// Commands

func (s *State) Navigate(listID, cardID int) {
	s.SelectedList = listID
	s.SelectedCard = cardID
}

func (s *State) AppendErr(err error) {
	s.ErrorList = append(s.ErrorList, err)
}

func (s *State) SetLoading(loading bool) {
	s.loading = loading
}
