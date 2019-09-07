package components

import (
	"github.com/giannimassi/trello-tui/pkg/gui/panel"
	"github.com/jesseduffield/gocui"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

const (
	headerHeight = 0.05
)

type Container struct {
	pp     panel.Parent
	header *Header
	lists  []*List
}

func NewContainer(pp panel.Parent) *Container {
	return &Container{
		pp:     pp,
		header: NewHeader(pp, 0, 0, 1, headerHeight),
	}
}

func (c *Container) Draw(g *gocui.Gui, ctx *Context) error {
	if err := c.header.Draw(g, ctx); err != nil {
		return errors.Wrapf(err, "while drawing header")
	}

	if ctx.Loading() {
		return nil
	}

	if len(c.lists) != ctx.ListsLen() {
		c.updateLists(g, ctx)
	}

	for i := range c.lists {
		if err := c.lists[i].Draw(g, ctx); err != nil {
			return errors.Wrapf(err, "while drawing list %d", i)
		}
	}

	return nil
}

func (c *Container) updateLists(g *gocui.Gui, ctx *Context) {
	var (
		listsLen = ctx.ListsLen()
		x0, w    float64
	)

	log.Warn().Int("old", len(c.lists)).Int("new", listsLen).Msg("updating lists")
	for i := 0; i < listsLen || i < len(c.lists); i++ {
		x0 = float64(i) / float64(listsLen)
		w = 1 / float64(listsLen)

		switch {
		case i >= len(c.lists):
			// Add list
			log.Warn().Int("i", i).Msg("adding list")
			c.lists = append(c.lists, NewList(i, x0, headerHeight, w, 1-headerHeight))
			c.lists[i].Panel = c.lists[i].Panel.WithParent(c.pp)
		case i >= listsLen:
			// Delete list
			log.Warn().Int("i", i).Msg("deleting list")
			l := c.lists[i]
			if err := g.DeleteView(l.Name()); err != nil {
				log.Error().Err(err).Str("name", l.Name()).Msg("Unexpected err while deleting view")
			}
		default:
			// Resize list
			log.Warn().Int("i", i).Msg("resizing list")
			c.lists[i].Panel.Move(x0, headerHeight)
			c.lists[i].Panel.Resize(w, 1-headerHeight)
		}
	}

	c.lists = c.lists[:listsLen]
}
