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

	ErrorList []error
	Nav       NavigationPosition
	BoardState BoardState
}

func NewState() *State {
	return &State{
		Updated: time.Now(),
		Nav: NavigationPosition{
			SelectedCardID: -1,
		},
	}
}

func (s *State) Name() string {
	switch s.BoardState {
	case BoardLoaded:
		return s.Board.Name
	case BoardLoading:
		return s.Nav.SelectedBoard
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
	return s.ErrorList
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

func (s *State) NavPosition() NavigationPosition {
	return s.Nav
}

type NavigationPosition struct {
	SelectedBoard     string
	SelectedListIndex int
	SelectedCardID    int
	SelectedCardState CardState
}

func (n *NavigationPosition) IsListSelected(idx int) bool {
	return n.SelectedListIndex == idx
}

func (n *NavigationPosition) IsCardSelected(id int) bool {
	return n.SelectedCardID == id
}

// Commands

func (s *State) KeyPressed(k gocui.Key, m gocui.Modifier) {
	switch k {
	case gocui.KeyArrowLeft:
		s.moveInBoard(-1)
	case gocui.KeyArrowRight:
		s.moveInBoard(1)
	case gocui.KeyArrowUp:
		s.moveInList(-1)
	case gocui.KeyArrowDown:
		s.moveInList(1)
	case gocui.KeyEnter:
		log.Warn().Msg("Enter not implemented")
	case gocui.KeyEsc:
		log.Warn().Msg("Esc not implemented")
	}
}

func (s *State) moveInBoard(offset int) {
	s.Nav.SelectedListIndex = (s.ListsLen() + s.Nav.SelectedListIndex + offset) % s.ListsLen()
	s.Nav.SelectedCardID = s.ListCardsIds(s.Nav.SelectedListIndex)[0]
}

func (s *State) moveInList(offset int) {
	var cardIDS = s.ListCardsIds(s.Nav.SelectedListIndex)
	for i, v := range cardIDS {
		if v == s.Nav.SelectedCardID {
			s.Nav.SelectedCardID = cardIDS[(len(cardIDS)+i+offset)%len(cardIDS)]
			break
		}
	}
}

func (s *State) AppendErr(err error) {
	s.ErrorList = append(s.ErrorList, err)
}

func (s *State) SetBoardState(boardState BoardState) {
	prevState := s.BoardState
	s.Updated = time.Now()
	switch boardState {
	case BoardUninitialized:
		s.Board = trello.Board{}
		s.Lists = []trello.List{}
		s.Cards = []trello.Card{}
		s.Nav.SelectedListIndex = -1
		s.Nav.SelectedCardID = -1
		s.Nav.SelectedCardState = CardUninitialized
	case BoardLoading:
		s.Board = trello.Board{}
		s.Lists = []trello.List{}
		s.Cards = []trello.Card{}
		s.Nav.SelectedListIndex = -1
		s.Nav.SelectedCardID = -1
		s.Nav.SelectedCardState = CardLoading
	case BoardNotFound:
		s.Board = trello.Board{}
		s.Lists = []trello.List{}
		s.Cards = []trello.Card{}
		s.Nav.SelectedListIndex = -1
		s.Nav.SelectedCardID = -1
		s.Nav.SelectedCardState = CardNotFound
	case BoardLoaded:
		s.Nav.SelectedCardState = CardLoaded
	}
	s.BoardState = boardState
	log.Debug().Int("old", int(prevState)).Int("new", int(s.BoardState)).Msg("board state changed")
}
