package panel

import (
	"fmt"
	"math/rand"

	"github.com/jroimartin/gocui"
	"github.com/pkg/errors"
)

type Parent interface {
	Size() (w, h int)
	Position() (x0, y0 float64)
}

type Child interface {
	Name() string
	SetParent(p Parent)
	Layout(g *gocui.Gui) error
}

type position interface {
	Apply(x0, y0, x1, y1 *int)
	Move(x0, y0 float64)
}

type size interface {
	Apply(x0, y0, x1, y1 *int)
	Resize(w, h float64)
}

type Panel struct {
	*gocui.View
	name           string
	x0, y0, x1, y1 int

	size     size
	position position

	parent   Parent
	children []Child
}

func AbsolutePanel(name string, x0, y0, w, h int) *Panel {
	return &Panel{
		name:     name,
		position: &absolutePosition{x0, y0},
		size:     &absoluteSize{w, h},
	}
}

func RelativePanel(name string, x0, y0, w, h float64) *Panel {
	if name == "" {
		name = fmt.Sprintf("v%00002d", rand.Intn(99999))
	}

	p := &Panel{
		name: name,
	}
	p.position = newRelativePosition(p.Parent, x0, y0)
	p.size = newRelativeSize(p.Parent, w, h)
	return p
}

func (p *Panel) Reflow() {
	p.position.Apply(&p.x0, &p.y0, &p.x1, &p.y1)
	p.size.Apply(&p.x0, &p.y0, &p.x1, &p.y1)
}

func (p *Panel) Layout(g *gocui.Gui) error {
	v, err := g.SetView(p.name, p.x0, p.y0, p.x1, p.y1)
	if err != nil && err.Error() != gocui.ErrUnknownView.Error() {
		return err
	}
	p.View = v
	p.View.Wrap = true
	p.View.Autoscroll = true
	for _, c := range p.children {
		err := c.Layout(g)
		if err != nil {
			return errors.Wrapf(err, "while drawing child %s", c.Name())
		}
	}

	return nil
}

func (p *Panel) Parent() Parent {
	return p.parent
}

func (p *Panel) Position() (float64, float64) {
	return float64(p.x0), float64(p.y0)
}

func (p *Panel) Size() (int, int) {
	return p.x1 - p.x0, p.y1 - p.y0
}
func (p *Panel) Move(x0, y0 float64) {
	p.position.Move(x0, y0)
}

func (p *Panel) Resize(w, h float64) {
	p.size.Resize(w, h)
}

func (p *Panel) SetParent(pp Parent) {
	p.parent = pp
}

func (p *Panel) WithChildren(panels ...Child) *Panel {
	p.children = panels
	for _, c := range p.children {
		c.SetParent(p)
	}
	return p
}

func (p *Panel) WithParent(pp Parent) *Panel {
	p.SetParent(pp)
	return p
}
