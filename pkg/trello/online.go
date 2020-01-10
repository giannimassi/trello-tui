package trello

import (
	"strings"

	"github.com/VojtechVitek/go-trello"
	"github.com/giannimassi/trello-tui/pkg/store"
	"github.com/rs/zerolog/log"
)

type boardOnline struct {
	boardLoading
	Board              *trello.Board
	Lists              []trello.List
	CardIDsByListIndex [][]int
	CardsByID          map[int]trello.Card
	ActionsByCardID    map[int][]trello.Action
}

func (b *boardOnline) online(newBoard *trello.Board, lists []trello.List, cards []trello.Card) store.State {
	b.Board = newBoard
	b.Lists = lists
	listIndexes := make(map[string]int, len(lists))
	for i, ls := range lists {
		listIndexes[ls.Id] = i
	}

	b.CardIDsByListIndex = make([][]int, len(lists))
	b.CardsByID = make(map[int]trello.Card, len(cards))
	for _, c := range cards {
		listIndex := listIndexes[c.IdList]
		if b.CardIDsByListIndex[listIndex] == nil {
			b.CardIDsByListIndex[listIndex] = make([]int, 0)
		}
		b.CardIDsByListIndex[listIndex] = append(b.CardIDsByListIndex[listIndex], c.IdShort)
		b.CardsByID[c.IdShort] = c
	}

	return b
}

func (b *boardOnline) offline(err error) store.State {
	offline := &boardOffline{
		boardOnline: *b,
		err:         err,
	}
	return offline
}

func (b *boardOnline) setCardActions(id int, actions []trello.Action) store.State {
	if b.ActionsByCardID == nil {
		b.ActionsByCardID = make(map[int][]trello.Action)
	}
	b.ActionsByCardID[id] = actions
	return b
}

func (b *boardOnline) HeaderTitle() string {
	return b.boardName + " - online"
}

func (b *boardOnline) HeaderSubtitle() string {
	return b.Board.Desc
}

func (b *boardOnline) ListsLen() int {
	return len(b.Lists)
}

func (b *boardOnline) ListName(idx int) string {
	if idx >= len(b.Lists) {
		return ""
	}
	return b.Lists[idx].Name
}

func (b *boardOnline) ListCardsIds(idx int) []int {
	if idx >= len(b.Lists) {
		return []int{}
	}

	return b.CardIDsByListIndex[idx]
}

func (b *boardOnline) CardName(id int) string {
	c, found := b.CardsByID[id]
	if !found {
		log.Error().Int("id", id).Msg("Card not found")
		return ""
	}
	return c.Name
}

func (b *boardOnline) CardLabelsStr(id int) string {
	c, found := b.CardsByID[id]
	if !found {
		log.Error().Int("id", id).Msg("Card not found")
		return ""
	}
	if len(c.Labels) == 0 {
		return ""
	}
	var strs []string
	for _, lbl := range c.Labels {
		textColor := "black"
		if lbl.Color == "black" {
			textColor = "white"
		}
		strs = append(strs, "["+textColor+":"+lbl.Color+"] "+lbl.Name+" [-:-:-]")
	}
	return strings.Join(strs, " ")
}

func (b *boardOnline) Description(id int) string {
	c, found := b.CardsByID[id]
	if !found {
		log.Error().Int("id", id).Msg("Card not found")
		return ""
	}
	return c.Desc
}

func (b *boardOnline) CardComments(id int) []string {
	c, found := b.CardsByID[id]
	if !found {
		log.Error().Int("id", id).Msg("Card not found")
		return []string{}
	}
	b.cardCommentsRequested(&c)
	comments := make([]string, 0)
	for _, c := range b.ActionsByCardID[id] {
		if c.Type != "commentCard" {
			continue
		}
		comments = append(comments, c.Data.Text+"\n[yellow]"+c.MemberCreator.Username+"[-:-:-]")
	}

	return comments
}
