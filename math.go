package gorat

import (
	"github.com/go-gl/mathgl/mgl32"
	"math"
)

func minimal(f32s ... float32) (min float32) {
	if len(f32s) < 1{
		panic("At least one data")
	}
	min = f32s[0]
	for _, v := range f32s {
		if v < min{
			min = v
		}
	}
	return
}
func maximum(f32s ... float32) (max float32) {
	if len(f32s) < 1{
		panic("At least one data")
	}
	max = f32s[0]
	for _, v := range f32s {
		if v > max{
			max = v
		}
	}
	return
}
func floorInt(f32 float32) int {
	return int(f32)
}
func ceilInt(f32 float32) int {
	return int(f32 + 0.999999999)
}

func devSquared(a, b, c mgl32.Vec2) float32 {
	devx := a[0] - 2*b[0] + c[0]
	devy := a[1] - 2*b[1] + c[1]
	return devx*devx + devy*devy
}
func lerp(t float32, p, q mgl32.Vec2) mgl32.Vec2 {
	return [2]float32{p[0] + t*(q[0]-p[0]), p[1] + t*(q[1]-p[1])}
}
func floor(f32 float32) float32 {
	return float32(math.Floor(float64(f32)))
}
func ceil(f32 float32) float32 {
	return float32(math.Ceil(float64(f32)))
}
func min(a,b float32) float32 {
	if a < b{
		return a
	}
	return b
}
func max(a,b float32) float32 {
	if a < b{
		return b
	}
	return a
}
func iclamp(a, min, max int) int {
	if a < min{
		return min
	}
	if a > max{
		return max
	}
	return a
}