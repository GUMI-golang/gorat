package SoftwareRasterizer

import (
	"image"
	"github.com/GUMI-golang/gorat"
)

func Quad(to *image.RGBA) {
	rat := gorat.NewSoftwareRasterizerRGBA(to)

	rat.MoveTo([2]float32{50,50}, nil)
	rat.QuadTo(gorat.Vec2(670, 50), gorat.Vec2(670, 670), nil)
	rat.LineTo(gorat.Vec2(50, 670), nil)
	//rat.LineTo([2]float32{90,170}, nil)
	// rat.Close()
	// .Close() call auto
	rat.FillColor()
}
func QuadAA(to *image.RGBA) {
	rat := gorat.NewSoftwareRasterizerRGBA(to)
	rat.SetAntiAliasing(gorat.AntiAliasing4x)
	rat.MoveTo([2]float32{50,50}, nil)
	rat.QuadTo(gorat.Vec2(670, 50), gorat.Vec2(670, 670), nil)
	rat.LineTo(gorat.Vec2(50, 670), nil)
	//rat.LineTo([2]float32{90,170}, nil)
	// rat.Close()
	// .Close() call auto
	rat.FillColor()
}