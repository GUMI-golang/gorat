package main

import (
	"image"
	"github.com/GUMI-golang/gorat/example/SoftwareRasterizer"
	"github.com/GUMI-golang/gumi/gcore"
	"testing"
	"fmt"
	"golang.org/x/image/vector"
)

var SRExamples = []func(rgba *image.RGBA){
	// Triangle
	SoftwareRasterizer.Triangle,
	SoftwareRasterizer.TriangleSubImage,
	SoftwareRasterizer.TriangleCWCCWOverlap,
	// Line
	SoftwareRasterizer.LineOne,
	// Quad
	SoftwareRasterizer.Quad,
}

var (
	size = gcore.FixedSize{1280, 720}
	index = 0
	out = "aout"
)
func main() {
	//var rgba *image.RGBA
	//var result testing.BenchmarkResult

	//
	//rgba = image.NewRGBA(size.Rect())
	//result = testing.Benchmark(func(b *testing.B) {
	//	SoftwareRasterizerPrototype.Quad(rgba)
	//})
	//fmt.Println(result.String())
	//gcore.Capture(out + "gorat", rgba)
	//
	//rgba = image.NewRGBA(size.Rect())
	//result = testing.Benchmark(func(b *testing.B) {
	//	ctx := gg.NewContextForRGBA(rgba)
	//	ctx.SetColor(image.Black)
	//	ctx.MoveTo(50, 50)
	//	ctx.QuadraticTo(670, 50, 670, 670)
	//	ctx.LineTo(50, 670)
	//	ctx.Fill()
	//})
	//fmt.Println(result.String())
	//gcore.Capture(out + "gg", rgba)
	//
	//rgba = image.NewRGBA(size.Rect())
	//blackUniform := image.NewUniform(image.Black)
	//result = testing.Benchmark(func(b *testing.B) {
	//	ctx := vector.NewRasterizer(size.Size())
	//
	//	ctx.MoveTo(50, 50)
	//	ctx.QuadTo(670, 50, 670, 670)
	//	ctx.LineTo(50, 670)
	//	ctx.ClosePath()
	//	ctx.Draw(rgba, rgba.Rect, blackUniform, image.ZP)
	//})
	//fmt.Println(result.String())
	//gcore.Capture(out + "vector", rgba)
}

