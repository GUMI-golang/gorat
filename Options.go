package gorat

import "image/color"

type Options struct {
	join   StrokeJoin
	cap    StrokeCap
	width  float32
	color  ColorFiller
	filler Filler
}

func (s *Options) DefaultOption() {
	s.join = StrokeJoinBevel
	s.cap = StrokeCapButt
	s.width = 1
	s.color = NewColorFiller(0, 0, 0, 255)
	s.filler = NewColorFiller(0, 0, 0, 255)
}
func (s *Options) Restore(opt Options) {
	s.filler = opt.filler
}
func (s *Options) Clone() Options {
	return Options{
		filler: s.filler,
	}
}

// getter
func (s *Options) GetFiller() Filler {
	return s.filler
}
func (s *Options) GetStrokeWidth() float32 {
	return s.width
}
func (s *Options) GetStrokeJoin() StrokeJoin {
	return s.join
}
func (s *Options) GetStrokeCap() StrokeCap {
	return s.cap
}
func (s *Options) GetStrokeColor() color.Color{
	return color.RGBA(s.color)
}

// setter
func (s *Options) SetFiller(f Filler) {
	if f == nil {
		f = NewColorFiller(0, 0, 0, 255)
	}
	s.filler = f
}
func (s *Options) SetStrokeWidth(w float32) {
	s.width = w
}
func (s *Options) SetStrokeJoin(j StrokeJoin) {
	s.join = j
}
func (s *Options) SetStrokeCap(c StrokeCap) {
	s.cap = c
}
func (s *Options) SetStrokeColor(c color.Color) {
	s.color = ColorFillerModel.Convert(c).(ColorFiller)
}