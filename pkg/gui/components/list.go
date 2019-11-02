package components

import (
	"bytes"
	"fmt"
	"math"

	"github.com/jroimartin/gocui"
	"github.com/pkg/errors"

	"github.com/giannimassi/trello-tui/pkg/gui/panel"
	"github.com/giannimassi/trello-tui/pkg/gui/theme"
)

const (
	cardContentHeight = 6
)

type ListState interface {
	IsListSelected(idx int) bool
	ListName(idx int) string
	ListCardsIds(idx int) []int
	FirstVisibleCardIndex(listIndex, cardsPerPage int) int
	CardName(id int) string
	IsCardSelected(id int) bool
}

type List struct {
	index int
	*panel.Panel
	visible bool
}

func NewList(idx int, x0, y0, w, h float64) List {
	return List{idx, panel.RelativePanel(fmt.Sprintf("list-%02d", idx), x0, y0, w, h), false} //WithChildren(
}

func (l *List) Draw(g *gocui.Gui, t Theme, s ListState) error {
	if err := l.Panel.Layout(g); err != nil {
		return errors.Wrapf(err, "list layout failure")
	}

	l.Title = " " + s.ListName(l.index) + " "
	l.Autoscroll = true
	if s.IsListSelected(l.index) {
		l.Select(g)
	}

	l.SelBgColor = gocui.ColorMagenta
	l.SelFgColor = gocui.ColorGreen

	w, h := l.Size()
	l.Panel.Clear()
	if !l.visible {
		return nil
	}

	var (
		cardsPerPage = h / cardContentHeight
		cardIds      = s.ListCardsIds(l.index)
	)
	firstCardIndex := s.FirstVisibleCardIndex(l.index, cardsPerPage)
	_, _ = l.Panel.Write([]byte("\n\n"))
	for i, id := range cardIds {
		if i < firstCardIndex || i > firstCardIndex+cardsPerPage {
			continue
		}
		n, title := normalizeTitle(s.CardName(id), w-1, cardContentHeight-2)
		title = fmt.Sprintf("#%d id: %d - ", i, id) + title
		_, _ = t.Color(theme.CardTitleClass, s.IsCardSelected(id)).Fprint(l.Panel, " "+title+" ")
		_, _ = l.Panel.Write(bytes.Repeat([]byte("\n"), cardContentHeight-n-1))
		_, _ = l.Panel.Write(bytes.Repeat([]byte("-"), w-1))
	}
	return nil
}

func normalizeTitle(title string, w, l int) (int, string) {
	if len(title) > w*(l-1) {
		w = w - 3
		return l, title[:w*(l-1)] + "..."
	}
	return int(math.Ceil(float64(len(title)) / float64(w))), title
}
