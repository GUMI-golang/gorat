package gorat
//
////
//import (
//	"github.com/go-gl/mathgl/mgl32"
//	"image"
//	"image/color"
//	"image/draw"
//	"math"
//)
//
//// 성능이 걱정된다....
//// 예측대로 성능이 바닥...
//type SoftwareRasterizerPrototype struct {
//	//vertice []mgl32.VecN
//	//indexes [][3]int
//	//
//	works [][][6]float32
//	//
//	currentWorks [][6]float32
//	//
//	Options
//	buf *image.RGBA
//}
//
//func NewSoftwareRasterizer(w, h int) *SoftwareRasterizerPrototype {
//	return NewSoftwareRasterizerRGBA(image.NewRGBA(image.Rect(0, 0, w, h)))
//}
//func NewSoftwareRasterizerRGBA(rgba *image.RGBA) *SoftwareRasterizerPrototype {
//	temp := &SoftwareRasterizerPrototype{
//		buf:     rgba,
//	}
//	temp.DefaultOption()
//	return temp
//}
//func (s *SoftwareRasterizerPrototype) last() (last mgl32.Vec2) {
//	copy(last[:], s.currentWorks[len(s.currentWorks) - 1][:])
//	return
//}
//
//func (s *SoftwareRasterizerPrototype) MoveTo(to mgl32.Vec2, info *mgl32.Vec4) {
//	const capacity = 32
//	if s.currentWorks != nil {
//		s.Close()
//	}
//	s.currentWorks = make([][6]float32, 1, capacity)
//	s.currentWorks[0] = [6]float32{to[0], to[1], info[0], info[1], info[2], info[3]}
//}
//func (s *SoftwareRasterizerPrototype) LineTo(to mgl32.Vec2, info *mgl32.Vec4) {
//	if len(s.currentWorks) == 0 {
//		s.MoveTo(to, info)
//		return
//	}
//	if info == nil {
//		temp := s.Options.color
//		info = &temp
//	}
//	s.currentWorks = append(s.currentWorks, [6]float32{to[0], to[1], info[0], info[1], info[2], info[3]})
//}
//func (s *SoftwareRasterizerPrototype) QuadTo(pivot, to mgl32.Vec2, info *mgl32.Vec4) {
//	if info == nil {
//		temp := s.Options.color
//		info = &temp
//	}
//	// Come from golang.org/vecx/image/vector Raster.QuadTo
//	from := s.last()
//	devsq := devSquared(from, pivot, to)
//	if devsq >= 0.333 {
//		const tol = 3
//		n := 1 + int(math.Sqrt(math.Sqrt(tol*float64(devsq))))
//		t, nInv := float32(0), 1/float32(n)
//		for i := 0; i < n-1; i++ {
//			t += nInv
//			ab := lerp(t, from, pivot)
//			bc := lerp(t, pivot, to)
//			s.LineTo(lerp(t, ab, bc), info)
//		}
//	}
//	s.LineTo(to, info)
//}
//func (s *SoftwareRasterizerPrototype) CubeTo(pivot1, pivot2, to mgl32.Vec2, info *mgl32.Vec4) {
//	if info == nil {
//		temp := s.Options.color
//		info = &temp
//	}
//	// Come from golang.org/vecx/image/vector Raster.QuadTo
//	from := s.last()
//	devsq := devSquared(from, pivot1, to)
//	if devsqAlt := devSquared(from, pivot2, to); devsq < devsqAlt {
//		devsq = devsqAlt
//	}
//	if devsq >= 0.333 {
//		const tol = 3
//		n := 1 + int(math.Sqrt(math.Sqrt(tol*float64(devsq))))
//		t, nInv := float32(0), 1/float32(n)
//		for i := 0; i < n-1; i++ {
//			t += nInv
//			ab := lerp(t, from, pivot1)
//			bc := lerp(t, pivot1, pivot2)
//			cd := lerp(t, pivot2, to)
//			abc := lerp(t, ab, bc)
//			bcd := lerp(t, bc, cd)
//			s.LineTo(lerp(t, abc, bcd), info)
//		}
//	}
//	s.LineTo(to, info)
//}
//func (s *SoftwareRasterizerPrototype) Close() {
//	s.works = append(s.works, s.currentWorks)
//	s.currentWorks = nil
//}
//
//func (s *SoftwareRasterizerPrototype) Stroke() {
//	// If not closed, close first
//	if s.currentWorks != nil {
//		s.Close()
//	}
//	backup := s.works
//	s.works = nil
//	// Half width
//	hw := s.GetStrokeWidth() / 2
//	j, c := s.GetStrokeJoin(), s.GetStrokeCap()
//	for _, work := range backup {
//		length := len(work)
//		if length < 2 {
//			continue
//		}
//		var pO, cO, pV, cV mgl32.Vec2
//		var pC, cC mgl32.Vec4
//		pV = Vec2(work[0][0], work[0][1])
//		pC = Vec4(work[0][2], work[0][3], work[0][4], work[0][5])
//		for i := 1; i < length; i++ {
//			cV = Vec2(work[i][0], work[i][1])
//			cC = Vec4(work[i][2], work[i][3], work[i][4], work[i][5])
//			cO = ortho(pV, cV).Mul(hw)
//			// line rect
//			s.MoveTo(pV.Add(cO), &pC)
//			s.LineTo(cV.Add(cO), &cC)
//			s.LineTo(cV.Sub(cO), &cC)
//			s.LineTo(pV.Sub(cO), &pC)
//			//s.Close()
//			// line rect close
//			switch i {
//			case 1:
//				//start cap
//				switch c {
//
//				}
//			case length - 1:
//				//endcap
//				switch c {
//
//				}
//				fallthrough
//			default:
//				// join
//				var a, b, c mgl32.Vec2
//				a = pV
//				if Cross(pO, cO) > 0{
//					b= pV.Sub(pO)
//					c= pV.Sub(cO)
//				}else {
//					b= pV.Add(pO)
//					c= pV.Add(cO)
//				}
//				switch j {
//				case StrokeJoinRound:
//					bc := c.Sub(b)
//					ppV := Vec2(work[i-2][0], work[i-2][0])
//					v1 := pV.Sub(ppV)
//					v2 := pV.Sub(cV)
//					d := b.Add(v1.Mul(bc[0] / v1.Sub(v2)[0]))
//					s.MoveTo(a, &pC)
//					s.LineTo(b, &pC)
//					s.QuadTo(d, c, &cC)
//				case StrokeJoinMiter:
//					bc := c.Sub(b)
//					ppV := Vec2(work[i-2][0], work[i-2][0])
//					v1 := pV.Sub(ppV)
//					v2 := pV.Sub(cV)
//					d := b.Add(v1.Mul(bc[0] / v1.Sub(v2)[0]))
//					s.MoveTo(a, &pC)
//					s.LineTo(b, &pC)
//					s.LineTo(d, &cC)
//					s.LineTo(c, &cC)
//				case StrokeJoinBevel:
//					s.MoveTo(a, &pC)
//					s.LineTo(b, &pC)
//					s.LineTo(c, &pC)
//				}
//			}
//			pO, pV, pC = cO, cV, cC
//		}
//	}
//	b := s.Backup()
//	defer s.Restore(b)
//	// lowest calc
//	s.SetOverlap(OverlapPlus)
//	s.FillColor()
//}
//func ortho(a, b mgl32.Vec2) mgl32.Vec2 {
//	return mgl32.Vec2{-(b[1] - a[1]), b[0] - a[0]}.Normalize()
//}
//func reortho(a mgl32.Vec2) mgl32.Vec2 {
//	return mgl32.Vec2{a[1], -a[0]}.Normalize()
//}
//
//type tribufElemColor struct {
//	Triangle
//	Colors [3]mgl32.Vec4
//}
//
//func (s *SoftwareRasterizerPrototype) FillColor() {
//	// If not closed, close first
//	if s.currentWorks != nil {
//		s.Close()
//	}
//	// Triangulate
//	var tribufCW []tribufElemColor // CW Triangle save here
//	var tribufCCW []tribufElemColor
//	for _, work := range s.works {
//		// filling need points at least 3
//		if len(work) < 3 {
//			continue
//		}
//		var length = len(work)
//		var start, before = mgl32.Vec2{work[0][0], work[0][1]}, mgl32.Vec2{work[1][0], work[1][1]}
//		var curr mgl32.Vec2
//		var startC, beforeC = mgl32.Vec4{work[0][2], work[0][3], work[0][4], work[0][5]}, mgl32.Vec4{work[1][2], work[1][3], work[1][4], work[1][5]}
//		var currC mgl32.Vec4
//		for i := 2; i < length; i++ {
//			curr = mgl32.Vec2{work[i][0], work[i][1]}
//			currC = mgl32.Vec4{work[i][2], work[i][3], work[i][4], work[i][5]}
//			temp := tribufElemColor{
//				Triangle: Triangle{
//					A: start,
//					B: before,
//					C: curr,
//				},
//				Colors: [3]mgl32.Vec4{
//					startC,
//					beforeC,
//					currC,
//				},
//			}
//			if temp.RotateDirection() == CW {
//				// CW
//				tribufCW = append(tribufCW, temp)
//			} else {
//				// CCW
//				tribufCCW = append(tribufCCW, temp)
//			}
//			before, beforeC = curr, currC
//		}
//	}
//	// copy result to workspace
//	var workspace = image.NewRGBA(s.buf.Rect)
//	// CW triangle render first
//	for _, tri := range tribufCW {
//		trColorCW(workspace, tri.Triangle, tri.Colors, &s.Options)
//	}
//	// CCW triangle render first
//	for _, tri := range tribufCCW {
//		trColorCCW(workspace, tri.Triangle, tri.Colors, &s.Options)
//	}
//	// Commit change
//	draw.Draw(s.buf, s.buf.Rect, workspace, workspace.Rect.Min, draw.Over)
//}
//func (s *SoftwareRasterizerPrototype) Clear() {
//	draw.Draw(s.buf, s.buf.Rect, image.NewUniform(color.Transparent), image.ZP, draw.Src)
//}
//
//func trColorCW(buf *image.RGBA, triangle Triangle, colors [3]mgl32.Vec4, o *Options) {
//	const (
//		r, g, b, a = 0, 1, 2, 3
//	)
//	aabb := triangle.AABB()
//	for y := floorInt(aabb.Min[1]); y < ceilInt(aabb.Max[1]); y++ {
//		for x := floorInt(aabb.Min[0]); x < ceilInt(aabb.Max[0]); x++ {
//			inside, w, u, v := triangle.AAUVW(mgl32.Vec2{
//				float32(x), float32(y),
//			}, o.aa)
//			if inside{
//				if !(buf.Rect.Min.X <= x && x < buf.Rect.Max.X && buf.Rect.Min.Y <= y && y < buf.Rect.Max.Y){
//					continue
//				}
//				off := buf.PixOffset(x+buf.Rect.Min.X, y + +buf.Rect.Min.Y)
//				uc := colors[0].Mul(u)
//				vc := colors[1].Mul(v)
//				wc := colors[2].Mul(w)
//				buf.Pix[off+r] = colorAddClamp(buf.Pix[off+r], uint8((wc[r] + uc[r] + vc[r]) * 255))
//				buf.Pix[off+g] = colorAddClamp(buf.Pix[off+g], uint8((wc[g] + uc[g] + vc[g]) * 255))
//				buf.Pix[off+b] = colorAddClamp(buf.Pix[off+b], uint8((wc[b] + uc[b] + vc[b]) * 255))
//				buf.Pix[off+a] = colorAddClamp(buf.Pix[off+a], uint8((wc[a] + uc[a] + vc[a]) * 255))
//			}
//		}
//	}
//}
//func trColorCCW(buf *image.RGBA, triangle Triangle, colors [3]mgl32.Vec4, o *Options) {
//	const (
//		r, g, b, a = 0, 1, 2, 3
//	)
//	aabb := triangle.AABB()
//	for y := floorInt(aabb.Min[1]); y < ceilInt(aabb.Max[1]); y++ {
//		for x := floorInt(aabb.Min[0]); x < ceilInt(aabb.Max[0]); x++ {
//			inside, w, u, v := triangle.AAUVW(mgl32.Vec2{
//				float32(x), float32(y),
//			}, o.aa)
//			if inside{
//				if !(buf.Rect.Min.X <= x && x < buf.Rect.Max.X && buf.Rect.Min.Y <= y && y < buf.Rect.Max.Y){
//					continue
//				}
//				off := buf.PixOffset(x+buf.Rect.Min.X, y + +buf.Rect.Min.Y)
//
//				if buf.Pix[off+a] == 0 {
//					// If CW work not rastered this Pixel
//					// raster here
//					uc := colors[0].Mul(u)
//					vc := colors[1].Mul(v)
//					wc := colors[2].Mul(w)
//					buf.Pix[off+r] = uint8((wc[r] + uc[r] + vc[r]) * 255)
//					buf.Pix[off+g] = uint8((wc[g] + uc[g] + vc[g]) * 255)
//					buf.Pix[off+b] = uint8((wc[b] + uc[b] + vc[b]) * 255)
//					buf.Pix[off+a] = uint8((wc[a] + uc[a] + vc[a]) * 255)
//				} else {
//					// If CW work rastered this Pixel
//					switch o.overlap {
//					case OverlapPlus:
//						uc := colors[0].Mul(u)
//						vc := colors[1].Mul(v)
//						wc := colors[2].Mul(w)
//						buf.Pix[off+r] = colorAddClamp(buf.Pix[off+r], uint8((wc[r]+uc[r]+vc[r])*255))
//						buf.Pix[off+g] = colorAddClamp(buf.Pix[off+g], uint8((wc[g]+uc[g]+vc[g])*255))
//						buf.Pix[off+b] = colorAddClamp(buf.Pix[off+b], uint8((wc[b]+uc[b]+vc[b])*255))
//						buf.Pix[off+a] = colorAddClamp(buf.Pix[off+a], uint8((wc[a]+uc[a]+vc[a])*255))
//					case OverlapCWMinusCCW:
//						uc := colors[0].Mul(u)
//						vc := colors[1].Mul(v)
//						wc := colors[2].Mul(w)
//						buf.Pix[off+r] = colorSubClamp(buf.Pix[off+r], uint8((wc[r]+uc[r]+vc[r])*255))
//						buf.Pix[off+g] = colorSubClamp(buf.Pix[off+g], uint8((wc[g]+uc[g]+vc[g])*255))
//						buf.Pix[off+b] = colorSubClamp(buf.Pix[off+b], uint8((wc[b]+uc[b]+vc[b])*255))
//						buf.Pix[off+a] = colorSubClamp(buf.Pix[off+a], uint8((wc[a]+uc[a]+vc[a])*255))
//					case OverlapCCWMinusCW:
//						uc := colors[0].Mul(u)
//						vc := colors[1].Mul(v)
//						wc := colors[2].Mul(w)
//						buf.Pix[off+r] = colorSubClamp(uint8((wc[r]+uc[r]+vc[r])*255), buf.Pix[off+r])
//						buf.Pix[off+g] = colorSubClamp(uint8((wc[g]+uc[g]+vc[g])*255), buf.Pix[off+g])
//						buf.Pix[off+b] = colorSubClamp(uint8((wc[b]+uc[b]+vc[b])*255), buf.Pix[off+b])
//						buf.Pix[off+a] = colorSubClamp(uint8((wc[a]+uc[a]+vc[a])*255), buf.Pix[off+a])
//					case OverlapVacate:
//						buf.Pix[off+r] = 0
//						buf.Pix[off+g] = 0
//						buf.Pix[off+b] = 0
//						buf.Pix[off+a] = 0
//					case OverlapCCW:
//						uc := colors[0].Mul(u)
//						vc := colors[1].Mul(v)
//						wc := colors[2].Mul(w)
//						buf.Pix[off+r] = uint8((wc[r] + uc[r] + vc[r]) * 255)
//						buf.Pix[off+g] = uint8((wc[g] + uc[g] + vc[g]) * 255)
//						buf.Pix[off+b] = uint8((wc[b] + uc[b] + vc[b]) * 255)
//						buf.Pix[off+a] = uint8((wc[a] + uc[a] + vc[a]) * 255)
//					case OverlapCW:
//					case OverlapAverage:
//
//						uc := colors[0].Mul(u)
//						vc := colors[1].Mul(v)
//						wc := colors[2].Mul(w)
//						buf.Pix[off+r] = colorAverage(uint8((wc[r]+uc[r]+vc[r])*255), buf.Pix[off+r])
//						buf.Pix[off+g] = colorAverage(uint8((wc[g]+uc[g]+vc[g])*255), buf.Pix[off+g])
//						buf.Pix[off+b] = colorAverage(uint8((wc[b]+uc[b]+vc[b])*255), buf.Pix[off+b])
//						buf.Pix[off+a] = colorAverage(uint8((wc[a]+uc[a]+vc[a])*255), buf.Pix[off+a])
//					}
//				}
//			}
//		}
//	}
//}
//func colorAddClamp(a, b uint8) uint8 {
//	temp := int32(a) + int32(b)
//	if temp > 255 {
//		return 255
//	}
//	return uint8(temp)
//}
//func colorSubClamp(a, b uint8) uint8 {
//	temp := int16(a) - int16(b)
//	if temp < 0 {
//		return 0
//	}
//	return uint8(temp)
//}
//func colorAverage(a, b uint8) uint8 {
//	return uint8((uint16(a) + uint16(b)) / 2)
//}
//func (s *SoftwareRasterizerPrototype) TrUV(triangle Triangle, uvs [3]mgl32.Vec2, img image.Image) {
//	const (
//		R, G, B, A = 0, 1, 2, 3
//		U, V       = 0, 1
//	)
//	sz := img.Bounds().Size()
//	aabb := triangle.AABB()
//	for y := floorInt(aabb.Min[1]); y < ceilInt(aabb.Max[1]); y++ {
//		for x := floorInt(aabb.Min[0]); x < ceilInt(aabb.Max[0]); x++ {
//			inside, w, u, v := triangle.UVW(mgl32.Vec2{
//				float32(x), float32(y),
//			})
//			if inside {
//				off := s.buf.PixOffset(x, y)
//				auv := uvs[0].Mul(u)
//				buv := uvs[1].Mul(v)
//				cuv := uvs[2].Mul(w)
//				coord_u := auv[U] + buv[U] + cuv[U]
//				coord_v := auv[V] + buv[V] + cuv[V]
//				pix := color.RGBAModel.Convert(img.At(int(float32(sz.X)*coord_u), int(float32(sz.Y)*coord_v))).(color.RGBA)
//				s.buf.Pix[off+R] = pix.R
//				s.buf.Pix[off+G] = pix.G
//				s.buf.Pix[off+B] = pix.B
//				s.buf.Pix[off+A] = pix.A
//			}
//		}
//	}
//}
//func (s *SoftwareRasterizerPrototype) Image() image.Image {
//	return s.buf
//}
