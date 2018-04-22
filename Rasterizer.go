package gorat

import (
	"golang.org/x/image/math/fixed"
	"image/color"
	"image/draw"
)

//type RasterizerNotifier interface {
//	Rasterizer
//	SwitchDo()
//	SwitchPost()
//}
type Rasterizer interface {
	Resize(w, h int)
	SubRasterizer
}
type SubRasterizer interface {
	//StartRaw() (data []byte, stride int)
	//EndRaw()
	Size() fixed.Rectangle52_12
	Draw(img draw.Image, op draw.Op)
	SubRasterizer(subbound fixed.Rectangle52_12) SubRasterizer
	NewPath() VectorPath
	VectorDrawer
}
type VectorPath interface {
	VectorDrawer
	Commit()
}
type VectorDrawer interface {
	IsBegin() bool
	MoveTo(to fixed.Point52_12)
	LineTo(to fixed.Point52_12)
	QuadTo(pivot fixed.Point52_12, to fixed.Point52_12)
	CubeTo(pivot1 fixed.Point52_12, pivot2 fixed.Point52_12, to fixed.Point52_12)
	Close()
	Stroke()
	Fill()
	Clear(bound fixed.Rectangle52_12)
	// TODO
	// Contour(cont Contour)
	//
	StrokeDash(dashes ...fixed.Int52_12)
	StrokeJoin(join StrokeJoin)
	StrokeCap(cap StrokeCap)
	StrokeWidth(width fixed.Int52_12)
	Color(color color.Color)
	FillStyle(fstyle Filler)
	Options() Options
	Revert(opts Options)

}

var ClearAll = fixed.Rectangle52_12{
	Min: fixed.Point52_12{
		X: 0,
		Y: 0,
	},
	Max: fixed.Point52_12{
		X: 1 << 63 - 1,
		Y: 1 << 63 - 1,
	},
	//maxI
}

