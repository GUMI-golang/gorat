package textrat

import (
	"golang.org/x/image/math/fixed"
	"golang.org/x/image/font"
	"github.com/GUMI-golang/gorat"
	"github.com/GUMI-golang/gumi/gcore"
)

type Font interface {
	Name() FontName
	Size() int
	SetSize(size int)
	Hint() font.Hinting
	SetHint(hint font.Hinting)
	//
	Text(ctx gorat.VectorDrawer, text string, point fixed.Point52_12, align gcore.Align)
	TextInRect(ctx gorat.VectorDrawer, text string, point fixed.Rectangle52_12, align gcore.Align)
	PathText(ctx gorat.VectorDrawer, text string, point fixed.Point52_12, align gcore.Align)
	PathTextInRect(ctx gorat.VectorDrawer, text string, point fixed.Rectangle52_12, align gcore.Align)
	MeasureText(text string) (fixed.Point52_12)
}
