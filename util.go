package gorat

import (
	"golang.org/x/image/math/fixed"
	"image"
	"github.com/go-gl/mathgl/mgl32"
	"image/color"
)

const (
	modifier64 = 1 << 12
)


func Fixed32ToFixed64(int26_6 fixed.Int26_6) fixed.Int52_12 {

	return fixed.Int52_12(int26_6) << 6
}
func Fixed64ToFixed32(i fixed.Int52_12) fixed.Int26_6 {
	return fixed.Int26_6(i >> 6)
}
//func Fixed32Tofloat32(int26_6 fixed.Int26_6) float32 {
//	return float32(int26_6) / modifier32
//}
func Fixed64ToFloat64(i fixed.Int52_12) float64 {
	return float64(i) / modifier64
}
func fixed64sToFloat64s(fixes ... fixed.Int52_12) []float64 {
	temp := make([]float64, len(fixes))
	for i, v := range fixes {
		temp[i] = Fixed64ToFloat64(v)
	}
	return temp
}
func Float64ToFixed64(f float64) fixed.Int52_12{
	return fixed.Int52_12(int64(f * modifier64))
}

func Float32sToFloat64s(float32s ...float32) []float64 {
	res := make([]float64, len(float32s))
	for i, v := range float32s {
		res[i] = float64(v)
	}
	return res
}

func FixedRectToRect(rect fixed.Rectangle52_12) image.Rectangle {
	return image.Rect(
		rect.Min.X.Round(),
		rect.Min.Y.Round(),
		rect.Max.X.Round(),
		rect.Max.Y.Round(),
	)
}
func RectToFixedRect(rect image.Rectangle) fixed.Rectangle52_12 {
	return fixed.Rectangle52_12{
		Min: Pt(rect.Min.X, rect.Min.Y),
		Max: Pt(rect.Max.X, rect.Max.Y),
	}
}

//func ChangeRect(rect fixed.Rectangle52_12) image.Rectangle {
//	return image.ChangeRect(
//		rect.Min.X.Round(),
//		rect.Min.Y.Round(),
//		rect.Max.X.Round(),
//		rect.Max.Y.Round(),
//	)
//}
func I(x int) fixed.Int52_12 {
	return fixed.Int52_12(x) << 12
}
func Pt(x, y int) fixed.Point52_12 {
	return fixed.Point52_12{
		X: fixed.Int52_12(x << 12),
		Y: fixed.Int52_12(y << 12),
	}
}

func Clamp(v, min, max float32) float32 {
	if v < min{
		return min
	}
	if v > max{
		return max
	}
	return v
}
func Vec4Color(c mgl32.Vec4) color.RGBA {
	return color.RGBA{
		R: uint8(Clamp(c[0], 0, 1) *255),
		G: uint8(Clamp(c[1], 0, 1)*255),
		B: uint8(Clamp(c[2], 0, 1)*255),
		A: uint8(Clamp(c[3], 0, 1)*255),
	}
}
func ColorVec4(c color.RGBA) mgl32.Vec4 {
	return mgl32.Vec4{
		float32(c.R) / 255,
		float32(c.G) / 255,
		float32(c.B) / 255,
		float32(c.A) / 255,
	}
}
func Vec2(x, y float32) mgl32.Vec2 {
	return mgl32.Vec2{x, y}
}
func Vec4(x, y, z, w float32) mgl32.Vec4 {
	return mgl32.Vec4{x, y, z, w}
}