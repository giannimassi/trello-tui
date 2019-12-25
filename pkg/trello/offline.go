package trello

type boardOffline struct {
	boardOnline
	err error
}

func (b *boardOffline) HeaderTitle() string {
	return b.boardName + " - offline"
}

func (b *boardOffline) HeaderSubtitle() string {
	if b.err != nil {
		return "Could not load board.\nDetails:\n" + b.err.Error()
	}
	return "Could not load board (unknown error)."
}
