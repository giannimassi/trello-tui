package panel

import "math"

type ParentFunc func() Parent

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

func (a *absolutePosition) Move(x0, y0 float64) {
	a.x0 = int(x0)
	a.y0 = int(y0)
}

type relativePosition struct {
	absolutePosition
	Parent   ParentFunc
	x0R, y0R float64
}

func newRelativePosition(pFunc ParentFunc, x0R, y0R float64) *relativePosition {
	return &relativePosition{
		Parent: pFunc,
		x0R:    x0R,
		y0R:    y0R,
	}
}

func (r *relativePosition) Apply(x0, y0, x1, y1 *int) {
	r.absolutePosition.Move(r.evalPosition())
	r.absolutePosition.Apply(x0, y0, x1, y1)
}

func (r *relativePosition) Move(x0R, y0R float64) {
	r.x0R = x0R
	r.y0R = y0R
}

func (r *relativePosition) evalPosition() (x0, y0 float64) {
	p := r.Parent()
	X0, Y0 := p.Position()
	W, H := p.Size()
	return X0 + applyRatioToInt(W, r.x0R), Y0 + applyRatioToInt(H, r.y0R)
}

func applyRatioToInt(v int, r float64) float64 {
	return math.Floor(float64(v) * r)
}
