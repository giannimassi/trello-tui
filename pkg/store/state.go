package store

// State describes the interface required for the gui
type State interface {
	ViewState
}

// ViewState describes the interface required for the view component
type ViewState interface {
	HeaderState
	ListsState
}

// HeaderState describes the interface required for the header component
type HeaderState interface {
	HeaderTitle() string
	HeaderSubtitle() string
}

// ListsState describes the interface required for the ListContainer component
type ListsState interface {
	ListsLen() int
	SingleListState
}

// SingleListState describes the interface required for the list component
type SingleListState interface {
	ListName(idx int) string
	ListCardsIds(idx int) []int
	CardState
}

// CardState describes the interface required for the selected card component
type CardState interface {
	CardName(id int) string
	CardLabelsStr(id int) string
	Description(id int) string
}
