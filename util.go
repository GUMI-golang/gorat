package gorat

import (
	"github.com/go-gl/mathgl/mgl32"
)


func Vec2(x, y float32) mgl32.Vec2 {
	return mgl32.Vec2{x, y}
}
func Vec4(x, y, z, w float32) mgl32.Vec4 {
	return mgl32.Vec4{x, y, z, w}
}
const (
	vecx = 0
	vecy = 1
	//
	pixr = 0
	pixg = 1
	pixb = 2
	pixa = 3
)