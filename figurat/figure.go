package figurat

import (
	"github.com/GUMI-golang/gumi/gorat"
	"golang.org/x/image/math/fixed"
	"math"
)

type Figure struct {
	raster gorat.Rasterizer
	transform_rotation Angle
}

func NewFigure(rasterizer gorat.VectorDrawer) Figure {
	return Figure{
		raster: rasterizer,
	}
}

//
func (s Figure) PathRect(rect fixed.Rectangle52_12) {
	s.raster.MoveTo(rect.Min)
	s.raster.LineTo(fixed.Point52_12{
		X:rect.Min.X,
		Y:rect.Max.Y,
	})
	s.raster.LineTo(rect.Max)
	s.raster.LineTo(fixed.Point52_12{
		X:rect.Max.X,
		Y:rect.Min.Y,
	})
}

const (
	circleprecision = 24
)
func (s Figure) PathEllipticalArc(center fixed.Point52_12, radius fixed.Point52_12, from, to Angle) {
	const n = 16
	var ffrom, fto = fixed.Int52_12(from), fixed.Int52_12(to)
	for i := 0; i < n; i++ {
		p1 := fixed.Int52_12((i + 0) / n)
		p2 := fixed.Int52_12((i + 1) / n)

		a1 := ffrom + (fto - ffrom) * p1
		a2 := ffrom + (fto - ffrom) * p2
		x0 := center.X + gorat.Mul(radius.X, gorat.Cos(a1))
		y0 := center.Y + gorat.Mul(radius.Y, gorat.Sin(a1))
		x1 := center.X + gorat.Mul(radius.X, gorat.Cos(a1+(a2-a1)/2))
		y1 := center.Y + gorat.Mul(radius.Y, gorat.Sin(a1+(a2-a1)/2))
		x2 := center.X + gorat.Mul(radius.X, gorat.Cos(a2))
		y2 := center.Y + gorat.Mul(radius.Y, gorat.Sin(a2))
		cx := 2*x1 - x0/2 - x2/2
		cy := 2*y1 - y0/2 - y2/2
		if i == 0 {
			if s.raster.IsBegin(){
				s.raster.LineTo(fixed.Point52_12{
					X:x0,
					Y:y0,
				})
			}else {
				s.raster.MoveTo(fixed.Point52_12{
					X:x0,
					Y:y0,
				})
			}
		}
		s.raster.QuadTo(fixed.Point52_12{
			X:cx,
			Y:cy,
		}, fixed.Point52_12{
			X:x2,
			Y:y2,
		})
	}

	//s.raster.MoveTo(rect.Min)
	//s.raster.LineTo(fixed.Point52_12{
	//	X:rect.Min.X,
	//	Y:rect.Max.Y,
	//})
	//s.raster.LineTo(rect.Max)
	//s.raster.LineTo(fixed.Point52_12{
	//	X:rect.Max.X,
	//	Y:rect.Min.Y,
	//})
}
func (s Figure) PathCircularArc(center fixed.Point52_12, radius fixed.Int52_12, from, to Angle) {
	s.PathEllipticalArc(center, fixed.Point52_12{radius, radius}, from, to)
}
func (s Figure) PathEllipse(center fixed.Point52_12, radius fixed.Point52_12) {
	s.PathEllipticalArc(center, radius, AngleZero, 2 * AnglePi)
}
func (s Figure) PathCircle(center fixed.Point52_12, radius fixed.Int52_12) {
	s.PathCircularArc(center, radius, AngleZero, 2 * AnglePi)
}
// Regular
func (s Figure) PathRegularTriangle(center fixed.Point52_12, radius fixed.Int52_12) {

}
func (s Figure) PathRegularRect(center fixed.Point52_12, radius fixed.Int52_12) {

}
func (s Figure) PathRegularPentagon(center fixed.Point52_12, radius fixed.Int52_12) {

}
func (s Figure) PathRegularHexagon(center fixed.Point52_12, radius fixed.Int52_12) {

}
func (s Figure) PathRegularPolygon(edgecount int, center fixed.Point52_12, radius fixed.Int52_12) {

}
// Inbound
func (s Figure) PathInboundTriangle(rect fixed.Rectangle52_12, angle Angle) {

}
func (s Figure) PathInboundPentagon(rect fixed.Rectangle52_12, angle Angle) {

}
func (s Figure) PathInboundHexagon(rect fixed.Rectangle52_12, angle Angle) {

}
func (s Figure) PathInboundPolygon(edgecount int, rect fixed.Rectangle52_12, angle Angle) {

}
