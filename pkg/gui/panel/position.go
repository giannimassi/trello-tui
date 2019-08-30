package panel

import "math"

type absolutePosition struct {
	x0, y0 int
}

func (a *absolutePosition) Apply(x0, y0, x1, y1 *int) {
	w := *x1 - *x0
	h := *y1 - *y0
	*x0 = a.x0
	*y0 = a.y0
	*x1 = a.x0 + w
	*y1 = a.y0 + h
}

type relativePosition struct {
	absolutePosition
	parent   Parent
	x0R, y0R float64
}

func newRelativePosition(parent Parent, x0R, y0R float64) *relativePosition {
	return &relativePosition{
		parent: parent,
		x0R:    x0R,
		y0R:    y0R,
	}
}

func (r *relativePosition) Apply(x0, y0, x1, y1 *int) {
	r.x0, r.y0 = r.evalPosition()
	r.absolutePosition.Apply(x0, y0, x1, y1)
}

func (r *relativePosition) evalPosition() (x0, y0 int) {
	X0, Y0 := r.parent.Position()
	W, H := r.parent.Size()
	return X0 + applyRatio(W, r.x0R), Y0 + applyRatio(H, r.y0R)
}

func applyRatio(v int, r float64) int {
	vr := float64(v) * r
	vf := math.Floor(vr)
	vi := int(vf)
	return vi
}
