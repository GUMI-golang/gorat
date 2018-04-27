package textrat

import (
	"github.com/GUMI-golang/gorat"
	"github.com/GUMI-golang/gumi/gcore"
	"github.com/go-gl/mathgl/mgl32"
	"golang.org/x/image/font"
	"image"
)

type Font interface {
	Name() FontName
	Size() int
	SetSize(size int)
	Hint() font.Hinting
	SetHint(hint font.Hinting)
	//
	Text(ctx gorat.VectorDrawer, text string, point mgl32.Vec2, align gcore.Align)
	TextInRect(ctx gorat.VectorDrawer, text string, rect image.Rectangle, align gcore.Align)
	PathText(ctx gorat.VectorDrawer, text string, point mgl32.Vec2, align gcore.Align)
	PathTextInRect(ctx gorat.VectorDrawer, text string, rect image.Rectangle, align gcore.Align)
	MeasureText(text string) mgl32.Vec2
}
