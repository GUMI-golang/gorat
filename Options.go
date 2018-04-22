package gorat

type Options struct {
	strokeJoin  StrokeJoin
	strokeCap   StrokeCap
	strokeWidth float32
	//
	overlap Overlap
	filler  Filler
	aa      AntiAliasing
	// TODO
	//

}


func (s *Options ) DefaultOption()  {
	s.filler = NewColorFiller(0,0,0,255)
}
func (s Options) Backup() Options {
	return Options{
		strokeJoin:  s.strokeJoin,
		strokeCap:   s.strokeCap,
		strokeWidth: s.strokeWidth,
		filler:      s.filler,
		overlap:     s.overlap,
	}
}
func (s *Options) Restore(o Options) {
	s.strokeJoin = o.strokeJoin
	s.strokeCap = o.strokeCap
	s.strokeWidth = o.strokeWidth
	s.filler = o.filler
	s.overlap = o.overlap
}

// getter
func (s Options) GetStrokeJoin() StrokeJoin {
	return s.strokeJoin
}
func (s Options) GetStrokeCap() StrokeCap {
	return s.strokeCap
}
func (s Options) GetStrokeWidth() float32 {
	return s.strokeWidth
}
func (s Options) GetFiller() Filler {
	return s.filler
}
func (s Options) GetOverlap() Overlap {
	return s.overlap
}
func (s Options) GetAntiAliasing() AntiAliasing {
	return s.aa
}

// setter
func (s *Options) SetStrokeJoin(j StrokeJoin) {
	s.strokeJoin = j
}
func (s *Options) SetStrokeCap(c StrokeCap) {
	s.strokeCap = c
}
func (s *Options) SetStrokeWidth(w float32) {
	s.strokeWidth = w
}
func (s *Options) SetFiller(f Filler) {
	if f == nil{
		f = NewColorFiller(0,0,0,255)
	}
	s.filler = f
}
func (s *Options) SetOverlap(o Overlap) {
	s.overlap = o
}
func (s *Options) SetAntiAliasing(aa AntiAliasing) {
	s.aa = aa
}
