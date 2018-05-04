package gorat

import (
	"github.com/go-gl/mathgl/mgl32"
	"image"
	"image/color"
)

type Rasterizer interface {
	Setup(w, h int)
	SubRasterizer
}

type SubRasterizer interface {
	Root() Rasterizer
	Bound() image.Rectangle
	SubRasterizer(r image.Rectangle) SubRasterizer
	VectorDrawer
}

type VectorDrawer interface {
	// Infomation
	Size() (w, h float32)
	PreviousPoint() mgl32.Vec2
	Point() mgl32.Vec2
	// Path
	Reset()
	MoveTo(to mgl32.Vec2)
	LineTo(to mgl32.Vec2)
	QuadTo(pivot, to mgl32.Vec2)
	CubeTo(pivot1, pivot2, to mgl32.Vec2)
	CloseTo()
	// Draw
	Clear()
	Stroke()
	Fill()
	FillStroke()
	//
	// Options
	VectorOptions
}
type VectorOptions interface {
	DefaultOption()
	Restore(opt Options)
	Clone() Options
	// getter
	GetFiller() Filler
	GetStrokeWidth() float32
	GetStrokeJoin() StrokeJoin
	GetStrokeCap() StrokeCap
	GetStrokeColor() color.Color
	// setter
	SetFiller(f Filler)
	SetStrokeWidth(w float32)
	SetStrokeJoin(j StrokeJoin)
	SetStrokeCap(c StrokeCap)
	SetStrokeColor(c color.Color)
}
