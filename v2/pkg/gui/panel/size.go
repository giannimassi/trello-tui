package panel

type absoluteSize struct {
	w, h int
}

func (s *absoluteSize) Apply(x0, y0, x1, y1 *int) {
	*x1 = *x0 + s.w
	*y1 = *y0 + s.h
}

func (s *absoluteSize) Resize(w, h float64) {
	s.w = int(w)
	s.h = int(h)
}

type relativeSize struct {
	absoluteSize
	Parent ParentFunc
	wR, hR float64
}

func newRelativeSize(pFunc ParentFunc, widthR, heightR float64) *relativeSize {
	return &relativeSize{
		Parent: pFunc,
		wR:     widthR,
		hR:     heightR,
	}
}

func (s *relativeSize) Apply(x0, y0, x1, y1 *int) {
	s.absoluteSize.Resize(s.evalSize())
	s.absoluteSize.Apply(x0, y0, x1, y1)
}

func (s *relativeSize) Resize(wR, hR float64) {
	s.wR = wR
	s.hR = hR
}

func (s *relativeSize) evalSize() (w, h float64) {
	p := s.Parent()
	W, H := p.Size()
	return applyRatioToInt(W, s.wR), applyRatioToInt(H, s.hR)
}
