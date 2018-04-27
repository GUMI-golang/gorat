package textrat

import (
	"github.com/GUMI-golang/gorat"
	"github.com/GUMI-golang/gumi/gcore"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/math/fixed"
)

func raster(ctx gorat.VectorDrawer, buf *truetype.GlyphBuf, point mgl32.Vec2) {
	var start int
	for _, end := range buf.Ends {
		contour(ctx, buf.Points[start:end], point)
		start = end
	}
}
func Fint32ToFloat32(i fixed.Int26_6) float32 {
	return float32(i) / float32(0x40)
}
func contour(ctx gorat.VectorDrawer, points []truetype.Point, point mgl32.Vec2) {
	if len(points) == 0 {
		return
	}
	var first mgl32.Vec2
	var ifirst, ilast = 0, len(points)
	if points[0].Flags&0x01 != 0 {
		ifirst = 1
		first = mgl32.Vec2{
			point[0] + Fint32ToFloat32(points[0].X),
			point[1] - Fint32ToFloat32(points[0].Y),
		}
	} else {
		last := mgl32.Vec2{
			point[0] + Fint32ToFloat32(points[ilast-1].X),
			point[1] - Fint32ToFloat32(points[ilast-1].Y),
		}
		if points[ilast-1].Flags&0x01 != 0 {
			first = last
			ilast = ilast - 1
		} else {
			first = mgl32.Vec2{
				(first.X() + last.X()) / 2,
				(first.Y() + last.Y()) / 2,
			}
		}
	}
	//==================================
	// drawloop
	// start point
	ctx.MoveTo(first)
	var q0, q0on = first, true
	for i := ifirst; i < ilast; i++ {

		p := points[i]

		var q, qon = mgl32.Vec2{
			point.X() + Fint32ToFloat32(p.X),
			point.Y() - Fint32ToFloat32(p.Y),
		}, p.Flags&0x01 != 0
		if qon {
			if q0on {
				ctx.LineTo(q)
			} else {
				ctx.QuadTo(q0, q)
			}
		} else {
			if !q0on {
				ctx.QuadTo(q0, q0.Add(q).Mul(0.5))
			}
		}
		q0, q0on = q, qon
	}
	if q0on {
	} else {
		ctx.QuadTo(q0, first)
	}

}

func alignHelp(align gcore.Align, size mgl32.Vec2) (res mgl32.Vec2) {
	v, h := gcore.SplitAlign(align)
	switch v {
	case gcore.AlignTop:
		res[1] = size[1]
	case gcore.AlignVertical:
		res[1] = size[1] / 2
	case gcore.AlignBottom:

	}
	switch h {
	case gcore.AlignLeft:

	case gcore.AlignHorizontal:
		res[0] = -size[0] / 2
	case gcore.AlignRight:
		res[0] = -size[0]

	}
	return
}
