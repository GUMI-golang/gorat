package main

import (
	"github.com/GUMI-golang/gorat"
	"github.com/GUMI-golang/gumi/gcore"
	"runtime"
	"github.com/GUMI-golang/gorat/fwrat"
	"github.com/GUMI-golang/gorat/oglSupport/v43"
	"fmt"
)

var width, height = 256, 256
func main() {
	runtime.LockOSThread()
	// setup driver
	gcore.Must(gorat.SetupDriver(v43.Driver()))
	// image loading
	//cube := gcore.MustValue(os.Open("example/cubes_64.png")).(*os.File)
	//defer cube.Close()
	//img := gcore.MustValue(png.Decode(cube)).(image.Image)
	//
	// screen setup
	ctx := fwrat.OffscreenContext(width, height)
	fmt.Println("Context", ctx)
	// Make target texture
	// It use driver, but if you need you can use your GL_TEXTURE_2D, GL_RGBA32F image
	// Like
	// res := gorat.HardwareResult(<your image uint32(gl pointer) here>)
	res := gorat.Driver().Result(width, height)
	defer res.Delete()
	// gorat hardware delete gl object when grabage collecter remove *Hardware object
	hw0 := gorat.NewHardware(res)
	// filling
	hw0.MoveTo(gorat.Vec2(32,32))
	hw0.LineTo(gorat.Vec2(32, float32(height-32)))
	hw0.LineTo(gorat.Vec2(float32(width-32), float32(height-32)))
	hw0.LineTo(gorat.Vec2(32, float32(height-32)))
	//hw0.SetFiller(gorat.NewImageFiller(img, gorat.ImageFillerGausian))
	hw0.Fill()
	// Stroking
	//hw0.MoveTo(gorat.Vec2(32,32))
	//hw0.LineTo(gorat.Vec2(32, float32(height-32)))
	//hw0.LineTo(gorat.Vec2(float32(width)/2, float32(height)/2))
	//hw0.SetStrokeWidth(4)
	//hw0.SetStrokeJoin(gorat.StrokeJoinMiter)
	//hw0.SetStrokeCap(gorat.StrokeCapRound)
	//hw0.Stroke()


	// save result
	gcore.Capture("aout0", res.Get())
}
