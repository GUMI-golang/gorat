package GLRasterizer

import (
	"image"
	"image/draw"
	"image/color"
)

var test *image.RGBA

func init() {
	test = image.NewRGBA(image.Rect(0,0, 128, 128))
	draw.Draw(test, test.Rect, image.NewUniform(color.RGBA{255, 0, 0, 255}), image.ZP, draw.Src)
	for x := 32; x < 96; x ++{
		for y := 0; y < 64; y ++{
			test.Set(x, y, color.RGBA{0,255, 0, 255})
		}
	}
}
