package main
//
//import (
//	"image"
//	"github.com/GUMI-golang/gumi/gcore"
//	"github.com/GUMI-golang/gorat/example/GLRasterizer"
//	"runtime"
//)
//
////var SRExamples = []func(rgba *image.rgba){
////	// Triangle
////	SoftwareRasterizer.Triangle,
////	SoftwareRasterizer.TriangleSubImage,
////	SoftwareRasterizer.TriangleCWCCWOverlap,
////	// Line
////	SoftwareRasterizer.LineOne,
////	// Quad
////	SoftwareRasterizer.Quad,
////}
//
//var GLExamples = []func(rgba *image.RGBA){
//	GLRasterizer.Simple,
//	GLRasterizer.Simple2,
//}
//
//var (
//	size = gcore.FixedSize{1280, 720}
//	index = 1
//	out = "aout"
//)
//func main() {
//	runtime.LockOSThread()
//	var rgba *image.RGBA
//	//var result testing.BenchmarkResult
//	rgba = image.NewRGBA(size.Rect())
//	GLExamples[index](rgba)
//	//result = testing.Benchmark(func(b *testing.B) {
//	//	SoftwareRasterizerPrototype.Quad(rgba)
//	//})
//	//fmt.Println(result.String())
//	//gcore.Capture(out + "gorat", rgba)
//	//
//	//rgba = image.NewRGBA(size.Rect())
//	//result = testing.Benchmark(func(b *testing.B) {
//	//	ctx := gg.NewContextForRGBA(rgba)
//	//	ctx.SetColor(image.Black)
//	//	ctx.MoveTo(50, 50)
//	//	ctx.QuadraticTo(670, 50, 670, 670)
//	//	ctx.LineTo(50, 670)
//	//	ctx.Fill()
//	//})
//	//fmt.Println(result.String())
//	//gcore.Capture(out + "gg", rgba)
//	//
//	//rgba = image.NewRGBA(size.Rect())
//	//blackUniform := image.NewUniform(image.Black)
//	//result = testing.Benchmark(func(b *testing.B) {
//	//	ctx := vector.NewRasterizer(size.Size())
//	//
//	//	ctx.MoveTo(50, 50)
//	//	ctx.QuadTo(670, 50, 670, 670)
//	//	ctx.LineTo(50, 670)
//	//	ctx.ClosePath()
//	//	ctx.Draw(rgba, rgba.Rect, blackUniform, image.ZP)
//	//})
//	//fmt.Println(result.String())
//	//gcore.Capture(out + "vector", rgba)
//}
//
