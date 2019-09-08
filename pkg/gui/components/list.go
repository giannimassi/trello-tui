package components

import (
	"fmt"
	"strings"

	"github.com/giannimassi/trello-tui/pkg/gui/panel"
	"github.com/jroimartin/gocui"
	"github.com/pkg/errors"
)

type List struct {
	idx int
	*panel.Panel
}

func NewList(idx int, x0, y0, w, h float64) List {
	return List{idx, panel.RelativePanel(fmt.Sprintf("list-%02d", idx), x0, y0, w, h)} //WithChildren(
}

func (l *List) Draw(g *gocui.Gui, ctx *Context) error {
	if err := l.Panel.Layout(g); err != nil {
		return errors.Wrapf(err, "list layout failure")
	}
	var nav = ctx.NavPosition()

	l.Title = ctx.ListTitle(l.idx)
	l.Autoscroll = true
	l.SelBgColor = gocui.ColorGreen

	if nav.IsListSelected(l.idx) {
		l.View, _ = g.SetCurrentView(l.Name())
	}
	w, _ := l.Size()
	l.Panel.Clear()
	_, _ = l.Panel.Write([]byte("\n\n" + strings.Repeat("-", w-1)))
	for _, id := range ctx.ListCardsIds(l.idx) {
		_, _ = ctx.Color(CardTitleClass, nav.IsCardSelected(id)).Fprint(l.Panel, ctx.CardTitle(id))
		_, _ = l.Panel.Write([]byte("\n\n" + strings.Repeat("-", w-1)))
	}
	return nil
}
