package gorat

import (
	"math"
	"github.com/go-gl/mathgl/mgl32"
)


var nanvec2 = mgl32.Vec2{float32(math.NaN()), float32(math.NaN())}

func quadFromTo(from, pivot, to mgl32.Vec2) (res []mgl32.Vec2) {
	devsq := DevSquared(from, pivot, to)
	if devsq >= 0.333 {
		const tol = 3
		n := 1 + int(math.Sqrt(math.Sqrt(tol*float64(devsq))))
		t, nInv := float32(0), 1/float32(n)
		for i := 0; i < n-1; i++ {
			t += nInv
			fromPivot := Lerp(t, from, pivot)
			pivotTo := Lerp(t, pivot, to)
			res = append(res, Lerp(t, fromPivot, pivotTo))
		}
	}
	res = append(res, to)
	return res
}
func cubeFromTo(from, pivot1, pivot2, to mgl32.Vec2) (res []mgl32.Vec2) {
	devsq := DevSquared(from, pivot1, to)
	if devsqAlt := DevSquared(from, pivot2, to); devsq < devsqAlt {
		devsq = devsqAlt
	}
	if devsq >= 0.333 {
		const tol = 3
		n := 1 + int(math.Sqrt(math.Sqrt(tol*float64(devsq))))
		t, nInv := float32(0), 1/float32(n)
		for i := 0; i < n-1; i++ {
			t += nInv
			ab := Lerp(t, from, pivot1)
			bc := Lerp(t, pivot1, pivot2)
			cd := Lerp(t, pivot2, to)
			abc := Lerp(t, ab, bc)
			bcd := Lerp(t, bc, cd)
			res = append(res, Lerp(t, abc, bcd))
		}
	}
	res = append(res, to)
	return res
}