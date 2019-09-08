package state

import (
	"time"

	"github.com/VojtechVitek/go-trello"
	"github.com/jroimartin/gocui"
	"github.com/rs/zerolog/log"
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

	ErrorList []error
	Nav       NavigationPosition
	loading   bool
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
	return s.Board.Name
}

func (s *State) Description() string {
	return s.Board.Desc
}

func (s *State) ListsLen() int {
	return len(s.Lists)
}

func (s *State) ListNameByIndex(idx int) (string, bool) {
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

func (s *State) CardNameByID(cardID int) (string, bool) {
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

func (s *State) NavPosition() NavigationPosition {
	return s.Nav
}

type NavigationPosition struct {
	SelectedBoard     string
	SelectedListIndex int
	SelectedCardID    int
	SelectedCardState CardState
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
	log.Warn().Int("offset", offset).Int("prev-list-index", s.Nav.SelectedListIndex).Int("prev-card-id", s.Nav.SelectedCardID).Msg("Move in board")
	s.Nav.SelectedListIndex = (s.ListsLen() + s.Nav.SelectedListIndex + offset) % s.ListsLen()
	s.Nav.SelectedCardID = s.ListCardsIds(s.Nav.SelectedListIndex)[0]
	cardName, _ := s.CardNameByID(s.Nav.SelectedCardID)
	listName, _ := s.ListNameByIndex(s.Nav.SelectedListIndex)
	log.Warn().Int("list-index", s.Nav.SelectedListIndex).Int("card-id", s.Nav.SelectedCardID).Str("list-name", listName).Str("card-name", cardName).Msg("Moved in board")
}

func (s *State) moveInList(offset int) {
	log.Warn().Int("offset", offset).Int("prev-list-index", s.Nav.SelectedListIndex).Int("prev-card-id", s.Nav.SelectedCardID).Msg("Move in list")
	var cardIDS = s.ListCardsIds(s.Nav.SelectedListIndex)
	for i, v := range cardIDS {
		if v == s.Nav.SelectedCardID {
			s.Nav.SelectedCardID = cardIDS[(len(cardIDS)+i+offset)%len(cardIDS)]
			break
		}
	}
	cardName, _ := s.CardNameByID(s.Nav.SelectedCardID)
	listName, _ := s.ListNameByIndex(s.Nav.SelectedListIndex)
	log.Warn().Int("list-index", s.Nav.SelectedListIndex).Int("card-id", s.Nav.SelectedCardID).Str("list-name", listName).Str("card-name", cardName).Msg("Moved in list")
}

func (s *State) AppendErr(err error) {
	s.ErrorList = append(s.ErrorList, err)
}

func (s *State) SetLoading(loading bool) {
	s.loading = loading
}
