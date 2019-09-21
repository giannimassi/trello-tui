package state

import (
	"time"

	"github.com/VojtechVitek/go-trello"
	"github.com/rs/zerolog/log"
)

type SelectedCardState int
type BoardState int

const (

	// BoardState
	BoardUninitialized = iota
	BoardLoading
	BoardNotFound
	BoardLoaded

	// SelectedCardState
	CardUninitialized = iota
	CardLoading
	CardLoaded
	CardPopupOpen
)

type State struct {
	Updated time.Time

	Board trello.Board
	Lists []trello.List
	Cards []trello.Card

	errors []error
	NavigationPosition
}

func NewState() *State {
	return &State{
		Updated: time.Now(),
		NavigationPosition: NavigationPosition{
			SelectedCardID: -1,
		},
	}
}

func (s *State) Name() string {
	switch s.BoardState {
	case BoardLoaded:
		return s.Board.Name
	case BoardLoading:
		return s.NavigationPosition.SelectedBoard
	default:
		return ""
	}
}

func (s *State) Description() string {
	if s.BoardState != BoardLoaded {
		return ""
	}
	return s.Board.Desc
}

func (s *State) ListsLen() int {
	if s.BoardState != BoardLoaded {
		return 0
	}
	return len(s.Lists)
}

func (s *State) ListNameByIndex(idx int) (string, bool) {
	if idx >= len(s.Lists) || s.BoardState != BoardLoaded {
		return "", false
	}

	return s.Lists[idx].Name, true
}

func (s *State) ListCardsIds(idx int) []int {
	var ids []int
	if idx >= len(s.Lists) || s.BoardState != BoardLoaded {
		return ids
	}

	for _, c := range s.Cards {
		if c.IdList == s.Lists[idx].Id {
			ids = append(ids, c.IdShort)
		}
	}

	return ids
}

func (s *State) CardNameByID(cardID int) (string, bool) {
	if s.BoardState != BoardLoaded {
		return "", false
	}

	for _, c := range s.Cards {
		if c.IdShort == cardID {
			return c.Name, true
		}
	}

	return "", false
}

type Card struct {
	trello.Card
}

func (s *State) SelectedCard() Card {
	for _, c := range s.Cards {
		if c.IdShort == s.SelectedCardID {
			return Card{c}
		}
	}

	return Card{}
}

func (s *State) Errors() []error {
	return s.errors
}

// Commands
func (s *State) InitNavigation() {
	if s.NavigationPosition.isInitialized() {
		return
	}
	s.NavigationPosition.FirstCardIdxs = make([]int, len(s.Lists))
	s.NavigationPosition.selectFirstCardAvailable(s)
}

func (s *State) SetBoardState(boardState BoardState) {
	prevState := s.BoardState
	s.Updated = time.Now()
	switch boardState {
	case BoardUninitialized:
		s.Board = trello.Board{}
		s.Lists = []trello.List{}
		s.Cards = []trello.Card{}
		s.NavigationPosition.SelectedListIndex = -1
		s.NavigationPosition.SelectedCardID = -1
		s.NavigationPosition.SelectedCardState = CardUninitialized
	case BoardLoading:
		s.Board = trello.Board{}
		s.Lists = []trello.List{}
		s.Cards = []trello.Card{}
		s.NavigationPosition.SelectedListIndex = -1
		s.NavigationPosition.SelectedCardID = -1
		s.NavigationPosition.SelectedCardState = CardLoading
	case BoardNotFound:
		s.Board = trello.Board{}
		s.Lists = []trello.List{}
		s.Cards = []trello.Card{}
		s.NavigationPosition.SelectedListIndex = -1
		s.NavigationPosition.SelectedCardID = -1
		s.NavigationPosition.SelectedCardState = CardUninitialized
	case BoardLoaded:
		s.NavigationPosition.SelectedCardState = CardLoaded
	}
	s.BoardState = boardState
	log.Debug().Int("old", int(prevState)).Int("new", int(s.BoardState)).Msg("board state changed")
}

func (s *State) AppendErr(err error) {
	s.errors = append(s.errors, err)
}

func (s *State) MoveLeft() {
	if len(s.Lists) == 0 || len(s.Cards) < 1 {
		return
	}
	// move to first in previous list (first in current list if its the first on)
	for i := s.SelectedListIndex - 1; i >= 0; i-- {
		cardIDs := s.ListCardsIds(i)
		if len(cardIDs) != 0 {
			s.SelectedListIndex = i
			s.SelectedCardID = cardIDs[0]
			break
		}
	}
}

func (s *State) MoveRight() {
	if len(s.Lists) == 0 || len(s.Cards) < 1 {
		return
	}
	// move to first in next list (first in current list if its the last in board)
	for i := s.SelectedListIndex + 1; i < len(s.Lists); i++ {
		cardIDs := s.ListCardsIds(i)
		log.Print(cardIDs, i, s)
		if len(cardIDs) != 0 {
			s.SelectedListIndex = i
			s.SelectedCardID = cardIDs[0]
			break
		}
	}
}

func (s *State) MoveUp() {
	if len(s.Lists) == 0 || len(s.Cards) < 1 {
		return
	}
	// Move to previous card in list or stop
	cardIDs := s.ListCardsIds(s.SelectedListIndex)
	if i := cardIndexInListFromID(cardIDs, s.SelectedCardID); i > 0 {
		s.SelectedCardID = cardIDs[i-1]
	}
}

func (s *State) MoveDown() {
	if len(s.Lists) == 0 || len(s.Cards) < 1 {
		return
	}
	// move to next card or stop
	cardIDs := s.ListCardsIds(s.SelectedListIndex)
	if i := cardIndexInListFromID(cardIDs, s.SelectedCardID); i+1 < len(cardIDs) {
		s.SelectedCardID = cardIDs[i+1]
	}
}

func (s *State) OpenCardPopup() {
	if s.SelectedCardState == CardLoaded {
		s.SelectedCardState = CardPopupOpen
	}
}

func (s *State) CloseCardPopup() {
	if s.SelectedCardState == CardPopupOpen {
		s.SelectedCardState = CardLoaded
	}
}
