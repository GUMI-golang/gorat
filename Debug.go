package gorat

import (
	"image"
	"image/png"
	"os"
)

var DEBUG _Debug
type _Debug struct {}

func (_Debug) get2Image(w, h float32) [2]*image.RGBA {
	var temp [2]*image.RGBA
	temp[0] = image.NewRGBA(image.Rect(0,0, int(w), int(h)))
	temp[1] = image.NewRGBA(image.Rect(0,0, int(w), int(h)))
	return temp
}
func (_Debug) save2Image(temp [2]*image.RGBA, n0, n1 string)  {
	f0, err := os.OpenFile(n0+".png", os.O_CREATE | os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	f1, err := os.OpenFile(n1+".png", os.O_CREATE | os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	png.Encode(f0, temp[0])
	png.Encode(f1, temp[1])
}
func (_Debug) FillToFile(rasterizer VectorDrawer, stroking, filling string) {
	temp := DEBUG.get2Image(rasterizer.Size())
	rasterizer.debugFill(temp[0], temp[1])
	DEBUG.save2Image(temp, stroking, filling)
}
func (_Debug) StrokeToFile(rasterizer VectorDrawer, stroking, filling string) {
	temp := DEBUG.get2Image(rasterizer.Size())
	rasterizer.debugStroke(temp[0], temp[1])
	DEBUG.save2Image(temp, stroking, filling)
}
