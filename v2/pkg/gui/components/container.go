package components

import (
	"github.com/fatih/color"
	"github.com/jroimartin/gocui"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"github.com/giannimassi/trello-tui/v2/pkg/gui/panel"
)

const (
	headerHeight = 0.1
)

type Theme interface {
	Color(class string, isSelected bool) *color.Color
	ListsPerPage(maxW int) int
}

type ViewState interface {
	HeaderState
	ListState
	SelectedCardState

	IsCardPopupOpen() bool
	FirstVisibleListIndex(listsPerPage int) int
	ListsLen() int
}

type View struct {
	parent panel.Parent
	theme  Theme

	header Header
	lists  []List
	card   CardPopup
}

func NewView(pp panel.Parent, t Theme) *View {
	return &View{
		parent: pp,
		theme:  t,
		header: NewHeader(pp, 0, 0, 0.8, headerHeight),
		card:   NewCardPopup(pp, 0.1, 0.15, 0.8, 0.75),
	}
}

func (v *View) Draw(g *gocui.Gui, s ViewState) error {
	if err := v.header.Draw(g, v.theme, s); err != nil {
		return errors.Wrapf(err, "while drawing header")
	}

	v.updateLists(g, s)
	for i := range v.lists {
		if err := v.lists[i].Draw(g, v.theme, s); err != nil {
			return errors.Wrapf(err, "while drawing list %d", i)
		}
	}

	if err := v.card.SetVisible(g, s.IsCardPopupOpen()); err != nil {
		return errors.Wrapf(err, "while setting card popup visibility")
	}
	if err := v.card.Draw(g, v.theme, s); err != nil {
		return errors.Wrapf(err, "while drawing card popup")
	}

	return nil
}

func (v *View) updateLists(g *gocui.Gui, s ViewState) {
	var (
		listsLen     = s.ListsLen()
		maxW, _      = g.Size()
		listsPerPage = v.theme.ListsPerPage(maxW)
		w            = 1 / float64(listsPerPage)
		h            = 1 - headerHeight
		x0           float64
		y0           = headerHeight
	)
	for i := 0; i < listsLen || i < len(v.lists); i++ {
		x0 = float64(i-s.FirstVisibleListIndex(listsPerPage)) / float64(listsPerPage)
		switch {
		case i >= len(v.lists):
			// Add list
			v.lists = append(v.lists, NewList(i, x0, y0, w, h))
			v.lists[i].Panel = v.lists[i].Panel.WithParent(v.parent)
		case i >= listsLen:
			// Delete list
			l := v.lists[i]
			if err := g.DeleteView(l.Name()); err != nil {
				log.Error().Err(err).Str("name", l.Name()).Msg("Unexpected err while deleting view")
			}
		}
		// Resize list
		v.lists[i].Panel.Move(x0, y0)
		v.lists[i].Panel.Resize(w, h)

		if x0 >= 0 && x0 < 1 {
			v.lists[i].visible = true
		} else {
			v.lists[i].visible = false
		}
	}

	v.lists = v.lists[:listsLen]
}
