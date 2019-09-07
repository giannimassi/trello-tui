package components

import (
	"fmt"
	"strings"

	"github.com/giannimassi/trello-tui/pkg/gui/panel"
	"github.com/jesseduffield/gocui"
	"github.com/pkg/errors"
)

type List struct {
	idx int
	*panel.Panel
}

func NewList(idx int, x0, y0, w, h float64) *List {
	return &List{idx, panel.RelativePanel(fmt.Sprintf("list-%02d", idx), x0, y0, w, h)} //WithChildren(
}

func (l *List) Draw(g *gocui.Gui, ctx *Context) error {
	if err := l.Panel.Layout(g); err != nil {
		return errors.Wrapf(err, "list layout failure")
	}
	if l.Panel.View != nil {
		l.Title = ctx.ListTitle(l.idx)
		l.Autoscroll = true

		var content = "\n"
		w, _ := l.Size()
		line := strings.Repeat("-", w-1) + "\n"
		for _, id := range ctx.ListCardsIds(l.idx) {
			content = content + ctx.CardTitle(id) + "\n" + line
		}

		l.Panel.Clear()
		if _, err := fmt.Fprintf(l.Panel, content); err != nil {
			return err
		}
	}
	return nil
}
