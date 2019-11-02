package components

import (
	"bytes"
	"fmt"
	"math"

	"github.com/jroimartin/gocui"
	"github.com/pkg/errors"

	"github.com/giannimassi/trello-tui/pkg/gui/panel"
)

const (
	cardContentHeight = 6
)

type List struct {
	idx int
	*panel.Panel
	visible bool
}

func NewList(idx int, x0, y0, w, h float64) List {
	return List{idx, panel.RelativePanel(fmt.Sprintf("list-%02d", idx), x0, y0, w, h), false} //WithChildren(
}

func (l *List) Draw(g *gocui.Gui, ctx *Context) error {
	if err := l.Panel.Layout(g); err != nil {
		return errors.Wrapf(err, "list layout failure")
	}

	l.Title = ctx.ListName(l.idx)
	l.Autoscroll = true
	l.SelBgColor = gocui.ColorGreen

	if ctx.IsListSelected(l.idx) {
		l.View, _ = g.SetCurrentView(l.Name())
	}
	w, h := l.Size()
	l.Panel.Clear()
	if !l.visible {
		return nil
	}

	var (
		cardsPerPage = h / cardContentHeight
		cardIds      = ctx.ListCardsIds(l.idx)
	)
	ctx.UpdateFirstCardIndex(cardsPerPage+1, cardIds)
	firstCardIdx := ctx.FirstVisibleCardIndex(l.idx)

	_, _ = l.Panel.Write([]byte("\n\n"))
	for i, id := range cardIds {
		if i < firstCardIdx || i > firstCardIdx+cardsPerPage {
			continue
		}

		n, title := prepTitle(ctx.CardName(id), w-1, cardContentHeight-2)
		_, _ = ctx.Color(CardTitleClass, ctx.IsCardSelected(id)).Fprint(l.Panel, title)
		_, _ = l.Panel.Write(bytes.Repeat([]byte("\n"), cardContentHeight-n-1))
		_, _ = l.Panel.Write(bytes.Repeat([]byte("-"), w-1))
	}
	return nil
}

func prepTitle(title string, w, l int) (int, string) {
	if len(title) > w*(l-1) {
		return l, title[:w*(l-1)] + "..."
	}
	return int(math.Ceil(float64(len(title)) / float64(w))), title
}
