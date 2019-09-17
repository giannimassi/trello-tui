package state

import (
	"time"

	"github.com/VojtechVitek/go-trello"
	"github.com/jroimartin/gocui"
	"github.com/rs/zerolog/log"
)

type CardState int
type BoardState int

const (

	// BoardState
	BoardUninitialized = iota
	BoardLoading
	BoardNotFound
	BoardLoaded

	// CardState
	CardUninitialized = iota
	CardLoading
	CardNotFound
	CardLoaded
	CardSelected
)

type State struct {
	Updated time.Time

	Board trello.Board
	Lists []trello.List
	Cards []trello.Card

	errors []error
	NavigationPosition
	BoardState BoardState
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

func (s *State) Errors() []error {
	return s.errors
}

func (s *State) IsBoardLoading() bool {
	return s.BoardState == BoardLoading
}

func (s *State) IsBoardLoaded() bool {
	return s.BoardState == BoardLoaded
}

func (s *State) IsBoardNotFound() bool {
	return s.BoardState == BoardNotFound
}

func (s *State) InitNavigation() {
	if s.NavigationPosition.isInitialized() {
		return
	}
	s.NavigationPosition.FirstCardIdxs = make([]int, len(s.Lists))
	s.NavigationPosition.selectFirstCardAvailable(s)
}

// Commands
func (s *State) KeyPressed(k gocui.Key, m gocui.Modifier) {
	switch k {
	case gocui.KeyArrowLeft, gocui.KeyArrowRight, gocui.KeyArrowUp, gocui.KeyArrowDown:
		s.NavigationPosition.update(s, k)
	case gocui.KeyEnter:
		log.Warn().Msg("Enter not implemented")
	case gocui.KeyEsc:
		log.Warn().Msg("Esc not implemented")
	}
}

func (s *State) AppendErr(err error) {
	s.errors = append(s.errors, err)
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
		s.NavigationPosition.SelectedCardState = CardNotFound
	case BoardLoaded:
		s.NavigationPosition.SelectedCardState = CardLoaded
	}
	s.BoardState = boardState
	log.Debug().Int("old", int(prevState)).Int("new", int(s.BoardState)).Msg("board state changed")
}
