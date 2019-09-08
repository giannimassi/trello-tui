package components

import (
	"github.com/giannimassi/trello-tui/pkg/gui/panel"
	"github.com/jroimartin/gocui"
	"github.com/pkg/errors"
)

type Header struct {
	*panel.Panel
}

func NewHeader(pp panel.Parent, x0, y0, w, h float64) Header {
	return Header{panel.RelativePanel("header", x0, y0, w, h).WithParent(pp)}
}

func (h *Header) Draw(g *gocui.Gui, ctx *Context) error {
	if err := h.Panel.Layout(g); err != nil {
		return errors.Wrapf(err, "while laying ot header")
	}
	h.Panel.View.Title = ctx.HeaderTitle()
	h.Panel.Clear()

	if ctx.HasDescription() {
		if _, err := ctx.Color(BoardDescriptionClass, false).Fprint(h.Panel, ctx.Description()); err != nil {
			return err
		}
	}

	// Other info about board could go here (members, notifications maybe?)
	return nil
}
