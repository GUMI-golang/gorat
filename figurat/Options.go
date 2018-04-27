package figurat

import "github.com/GUMI-golang/gorat"

type Options struct {
	strokeJoin  StrokeJoin
	strokeCap   StrokeCap
	strokeWidth float32
	filler      gorat.Filler
}
type StrokeJoin uint8

const (
	StrokeJoinBevel StrokeJoin = iota
	StrokeJoinRound StrokeJoin = iota
	StrokeJoinMiter StrokeJoin = iota
)

type StrokeCap uint8

const (
	StrokeCapButt   StrokeCap = iota
	StrokeCapRound  StrokeCap = iota
	StrokeCapSqaure StrokeCap = iota
)

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
func (s Options) GetFiller() gorat.Filler {
	return s.filler
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
func (s *Options) SetFiller(f gorat.Filler) {
	if f == nil {
		f = gorat.NewColorFiller(0, 0, 0, 255)
	}
	s.filler = f
}
