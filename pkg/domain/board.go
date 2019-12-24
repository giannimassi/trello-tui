package domain

import (
	"sort"
	"time"
)

// Board describes a trello board
type Board struct {
	Updated     time.Time
	IsEmpty     bool
	Name        string
	Description string
	Lists       []List
}

// CardByID returns a card with the corresponding id if available
func (b *Board) CardByID(id int) (Card, bool) {
	for _, l := range b.Lists {
		if c, found := l.CardsByID[id]; found {
			return c, true
		}
	}

	return Card{}, false
}

// NewBoard returns a new instance of Board
func NewBoard(name, description string, lists []List, isEmpty bool) *Board {
	return &Board{
		Updated:     time.Now(),
		IsEmpty:     isEmpty,
		Name:        name,
		Description: description,
		Lists:       lists,
	}
}

// List describes a trello list
type List struct {
	ID        string
	Name      string
	CardsByID map[int]Card
	CartIds   []int
}

// NewList returns a new instance of List
func NewList(id, name string, cardsByID map[int]Card) List {
	cardIds := make([]int, 0, len(cardsByID))
	for id := range cardsByID {
		cardIds = append(cardIds, id)
	}
	sort.Slice(cardIds, func(i, j int) bool {
		return cardsByID[cardIds[i]].Pos < cardsByID[cardIds[j]].Pos
	})

	return List{
		ID:        id,
		Name:      name,
		CardsByID: cardsByID,
		CartIds:   cardIds,
	}
}

// Card describes a trello card which can be part of a trello list
type Card struct {
	ID          string
	Name        string
	Description string
	Pos         float64
	Labels      []CardLabel
}

// CardLabel describes a trello label which can be associated with a trello card
type CardLabel struct {
	Name  string
	Color string
}

// NewCard returns a news instance of the Card type
func NewCard(id, name, description string, pos float64, labels []CardLabel) Card {
	return Card{
		ID:          id,
		Name:        name,
		Description: description,
		Pos:         pos,
		Labels:      labels,
	}
}
