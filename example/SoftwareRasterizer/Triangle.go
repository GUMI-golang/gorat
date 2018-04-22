package SoftwareRasterizer

import (
	"image"
	"github.com/GUMI-golang/gorat"
	"image/draw"
)

func Triangle(to *image.RGBA) {
	// Software rasterizer support RGBA Direct and RGBA SubImage
	rat := gorat.NewSoftwareRasterizerRGBA(to)
	//
	rat.MoveTo([2]float32{10,10}, nil)
	rat.LineTo([2]float32{90,90}, nil)
	rat.LineTo([2]float32{10,90}, nil)
	// rat.Close()
	// .Close() call auto
	rat.FillColor()
}
func TriangleSubImage(to *image.RGBA) {
	// Software rasterizer support RGBA Direct and RGBA SubImage

	half := to.SubImage(image.Rect(
		to.Rect.Min.X + to.Rect.Dx() / 4,
		to.Rect.Min.Y + to.Rect.Dy() / 4,
		to.Rect.Max.X - to.Rect.Dx() / 4,
		to.Rect.Max.Y - to.Rect.Dy() / 4,
	)).(*image.RGBA)
	draw.Draw(half, half.Rect, image.NewUniform(image.White), image.ZP, draw.Src)
	rat := gorat.NewSoftwareRasterizerRGBA(half)
	//
	rat.MoveTo([2]float32{10,10}, nil)
	rat.LineTo([2]float32{90,90}, nil)
	rat.LineTo([2]float32{10,90}, nil)
	// rat.Close()
	// .Close() call auto
	rat.FillColor()
}
func TriangleCWCCWOverlap(to *image.RGBA) {
	rat := gorat.NewSoftwareRasterizerRGBA(to)
	// CW Triangle
	rat.MoveTo([2]float32{10,20}, nil)
	rat.LineTo([2]float32{90,100}, nil)
	rat.LineTo([2]float32{10,100}, nil)
	//rat.Close()
	// CCW Triangle
	rat.MoveTo([2]float32{20,10}, nil)
	rat.LineTo([2]float32{20,90}, nil)
	rat.LineTo([2]float32{100,90}, nil)
	//rat.Close()
	rat.FillColor()
}
func TriangleAA(to *image.RGBA) {
	// Software rasterizer support RGBA Direct and RGBA SubImage
	rat := gorat.NewSoftwareRasterizerRGBA(to)
	rat.SetAntiAliasing(gorat.AntiAliasing4x)
	//
	rat.MoveTo([2]float32{10,10}, nil)
	rat.LineTo([2]float32{90,90}, nil)
	rat.LineTo([2]float32{10,90}, nil)
	// rat.Close()
	// .Close() call auto
	rat.FillColor()
}