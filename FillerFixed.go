package gorat

import (
	"image"
	"image/color"
)

type FillerFixed struct {
	img *image.RGBA
}

func (s *FillerFixed) ColorAt(x, y int) color.Color {
	if inRect(s.img.Rect, image.Point{X:x, Y:y,}){
		return s.img.At(x, y)
	}
	return color.Transparent
}
func inRect(r image.Rectangle, pt image.Point) bool {
	return ((r.Min.X <= pt.X) && (pt.X < r.Max.X)) && ((r.Min.Y <= pt.Y) && (pt.Y < r.Max.Y))
}

