package panel

import (
	"github.com/jesseduffield/gocui"
)

type Parent interface {
	Size() (w, h int)
	Position() (x0, y0 int)
}

type viewApplier interface {
	Apply(x0, y0, x1, y1 *int)
}

type Panel struct {
	*gocui.View
	name           string
	x0, y0, x1, y1 int
	position, size viewApplier
	overlaps       byte
}

func NewAbsolutePanel(name string, x0, y0, w, h int) *Panel {
	return &Panel{
		name:     name,
		position: &absolutePosition{x0, y0},
		size:     &absoluteSize{w, h},
	}
}

func NewRelativePanel(name string, pp Parent, x0, y0, w, h float64) *Panel {
	return &Panel{
		name:     name,
		position: newRelativePosition(pp, x0, y0),
		size:     newRelativeSize(pp, w, h),
	}
}

func (p *Panel) Layout(g *gocui.Gui) error {
	p.position.Apply(&p.x0, &p.y0, &p.x1, &p.y1)
	p.size.Apply(&p.x0, &p.y0, &p.x1, &p.y1)
	if v, err := g.SetView(p.name, p.x0, p.y0, p.x1, p.y1, p.overlaps); err != nil {
		if err.Error() != gocui.ErrUnknownView.Error() {
			return err
		}
		p.View = v
	}
	return nil
}
