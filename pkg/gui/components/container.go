package components

import (
	"github.com/giannimassi/trello-tui/pkg/gui/panel"
	"github.com/jroimartin/gocui"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

const (
	headerHeight = 0.1
	interMargin  = 0.01
)

type Container struct {
	pp     panel.Parent
	header Header
	lists  []List
}

func NewContainer(pp panel.Parent) Container {
	return Container{
		pp:     pp,
		header: NewHeader(pp, 0, 0, 0.8, headerHeight-interMargin/2),
	}
}

func (c *Container) Draw(g *gocui.Gui, ctx *Context) error {
	c.header.Panel.Reflow()
	if err := c.header.Draw(g, ctx); err != nil {
		return errors.Wrapf(err, "while drawing header")
	}
	c.updateLists(g, ctx)

	for i := range c.lists {
		c.lists[i].Panel.Reflow()
		if err := c.lists[i].Draw(g, ctx); err != nil {
			return errors.Wrapf(err, "while drawing list %d", i)
		}
	}

	return nil
}

func (c *Container) updateLists(g *gocui.Gui, ctx *Context) {
	var (
		listsLen     = ctx.ListsLen()
		maxW, _      = g.Size()
		listsPerPage = listWidths.Match(maxW)
		w            = 1 / float64(listsPerPage)
		h            = 1 - headerHeight - interMargin/2
		x0           float64
		y0           = headerHeight + interMargin/2
	)

	for i := 0; i < listsLen || i < len(c.lists); i++ {
		x0 = float64(i) / float64(listsPerPage)
		switch {
		case i >= len(c.lists):
			// Add list
			c.lists = append(c.lists, NewList(i, x0, y0, w, h))
			c.lists[i].Panel = c.lists[i].Panel.WithParent(c.pp)
		case i >= listsLen:
			// Delete list
			l := c.lists[i]
			if err := g.DeleteView(l.Name()); err != nil {
				log.Error().Err(err).Str("name", l.Name()).Msg("Unexpected err while deleting view")
			}
		}
		// Resize list
		c.lists[i].Panel.Move(x0, y0)
		c.lists[i].Panel.Resize(w, h)
		c.lists[i].visible = x0 < 1
	}

	c.lists = c.lists[:listsLen]
}

type listWidthMap [][2]int

func (l *listWidthMap) Match(w int) int {
	for _, v := range listWidths {
		if w <= v[0] {
			return v[1]
		}
	}
	return 1
}

var listWidths = listWidthMap{
	{50, 1},
	{100, 2},
	{150, 3},
	{200, 4},
	{250, 5},
	{300, 6},
	{450, 7},
	{100000000, 8},
}
