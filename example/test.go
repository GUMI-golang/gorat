package main

import (
	"github.com/GUMI-golang/gorat"
	"github.com/GUMI-golang/gumi/gcore"
	"fmt"
	"os"
	"image/jpeg"
)

func main() {
	f, err := os.Open("example/cubes_1024.jpg")
	if err != nil {
		panic(err)
	}
	tex, err := jpeg.Decode(f)
	if err != nil {
		panic(err)
	}
	filler := gorat.NewImageFiller(tex, gorat.ImageFillerFixed)

	r := gorat.NewSoftware(gcore.FixedSize{400, 400}.Rect())
	r.SetFiller(gorat.NewColorFiller(255,0,0, 255))
	r.SetFiller(filler)
	//
	r.MoveTo(gorat.Vec2(20, 20))
	r.LineTo(gorat.Vec2(380, 380))
	r.LineTo(gorat.Vec2(20, 380))
	r.CloseTo()
	r.MoveTo(gorat.Vec2(40, 40))
	r.LineTo(gorat.Vec2(40, 360))
	r.LineTo(gorat.Vec2(360, 360))
	r.QuadTo(gorat.Vec2(360, 40), gorat.Vec2(40, 40))
	r.CloseTo()
	//r.Print()
	r.Fill()
	fmt.Println("==========================================")
	//r.Print()
	gcore.Capture("aout2", r.Image())
}
