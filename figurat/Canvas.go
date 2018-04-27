package figurat

import (
	"github.com/GUMI-golang/gorat"
	"github.com/go-gl/mathgl/mgl32"
	"math"
)

type Canvas struct {
	raster  gorat.VectorDrawer
	backup  [][]mgl32.Vec2
	working []mgl32.Vec2
}

func NewCanvas(rasterizer gorat.VectorDrawer) Canvas {
	return Canvas{
		raster: rasterizer,
	}
}

func (s *Canvas) MoveTo(to mgl32.Vec2) {
	const pointCapacity = 32
	s.working = make([]mgl32.Vec2, 0, pointCapacity)
	s.working = append(s.working, to)
}
func (s *Canvas) LineTo(to mgl32.Vec2) {
	s.working = append(s.working, to)
}
func (s *Canvas) QuadTo(pivot, to mgl32.Vec2) {
	// Come from golang.org/vecx/image/vector Raster.QuadTo
	from := s.working[len(s.working)-1]
	devsq := gorat.DevSquared(from, pivot, to)
	if devsq >= 0.333 {
		const tol = 3
		n := 1 + int(math.Sqrt(math.Sqrt(tol*float64(devsq))))
		t, nInv := float32(0), 1/float32(n)
		for i := 0; i < n-1; i++ {
			t += nInv
			ab := gorat.Lerp(t, from, pivot)
			bc := gorat.Lerp(t, pivot, to)
			s.LineTo(gorat.Lerp(t, ab, bc))
		}
	}
	s.LineTo(to)
}
func (s *Canvas) CubeTo(pivot1, pivot2, to mgl32.Vec2) {
	// Come from golang.org/vecx/image/vector Raster.QuadTo
	from := s.working[len(s.working)-1]
	devsq := gorat.DevSquared(from, pivot1, to)
	if devsqAlt := gorat.DevSquared(from, pivot2, to); devsq < devsqAlt {
		devsq = devsqAlt
	}
	if devsq >= 0.333 {
		const tol = 3
		n := 1 + int(math.Sqrt(math.Sqrt(tol*float64(devsq))))
		t, nInv := float32(0), 1/float32(n)
		for i := 0; i < n-1; i++ {
			t += nInv
			ab := gorat.Lerp(t, from, pivot1)
			bc := gorat.Lerp(t, pivot1, pivot2)
			cd := gorat.Lerp(t, pivot2, to)
			abc := gorat.Lerp(t, ab, bc)
			bcd := gorat.Lerp(t, bc, cd)
			s.LineTo(gorat.Lerp(t, abc, bcd))
		}
	}
	s.LineTo(to)
}
func (s *Canvas) CloseTo() {
	s.backup = append(s.backup, s.working)
	s.working = nil
}