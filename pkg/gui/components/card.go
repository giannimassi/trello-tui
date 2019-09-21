package components

import (
	"github.com/giannimassi/trello-tui/pkg/gui/panel"
	"github.com/jroimartin/gocui"
	"github.com/pkg/errors"
)

type CardPopup struct {
	*panel.Panel
	visible bool
}

func NewCardPopup(pp panel.Parent, x0, y0, w, h float64) CardPopup {
	return CardPopup{panel.RelativePanel("card-popup", x0, y0, w, h).WithParent(pp), false}
}

func (c *CardPopup) Draw(g *gocui.Gui, ctx *Context) error {
	if !c.visible {
		return nil
	}

	if err := c.Panel.Layout(g); err != nil {
		return errors.Wrapf(err, "card popup layout failure")
	}

	c.Title = ctx.CardPopupTitle()
	c.Autoscroll = true
	c.SelBgColor = gocui.ColorGreen

	return nil
}

// SetVisible set the value of visible field and deletes the view if necessary
func (c *CardPopup) SetVisible(g *gocui.Gui, visible bool) error {
	if visible == c.visible {
		return nil
	}
	if !visible {
		if err := g.DeleteView(c.Name()); err != nil {
			return errors.Wrapf(err, "while deleting card popup view")
		}
		c.View = nil
		c.visible = visible
	}
	c.visible = visible
	return nil
}
