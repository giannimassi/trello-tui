package domain

import (
	"sort"
	"time"
)

type Board struct {
	Updated     time.Time
	IsEmpty     bool
	Name        string
	Description string
	Lists       []List
}

func (b *Board) CardById(id int) (Card, bool) {
	for _, l := range b.Lists {
		if c, found := l.CardsByID[id]; found {
			return c, true
		}
	}

	return Card{}, false
}

func NewBoard(name, description string, lists []List, isEmpty bool) *Board {
	return &Board{
		Updated:     time.Now(),
		IsEmpty:     isEmpty,
		Name:        name,
		Description: description,
		Lists:       lists,
	}
}

type List struct {
	Id        string
	Name      string
	CardsByID map[int]Card
	CartIds   []int
}

func NewList(id, name string, cardsById map[int]Card) List {
	cardIds := make([]int, 0, len(cardsById))
	for id := range cardsById {
		cardIds = append(cardIds, id)
	}
	sort.Slice(cardIds, func(i, j int) bool {
		return cardsById[cardIds[i]].Pos < cardsById[cardIds[j]].Pos
	})

	return List{
		Id:        id,
		Name:      name,
		CardsByID: cardsById,
		CartIds:   cardIds,
	}
}

type Card struct {
	Id          string
	Name        string
	Description string
	Pos         float64
}

func NewCard(id, name, description string, pos float64) Card {
	return Card{
		Id:          id,
		Name:        name,
		Description: description,
		Pos:         pos,
	}
}
