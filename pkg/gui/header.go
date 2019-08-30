package gui

import (
	"fmt"

	"github.com/jesseduffield/gocui"
	"github.com/pkg/errors"

	"github.com/giannimassi/trello-cli/pkg/gui/panel"
)

type Header struct {
	panel *panel.Panel
}

func (h *Header) Layout(g *gocui.Gui, view *View) error {
	if err := h.panel.Layout(g); err != nil {
		return errors.Wrapf(err, "header layout failure")
	}
	h.panel.Clear()
	if _, err := fmt.Fprintf(h.panel, view.HeaderTitle()); err != nil {
		return err
	}
	return nil
}
