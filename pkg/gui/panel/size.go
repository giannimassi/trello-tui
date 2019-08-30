package panel

type absoluteSize struct {
	w, h int
}

func (s *absoluteSize) Apply(x0, y0, x1, y1 *int) {
	*x1 = *x0 + s.w
	*y1 = *y0 + s.h
}

type relativeSize struct {
	absoluteSize
	parent Parent
	wR, hR float64
}

func newRelativeSize(parent Parent, widthR, heightR float64) *relativeSize {
	return &relativeSize{
		parent: parent,
		wR:     widthR,
		hR:     heightR,
	}
}

func (s *relativeSize) Apply(x0, y0, x1, y1 *int) {
	s.w, s.h = s.evalSize()
	s.absoluteSize.Apply(x0, y0, x1, y1)
}

func (s *relativeSize) evalSize() (w, h int) {
	W, H := s.parent.Size()
	return applyRatio(W, s.wR), applyRatio(H, s.hR)
}
