package main

import (
	"github.com/GUMI-golang/gorat"
	"github.com/GUMI-golang/gumi/gcore"
	"github.com/go-gl/mathgl/mgl32"
	"image"
	"os"
)

func main() {
	rat := gorat.NewSoftwareRasterizer(100, 100)
	tex, _, err := image.Decode(gcore.MustValue(os.Open("dirt.jpg")).(*os.File))
	if err != nil {
		panic(err)
	}
	rat.TrUV(gorat.Triangle{
		A: mgl32.NewVecNFromData([]float32{10, 10}).Vec2(),
		B: mgl32.NewVecNFromData([]float32{90, 10}).Vec2(),
		C: mgl32.NewVecNFromData([]float32{10, 90}).Vec2(),
	}, [3]mgl32.Vec2{
		{0, 0},
		{1, 0},
		{0, 1},
	}, tex)
	gcore.Capture("aout", rat.Image())
}
