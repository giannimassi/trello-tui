package components

import (
	"github.com/jroimartin/gocui"
	"github.com/pkg/errors"

	"github.com/giannimassi/trello-tui/v2/pkg/gui/panel"
	"github.com/giannimassi/trello-tui/v2/pkg/gui/theme"
)

type HeaderState interface {
	HeaderTitle() string
	HeaderSubtitle() string
}

type Header struct {
	*panel.Panel
}

func NewHeader(pp panel.Parent, x0, y0, w, h float64) Header {
	return Header{panel.RelativePanel("header", x0, y0, w, h).WithParent(pp)}
}

func (h *Header) Draw(g *gocui.Gui, t Theme, s HeaderState) error {
	if err := h.Panel.Layout(g); err != nil {
		return errors.Wrapf(err, "while laying ot header")
	}
	h.Panel.View.Title = " " + s.HeaderTitle() + " "
	h.Panel.Clear()
	_, _ = h.Panel.Write([]byte("\n"))
	if _, err := t.Color(theme.BoardDescriptionClass, false).Fprint(h.Panel, "  "+s.HeaderSubtitle()); err != nil {
		return err
	}

	// TODO: other info about board could go here (members, notifications maybe?)
	return nil
}
