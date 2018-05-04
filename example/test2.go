package main

import (
	"github.com/GUMI-golang/gorat"
	"github.com/GUMI-golang/gumi/gcore"
	"runtime"
	"github.com/GUMI-golang/gorat/fwrat"
	"github.com/GUMI-golang/gorat/oglSupport/v43"
	"fmt"
	"image"
	"github.com/GUMI-golang/gorat/textrat"
	"image/color"
)

var width, height = 800, 600
func main() {
	runtime.LockOSThread()
	// setup driver
	octx := gcore.MustValue(fwrat.CreateOffscreenContext(v43.Driver())).(*fwrat.Offscreen)
	defer octx.Delete()
	gorat.Use(octx)
	res := octx.Driver().Result(width, height)
	defer res.Delete()
	// gorat hardware delete gl object when grabage collecter remove *Hardware object
	hw0 := gorat.NewHardware(res)
	// filling
	hw0.SetFiller(gorat.NewColorFiller(color.RGBA{255, 0, 0, 255}))
	hw0.MoveTo(gorat.Vec2(0,0))
	hw0.LineTo(gorat.Vec2(float32(width),0))
	hw0.LineTo(gorat.Vec2(float32(width), float32(height)))
	hw0.LineTo(gorat.Vec2(0, float32(height)))
	hw0.Fill()
	sub := hw0.SubRasterizer(image.Rect(32, 32, 800 - 32, 600-32))
	bd := sub.Bound()
	bd = bd.Sub(bd.Min)
	fmt.Println(bd)
	textrat.Default.SetSize(32)
	sub.SetFiller(gorat.NewColorFiller(color.RGBA{0, 255, 0, 128}))
	textrat.Default.TextInRect(sub, "Hello?", bd , gcore.AlignCenter)
	// save result
	gcore.Capture("aout0", res.Get())
}
