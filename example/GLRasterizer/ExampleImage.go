package main

import (
	"image"
	"image/draw"
	"image/color"
)

var test *image.RGBA

func init() {
	test = image.NewRGBA(image.Rect(0,0, 128, 128))
	draw.Draw(test, test.Rect, image.NewUniform(color.RGBA{255, 0, 0, 255}), image.ZP, draw.Src)
}
