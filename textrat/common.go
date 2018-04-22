package textrat

import (
	"github.com/GUMI-golang/gorat"
	"golang.org/x/image/math/fixed"
	"github.com/GUMI-golang/gumi/gcore"
	"github.com/golang/freetype/truetype"
)

func raster(ctx gorat.VectorDrawer, buf *truetype.GlyphBuf, point fixed.Point52_12) {
	var start int
	for _, end := range buf.Ends {
		contour(ctx, buf.Points[start:end], point)
		start = end
	}
}
func contour(ctx gorat.VectorDrawer, points []truetype.Point, point fixed.Point52_12) {
	if len(points) == 0 {
		return
	}
	var first fixed.Point52_12
	var ifirst, ilast = 0, len(points)
	if points[0].Flags&0x01 != 0 {
		ifirst = 1
		first = fixed.Point52_12{
			X: point.X + gorat.Fixed32ToFixed64(points[0].X),
			Y: point.Y - gorat.Fixed32ToFixed64(points[0].Y),
		}
	} else {
		last := fixed.Point52_12{
			X: point.X + gorat.Fixed32ToFixed64(points[ilast-1].X),
			Y: point.Y - gorat.Fixed32ToFixed64(points[ilast-1].Y),
		}
		if points[ilast-1].Flags&0x01 != 0 {
			first = last
			ilast = ilast - 1
		} else {
			first = fixed.Point52_12{
				X: (first.X + last.X) / 2,
				Y: (first.Y + last.Y) / 2,
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

		var q, qon = fixed.Point52_12{
			X: point.X + gorat.Fixed32ToFixed64(p.X),
			Y: point.Y - gorat.Fixed32ToFixed64(p.Y),
		}, p.Flags&0x01 != 0
		if qon {
			if q0on {
				ctx.LineTo(q)
			} else {
				ctx.QuadTo(q0, q)
			}
		} else {
			if !q0on {
				ctx.QuadTo(q0, q0.Add(q).Div(fixed.Int52_12(2 << 12)))
			}
		}
		q0, q0on = q, qon
	}
	if q0on {
	} else {
		ctx.QuadTo(q0, first)
	}

}

func alignHelp(align gcore.Align, size fixed.Point52_12) (res fixed.Point52_12) {
	v, h := gcore.SplitAlign(align)
	switch v {
	case gcore.AlignTop:
		res.Y = size.Y
	case gcore.AlignVertical:
		res.Y = size.Y / 2
	case gcore.AlignBottom:

	}
	switch h {
	case gcore.AlignLeft:

	case gcore.AlignHorizontal:
		res.X = -size.X / 2
	case gcore.AlignRight:
		res.X = -size.X

	}
	return
}
