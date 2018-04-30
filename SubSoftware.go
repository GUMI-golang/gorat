package gorat

//import (
//	"github.com/go-gl/mathgl/mgl32"
//	"image"
//	"image/color"
//	"image/draw"
//	"math"
//)
//
//type SoftwareSub struct {
//	root  *SoftwareRoot
//	bound image.Rectangle
//	// workspace
//	bufPix      []float32
//	start, last mgl32.Vec2
//	//
//	Options
//}
//func (s *SoftwareSub) PreviousPoint() mgl32.Vec2 {
//	panic("implement me")
//}
//
//func (s *SoftwareSub) Point() mgl32.Vec2 {
//	panic("implement me")
//}
//
//func (s *SoftwareSub) Stroke() {
//	panic("implement me")
//}
//func (s *SoftwareSub) Root() Rasterizer {
//	return s.root
//}
//func (s *SoftwareSub) Bound() image.Rectangle {
//	return s.bound
//}
//func (s *SoftwareSub) SubRasterizer(pixr image.Rectangle) SubRasterizer {
//	pixr = pixr.Add(s.bound.Min)
//	pixr = s.bound.Intersect(pixr)
//	temp := &SoftwareSub{
//		root:  s.root,
//		bound: pixr,
//	}
//	temp.Restore(s.Clone())
//	temp.Reset()
//	return temp
//}
//func (s *SoftwareSub) Clear() {
//	draw.Draw(s.root.img, s.bound, image.NewUniform(color.Transparent), image.ZP, draw.Src)
//}
//func (s *SoftwareSub) Size() (w, h float32) {
//	return float32(s.bound.Dx()), float32(s.bound.Dy())
//}
//func (s *SoftwareSub) Reset() {
//	s.start = Vec2(0, 0)
//	s.last = Vec2(0, 0)
//	if n := s.bound.Dx() * s.bound.Dy(); n > cap(s.bufPix) {
//		s.bufPix = make([]float32, n)
//	} else {
//		s.bufPix = s.bufPix[:n]
//		for i := range s.bufPix {
//			s.bufPix[i] = 0
//		}
//	}
//
//}
//

//
//func (s *SoftwareSub) MoveTo(p mgl32.Vec2) {
//	s.start = p
//	s.last = p
//}
//
//const closeToZero = 0.0000001
//
//func fromTo(from, to float32, max int) (pixa, pixb int) {
//	pixa, pixb = floorInt(from), ceilInt(to)
//	if pixb > max {
//		pixb = max
//	}
//	return
//}
//func (s *SoftwareSub) LineTo(p mgl32.Vec2) {
//	// Point setup
//	from := s.last
//	to := p
//	s.last = p
//	size := s.bound.Size()
//	//
//	// direction setup
//	var dir float32 = 1
//	if from[vecy] > to[vecy]  {
//		dir, from, to = -1, to, from
//	}
//	if to[vecy]-from[vecy] <= closeToZero {
//		return
//	}
//	// delta xCurr/ delta y
//	Δxy := (to[vecx] - from[vecx]) / (to[vecy] - from[vecy])
//	xCurr := from[vecx]
//	yFrom, yTo := fromTo(from[vecy], to[vecy], size.Y)
//	//
//	for y := yFrom; y < yTo; y++ {
//		Δy := min(float32(y+1), to[vecy]) - max(float32(y), from[vecy])
//		xNext := xCurr + Δy*Δxy // = xCurr + Δx
//		if y < 0 {
//			xCurr = xNext
//			continue
//		}
//		buf := s.bufPix[y*size.X:]
//		d := float32(Δy * dir)
//		x0, x1 := xCurr, xNext
//		if xCurr > xNext {
//			x0, x1 = x1, x0
//		}
//		x0i := floorInt(x0)
//		x0Floor := float32(x0i)
//		x1i := ceilInt(x1)
//		x1Ceil := float32(x1i)
//		if x1i <= x0i+1 {
//			xmf := float32(0.5*(xCurr+xNext)) - x0Floor
//			if i := iclamp(x0i+0, 0, size.X); i < len(buf) {
//				buf[i] += d - float32(d*xmf)
//			}
//			if i := iclamp(x0i+1, 0, size.X); i < len(buf) {
//				buf[i] += float32(d * xmf)
//			}
//		} else {
//			s := 1 / (x1 - x0)
//			x0f := x0 - x0Floor
//			oneMinusX0f := 1 - x0f
//			a0 := float32(0.5 * s * oneMinusX0f * oneMinusX0f)
//			x1f := x1 - x1Ceil + 1
//			am := float32(0.5 * s * x1f * x1f)
//
//			if i := iclamp(x0i, 0, size.X); i < len(buf) {
//				buf[i] += float32(d * a0)
//			}
//
//			if x1i == x0i+2 {
//				if i := iclamp(x0i+1, 0, size.X); i < len(buf) {
//					buf[i] += float32(d * (1 - a0 - am))
//				}
//			} else {
//				a1 := float32(s * (1.5 - x0f))
//				if i := iclamp(x0i+1, 0, size.X); i < len(buf) {
//					buf[i] += float32(d * (a1 - a0))
//				}
//				dTimesS := float32(d * s)
//				for xi := x0i + 2; xi < x1i-1; xi++ {
//					if i := iclamp(xi, 0, size.X); i < len(buf) {
//						buf[i] += dTimesS
//					}
//				}
//				a2 := a1 + float32(s*float32(x1i-x0i-3))
//				if i := iclamp(x1i-1, 0, size.X); i < len(buf) {
//					buf[i] += float32(d * (1 - a2 - am))
//				}
//			}
//
//			if i := iclamp(x1i, 0, size.X); i < len(buf) {
//				buf[i] += float32(d * am)
//			}
//		}
//		xCurr = xNext
//	}
//}
//func (s *SoftwareSub) QuadTo(pivot, to mgl32.Vec2) {
//	from := s.last
//	devsq := DevSquared(from, pivot, to)
//	if devsq >= 0.333 {
//		const tol = 3
//		n := 1 + int(math.Sqrt(math.Sqrt(tol*float64(devsq))))
//		t, nInv := float32(0), 1/float32(n)
//		for i := 0; i < n-1; i++ {
//			t += nInv
//			fromPivot := Lerp(t, from, pivot)
//			pivotTo := Lerp(t, pivot, to)
//			s.LineTo(Lerp(t, fromPivot, pivotTo))
//		}
//	}
//	s.LineTo(to)
//}
//func (s *SoftwareSub) CubeTo(pivot1, pivot2, to mgl32.Vec2) {
//	from := s.last
//	devsq := DevSquared(from, pivot1, to)
//	if devsqAlt := DevSquared(from, pivot2, to); devsq < devsqAlt {
//		devsq = devsqAlt
//	}
//	if devsq >= 0.333 {
//		const tol = 3
//		n := 1 + int(math.Sqrt(math.Sqrt(tol*float64(devsq))))
//		t, nInv := float32(0), 1/float32(n)
//		for i := 0; i < n-1; i++ {
//			t += nInv
//			ab := Lerp(t, from, pivot1)
//			bc := Lerp(t, pivot1, pivot2)
//			cd := Lerp(t, pivot2, to)
//			abc := Lerp(t, ab, bc)
//			bcd := Lerp(t, bc, cd)
//			s.LineTo(Lerp(t, abc, bcd))
//		}
//	}
//	s.LineTo(to)
//}
//func (s *SoftwareSub) CloseTo() {
//	s.LineTo(s.start)
//}
//func (s *SoftwareSub) Fill() {
//	//w, h := s.Size()
//	//for y := 0; y < int(h); y++ {
//	//	for x := 0; x < int(w); x++ {
//	//		fmt.Printf("%5.2f ", s.bufPix[x + y * int(w)])
//	//	}
//	//	fmt.Println()
//	//}
//	//fmt.Println("==================================================================================")
//	switch t := s.filler.(type) {
//	case ColorFiller:
//		s.uniformFill(float32(t.R), float32(t.G), float32(t.B), float32(t.A))
//		//case _ImageFillerFixed:
//	case *_ImageFillerFixed:
//		s.fixedFill(t.img)
//	default:
//		s.fillerFill()
//	}
//	//for y := 0; y < int(h); y++ {
//	//	for x := 0; x < int(w); x++ {
//	//		fmt.Printf("%5.2f ", s.bufPix[x + y * int(w)])
//	//	}
//	//	fmt.Println()
//	//}
//
//	return
//}
//
//const almostZero float32 = 0.000001
//
//func (s *SoftwareSub) uniformFill(cr, cg, cb, ca float32) {
//	acc := float32(0)
//	for i, v := range s.bufPix {
//		acc += v
//		pixa := acc
//		if pixa < 0 {
//			pixa = -pixa
//		}
//		if pixa > 1 {
//			pixa = 1
//		}
//
//		s.bufPix[i] = pixa
//	}
//	width := s.bound.Dx()
//	for x := s.bound.Min.X; x < s.bound.Max.X; x++ {
//		for y := s.bound.Min.Y; y < s.bound.Max.Y; y++ {
//			pixOffset := s.root.img.PixOffset(x, y)
//			bufOffset := (x - s.bound.Min.X) + width*(y-s.bound.Min.Y)
//			buf := s.bufPix[bufOffset]
//
//			if buf < almostZero {
//				continue
//			}
//
//			//s.root.img.Pix[pixOffset+pixr] = uint8(iclamp(int(buf*cr), 0, 255))
//			//s.root.img.Pix[pixOffset+pixg] = uint8(iclamp(int(buf*cg), 0, 255))
//			//s.root.img.Pix[pixOffset+pixb] = uint8(iclamp(int(buf*cb), 0, 255))
//			//s.root.img.Pix[pixOffset+pixa] = uint8(iclamp(int(buf*ca), 0, 255))
//			sr := iclamp(int(buf*float32(cr)), 0, 255) * 0x101
//			sg := iclamp(int(buf*float32(cg)), 0, 255) * 0x101
//			sb := iclamp(int(buf*float32(cb)), 0, 255) * 0x101
//			sa := iclamp(int(buf*float32(ca)), 0, 255) * 0x101
//			tempa := (math.MaxUint16 - sa) * 0x101
//			s.root.img.Pix[pixOffset+pixr] = uint8((int(s.root.img.Pix[pixOffset+pixr])*tempa/math.MaxUint16 + sr) >> 8)
//			s.root.img.Pix[pixOffset+pixg] = uint8((int(s.root.img.Pix[pixOffset+pixg])*tempa/math.MaxUint16 + sg) >> 8)
//			s.root.img.Pix[pixOffset+pixb] = uint8((int(s.root.img.Pix[pixOffset+pixb])*tempa/math.MaxUint16 + sb) >> 8)
//			s.root.img.Pix[pixOffset+pixa] = uint8((int(s.root.img.Pix[pixOffset+pixa])*tempa/math.MaxUint16 + sa) >> 8)
//		}
//	}
//
//	return
//}
//func (s *SoftwareSub) fixedFill(src *image.RGBA) {
//	acc := float32(0)
//	// Draw Rect
//	fillBound := image.Rectangle{
//		Min: image.Point{
//			X: math.MaxInt32,
//			Y: math.MaxInt32,
//		},
//		Max: image.Point{
//			X: math.MinInt32,
//			Y: math.MinInt32,
//		},
//	}
//	width := s.bound.Dx()
//	for i, v := range s.bufPix {
//		acc += v
//		pixa := acc
//		if pixa < 0 {
//			pixa = -pixa
//		}
//		if pixa > 1 {
//			pixa = 1
//		}
//		//
//		if pixa > almostZero {
//			y := i / width
//			x := i % width
//			if x < fillBound.Min.X {
//				fillBound.Min.X = x
//			}
//			if fillBound.Max.X < x {
//				fillBound.Max.X = x
//			}
//			if y < fillBound.Min.Y {
//				fillBound.Min.Y = y
//			}
//			if fillBound.Max.Y < y {
//				fillBound.Max.Y = y
//			}
//		}
//		s.bufPix[i] = pixa
//	}
//	if fillBound.Max.X < 0 {
//		return
//	}
//	//
//	for x := s.bound.Min.X; x < s.bound.Max.X; x++ {
//		for y := s.bound.Min.Y; y < s.bound.Max.Y; y++ {
//			pixOffset := s.root.img.PixOffset(x, y)
//			bufOffset := (x) + width*(y)
//			buf := s.bufPix[bufOffset]
//			if buf < almostZero {
//				continue
//			}
//			srcX, srcY := x - s.bound.Min.X -fillBound.Min.X, y-fillBound.Min.Y - s.bound.Min.Y
//			if src.Rect.Min.X <= srcX && srcX < src.Rect.Max.X && src.Rect.Min.Y <= srcY && srcY < src.Rect.Max.Y {
//				srcOffset := src.PixOffset(srcX, srcY)
//				sr := iclamp(int(buf*float32(src.Pix[srcOffset+pixr])), 0, 255) * 0x101
//				sg := iclamp(int(buf*float32(src.Pix[srcOffset+pixg])), 0, 255) * 0x101
//				sb := iclamp(int(buf*float32(src.Pix[srcOffset+pixb])), 0, 255) * 0x101
//				sa := iclamp(int(buf*float32(src.Pix[srcOffset+pixa])), 0, 255) * 0x101
//				tempa := (math.MaxUint16 - sa) * 0x101
//				s.root.img.Pix[pixOffset+pixr] = uint8((int(s.root.img.Pix[pixOffset+pixr])*tempa/math.MaxUint16 + sr) >> 8)
//				s.root.img.Pix[pixOffset+pixg] = uint8((int(s.root.img.Pix[pixOffset+pixg])*tempa/math.MaxUint16 + sg) >> 8)
//				s.root.img.Pix[pixOffset+pixb] = uint8((int(s.root.img.Pix[pixOffset+pixb])*tempa/math.MaxUint16 + sb) >> 8)
//				s.root.img.Pix[pixOffset+pixa] = uint8((int(s.root.img.Pix[pixOffset+pixa])*tempa/math.MaxUint16 + sa) >> 8)
//			}
//		}
//	}
//
//	return
//}
//func (s *SoftwareSub) fillerFill() {
//	acc := float32(0)
//	// Draw Rect
//	fillBound := image.Rectangle{
//		Min: image.Point{
//			X: math.MaxInt32,
//			Y: math.MaxInt32,
//		},
//		Max: image.Point{
//			X: math.MinInt32,
//			Y: math.MinInt32,
//		},
//	}
//	width := s.bound.Dx()
//	for i, v := range s.bufPix {
//		acc += v
//		pixa := acc
//		if pixa < 0 {
//			pixa = -pixa
//		}
//		if pixa > 1 {
//			pixa = 1
//		}
//		//
//		if pixa > almostZero {
//			y := i / width
//			x := i % width
//			if x <= fillBound.Min.X {
//				fillBound.Min.X = x
//			}
//			if fillBound.Max.X <= x {
//				fillBound.Max.X = x
//			}
//			if y <= fillBound.Min.Y {
//				fillBound.Min.Y = y
//			}
//			if fillBound.Max.Y <= y {
//				fillBound.Max.Y = y
//			}
//		}
//		s.bufPix[i] = pixa
//	}
//	if fillBound.Max.X < 0 {
//		return
//	}
//	//
//	if fwb, ok := s.filler.(FillerWithBound); ok {
//		fwb.to(fillBound)
//	}
//	//
//	for x := s.bound.Min.X; x < s.bound.Max.X; x++ {
//		for y := s.bound.Min.Y; y < s.bound.Max.Y; y++ {
//			pixOffset := s.root.img.PixOffset(x, y)
//			bufOffset := (x - s.bound.Min.X) + width*(y-s.bound.Min.Y)
//			buf := s.bufPix[bufOffset]
//			if buf <= almostZero {
//				continue
//			}
//			c := s.filler.rgba((x - s.bound.Min.X)-fillBound.Min.X , (y-s.bound.Min.Y)-fillBound.Min.Y)
//			if c.A == 0 {
//				continue
//			}
//			sr := iclamp(int(buf*float32(c.R)), 0, 255) * 0x101
//			sg := iclamp(int(buf*float32(c.G)), 0, 255) * 0x101
//			sb := iclamp(int(buf*float32(c.B)), 0, 255) * 0x101
//			sa := iclamp(int(buf*float32(c.A)), 0, 255) * 0x101
//			tempa := (math.MaxUint16 - sa) * 0x101
//			s.root.img.Pix[pixOffset+pixr] = uint8((int(s.root.img.Pix[pixOffset+pixr])*tempa/math.MaxUint16 + sr) >> 8)
//			s.root.img.Pix[pixOffset+pixg] = uint8((int(s.root.img.Pix[pixOffset+pixg])*tempa/math.MaxUint16 + sg) >> 8)
//			s.root.img.Pix[pixOffset+pixb] = uint8((int(s.root.img.Pix[pixOffset+pixb])*tempa/math.MaxUint16 + sb) >> 8)
//			s.root.img.Pix[pixOffset+pixa] = uint8((int(s.root.img.Pix[pixOffset+pixa])*tempa/math.MaxUint16 + sa) >> 8)
//		}
//	}
//
//	return
//}
