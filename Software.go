package gorat

import (
	"fmt"
	"github.com/go-gl/mathgl/mgl32"
	"image"
	"math"
)

type Software struct {
	// result after fill
	img *image.RGBA
	// workspace
	bufPix      []float32
	start, last mgl32.Vec2
	//
	Options
}

func NewSoftware(size image.Rectangle) *Software {
	res := new(Software)
	res.img = image.NewRGBA(size)
	res.DefaultOption()
	res.Reset()
	return res
}
func (s *Software) Reset() {
	s.start = Vec2(0, 0)
	s.last = Vec2(0, 0)
	if n := s.img.Rect.Dx() * s.img.Rect.Dy(); n > cap(s.bufPix) {
		s.bufPix = make([]float32, n)
	} else {
		s.bufPix = s.bufPix[:n]
		for i := range s.bufPix {
			s.bufPix[i] = 0
		}
	}

}

const (
	vecx = 0
	vecy = 1
	//
	r = 0
	g = 1
	b = 2
	a = 3
)

func (s *Software) MoveTo(p mgl32.Vec2) {
	s.start = p
	s.last = p
}

const closeToZero = 0.0000001

func fromTo(from, to float32, max int) (a, b int) {
	a, b = floorInt(from), ceilInt(to)
	if b > max {
		b = max
	}
	return
}
func (s *Software) LineTo(p mgl32.Vec2) {
	// Point setup
	from := s.last
	to := p
	s.last = p
	size := s.img.Rect.Size()
	//
	// direction setup
	var dir float32 = 1
	if from[vecy] > to[vecy] {
		dir, from, to = -1, to, from
	}
	if to[vecy]-from[vecy] <= closeToZero {
		return
	}
	// delta xCurr/ delta y
	Δxy := (to[vecx] - from[vecx]) / (to[vecy] - from[vecy])
	xCurr := from[vecx]
	yFrom, yTo := fromTo(from[vecy], to[vecy], size.Y)
	//
	for y := yFrom; y < yTo; y++ {
		Δy := min(float32(y+1), to[vecy]) - max(float32(y), from[vecy])
		xNext := xCurr + float32(Δy*Δxy) // = xCurr + Δx
		if y < 0 {
			xCurr = xNext
			continue
		}
		buf := s.bufPix[y*size.X:]
		d := float32(Δy * dir)
		x0, x1 := xCurr, xNext
		if xCurr > xNext {
			x0, x1 = x1, x0
		}
		x0i := floorInt(x0)
		x0Floor := float32(x0i)
		x1i := ceilInt(x1)
		x1Ceil := float32(x1i)
		if x1i <= x0i+1 {
			xmf := float32(0.5*(xCurr+xNext)) - x0Floor
			if i := iclamp(x0i+0, 0, size.X); i < len(buf) {
				buf[i] += d - float32(d*xmf)
			}
			if i := iclamp(x0i+1, 0, size.X); i < len(buf) {
				buf[i] += float32(d * xmf)
			}
		} else {
			s := 1 / (x1 - x0)
			x0f := x0 - x0Floor
			oneMinusX0f := 1 - x0f
			a0 := float32(0.5 * s * oneMinusX0f * oneMinusX0f)
			x1f := x1 - x1Ceil + 1
			am := float32(0.5 * s * x1f * x1f)

			if i := iclamp(x0i, 0, size.X); i < len(buf) {
				buf[i] += float32(d * a0)
			}

			if x1i == x0i+2 {
				if i := iclamp(x0i+1, 0, size.X); i < len(buf) {
					buf[i] += float32(d * (1 - a0 - am))
				}
			} else {
				a1 := float32(s * (1.5 - x0f))
				if i := iclamp(x0i+1, 0, size.X); i < len(buf) {
					buf[i] += float32(d * (a1 - a0))
				}
				dTimesS := float32(d * s)
				for xi := x0i + 2; xi < x1i-1; xi++ {
					if i := iclamp(xi, 0, size.X); i < len(buf) {
						buf[i] += dTimesS
					}
				}
				a2 := a1 + float32(s*float32(x1i-x0i-3))
				if i := iclamp(x1i-1, 0, size.X); i < len(buf) {
					buf[i] += float32(d * (1 - a2 - am))
				}
			}

			if i := iclamp(x1i, 0, size.X); i < len(buf) {
				buf[i] += float32(d * am)
			}
		}
		xCurr = xNext
	}
}
func (s *Software) QuadTo(pivot, to mgl32.Vec2) {
	from := s.last
	devsq := devSquared(from, pivot, to)
	if devsq >= 0.333 {
		const tol = 3
		n := 1 + int(math.Sqrt(math.Sqrt(tol*float64(devsq))))
		t, nInv := float32(0), 1/float32(n)
		for i := 0; i < n-1; i++ {
			t += nInv
			fromPivot := lerp(t, from, pivot)
			pivotTo := lerp(t, pivot, to)
			s.LineTo(lerp(t, fromPivot, pivotTo))
		}
	}
	s.LineTo(to)
}
func (s *Software) CubeTo(pivot1, pivot2, to mgl32.Vec2) {
	from := s.last
	devsq := devSquared(from, pivot1, to)
	if devsqAlt := devSquared(from, pivot2, to); devsq < devsqAlt {
		devsq = devsqAlt
	}
	if devsq >= 0.333 {
		const tol = 3
		n := 1 + int(math.Sqrt(math.Sqrt(tol*float64(devsq))))
		t, nInv := float32(0), 1/float32(n)
		for i := 0; i < n-1; i++ {
			t += nInv
			ab := lerp(t, from, pivot1)
			bc := lerp(t, pivot1, pivot2)
			cd := lerp(t, pivot2, to)
			abc := lerp(t, ab, bc)
			bcd := lerp(t, bc, cd)
			s.LineTo(lerp(t, abc, bcd))
		}
	}
	s.LineTo(to)
}
func (s *Software) CloseTo() {
	s.LineTo(s.start)
}
func (s *Software) Print() {
	w := s.img.Rect.Dx()
	h := s.img.Rect.Dy()
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			fmt.Printf("%5.2f ", s.bufPix[x+w*y])
		}
		fmt.Println()
	}

}
func (s *Software) Fill() {
	switch t := s.filler.(type) {
	case ColorFiller:
		s.uniformFill(float32(t.R), float32(t.G), float32(t.B), float32(t.A))
		//case _ImageFillerFixed:
	case *_ImageFillerFixed:
		s.fixedFill(t.img)
	default:
		s.fillerFill()

	}
	return
}

const almostZero float32 = 0.000001
func (s *Software) uniformFill(cr, cg, cb, ca float32) {
	acc := float32(0)
	for i, v := range s.bufPix {
		acc += v
		a := acc
		if a < 0 {
			a = -a
		}
		if a > 1 {
			a = 1
		}

		s.bufPix[i] = a
	}
	width := s.img.Rect.Dx()
	for x := s.img.Rect.Min.X; x < s.img.Rect.Max.X; x++ {
		for y := s.img.Rect.Min.Y; y < s.img.Rect.Max.Y; y++ {
			pixOffset := s.img.PixOffset(x, y)
			bufOffset := (x - s.img.Rect.Min.X) + width*(y-s.img.Rect.Min.Y)
			buf := s.bufPix[bufOffset]
			if buf < almostZero{
				continue
			}
			s.img.Pix[pixOffset+r] = uint8(iclamp(int(buf*cr), 0, 255))
			s.img.Pix[pixOffset+g] = uint8(iclamp(int(buf*cg), 0, 255))
			s.img.Pix[pixOffset+b] = uint8(iclamp(int(buf*cb), 0, 255))
			s.img.Pix[pixOffset+a] = uint8(iclamp(int(buf*ca), 0, 255))
		}
	}

	return
}
func (s *Software) fixedFill(src *image.RGBA) {
	acc := float32(0)
	// Draw Rect
	fillBound := image.Rectangle{
		Min: image.Point{
			X: math.MaxInt32,
			Y: math.MaxInt32,
		},
		Max:image.Point{
			X: math.MinInt32,
			Y: math.MinInt32,
		},
	}
	width := s.img.Rect.Dx()
	for i, v := range s.bufPix {
		acc += v
		a := acc
		if a < 0 {
			a = -a
		}
		if a > 1 {
			a = 1
		}
		//
		if a > almostZero {
			y := i / width
			x := i % width
			if x < fillBound.Min.X{
				fillBound.Min.X = x
			}
			if fillBound.Max.X < x{
				fillBound.Max.X = x
			}
			if y < fillBound.Min.Y{
				fillBound.Min.Y = y
			}
			if fillBound.Max.Y < y{
				fillBound.Max.Y = y
			}
		}
		s.bufPix[i] = a
	}
	if fillBound.Max.X < 0{
		return
	}
	//
	for x := s.img.Rect.Min.X; x < s.img.Rect.Max.X; x++ {
		for y := s.img.Rect.Min.Y; y < s.img.Rect.Max.Y; y++ {
			pixOffset := s.img.PixOffset(x, y)
			bufOffset := (x) + width*(y)
			buf := s.bufPix[bufOffset]
			if buf < almostZero{
				continue
			}
			srcX, srcY := x - fillBound.Min.X, y-fillBound.Min.Y
			if src.Rect.Min.X <= srcX && srcX < src.Rect.Max.X && src.Rect.Min.Y <= srcY && srcY < src.Rect.Max.Y{
				srcOffset := src.PixOffset(srcX, srcY)
				s.img.Pix[pixOffset+r] = uint8(iclamp(int(buf*float32(src.Pix[srcOffset +r])), 0, 255))
				s.img.Pix[pixOffset+g] = uint8(iclamp(int(buf*float32(src.Pix[srcOffset +g])), 0, 255))
				s.img.Pix[pixOffset+b] = uint8(iclamp(int(buf*float32(src.Pix[srcOffset +b])), 0, 255))
				s.img.Pix[pixOffset+a] = uint8(iclamp(int(buf*float32(src.Pix[srcOffset +a])), 0, 255))
			}
		}
	}

	return
}
func (s *Software) fillerFill() {
	acc := float32(0)
	// Draw Rect
	fillBound := image.Rectangle{
		Min: image.Point{
			X: math.MaxInt32,
			Y: math.MaxInt32,
		},
		Max:image.Point{
			X: math.MinInt32,
			Y: math.MinInt32,
		},
	}
	width := s.img.Rect.Dx()
	for i, v := range s.bufPix {
		acc += v
		a := acc
		if a < 0 {
			a = -a
		}
		if a > 1 {
			a = 1
		}
		//
		if a > 0 {
			y := i / width
			x := i % width
			if x < fillBound.Min.X{
				fillBound.Min.X = x
			}
			if fillBound.Max.X < x{
				fillBound.Max.X = x
			}
			if y < fillBound.Min.Y{
				fillBound.Min.X = y
			}
			if fillBound.Max.Y < y{
				fillBound.Max.Y = y
			}
		}
		s.bufPix[i] = a
	}
	if fillBound.Max.X < 0{
		return
	}
	//
	for x := s.img.Rect.Min.X; x < s.img.Rect.Max.X; x++ {
		for y := s.img.Rect.Min.Y; y < s.img.Rect.Max.Y; y++ {
			pixOffset := s.img.PixOffset(x, y)
			bufOffset := (x - s.img.Rect.Min.X) + width*(y-s.img.Rect.Min.Y)
			buf := s.bufPix[bufOffset]
			if buf < almostZero{
				continue
			}
			c := s.filler.RGBA(x - fillBound.Min.X, y-fillBound.Min.Y)
			if c.A == 0{
				continue
			}
			s.img.Pix[pixOffset+r] = uint8(iclamp(int(buf*float32(c.R)), 0, 255))
			s.img.Pix[pixOffset+g] = uint8(iclamp(int(buf*float32(c.G)), 0, 255))
			s.img.Pix[pixOffset+b] = uint8(iclamp(int(buf*float32(c.B)), 0, 255))
			s.img.Pix[pixOffset+a] = uint8(iclamp(int(buf*float32(c.A)), 0, 255))
		}
	}

	return
}
func (s *Software) Image() image.Image {
	return s.img
}
