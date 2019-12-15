package gui

// State describes the interface required for the gui
type State interface {
	UserActionHandler
	ViewState
}

// UserActionHandler describes the interface required for handling user actions
type UserActionHandler interface {
	LeftPressed()
	RightPressed()
	UpPressed()
	DownPressed()
	EnterPressed()
	BackPressed()
}

// ViewState describes the interface required for the view component
type ViewState interface {
	HeaderState
	ListState
	SelectedCardState

	IsCardPopupOpen() bool
	FirstVisibleListIndex(listsPerPage int) int
	ListsLen() int
}

// HeaderState describes the interface required for the header component
type HeaderState interface {
	HeaderTitle() string
	HeaderSubtitle() string
}

// ListState describes the interface required for the list component
type ListState interface {
	IsListSelected(idx int) bool
	ListName(idx int) string
	ListCardsIds(idx int) []int
	FirstVisibleCardIndex(listIndex, cardsPerPage int) int
	CardName(id int) string
	IsCardSelected(id int) bool
}

// SelectedCardState describes the interface required for the selected card component
type SelectedCardState interface {
	SelectedCardName() string
	SelectedCardDescription() string
}
