package trello

import (
	"fmt"
	"strings"

	"github.com/giannimassi/trello-tui/pkg/domain"
	"github.com/rs/zerolog/log"
)

type boardOnline struct {
	boardLoading
	Board *domain.Board
}

func (b *boardOnline) online(newBoard *domain.Board) board {
	b.Board = newBoard
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

func (b *boardOnline) ListsLen() int {
	return len(b.Board.Lists)
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
	c, found := b.Board.CardByID(id)
	if !found {
		log.Error().Int("id", id).Msg("Card not found")
		return ""
	}
	return fmt.Sprintf("%s", c.Name)
}

func (b *boardOnline) CardLabelsStr(id int) string {
	c, found := b.Board.CardByID(id)
	if !found {
		log.Error().Int("id", id).Msg("Card not found")
		return ""
	}
	var strs []string
	for _, lbl := range c.Labels {
		strs = append(strs, "[black:"+lbl.Color+"] "+lbl.Name+" [-:-:-]")
	}
	return strings.Join(strs, " ")
}

func (b *boardOnline) Description(id int) string {
	c, found := b.Board.CardByID(id)
	if !found {
		log.Error().Int("id", id).Msg("Card not found")
		return ""
	}
	return c.Description
}

func (b *boardOnline) CardComments(id int) []string {
	_, found := b.Board.CardByID(id)
	if !found {
		log.Error().Int("id", id).Msg("Card not found")
		return []string{}
	}

	return []string{
		"sdo fi jsdof nosd fosd",
		"dfoisj dofg sndofismo difs",
	}

	// return c.Comments
}

func cardIndexInListFromID(cardIds []int, id int) (int, bool) {
	for i, cardID := range cardIds {
		if cardID == id {
			return i, true
		}
	}
	return 0, false
}
