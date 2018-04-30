package gorat

import (
	"image"
	"image/color"
	"image/draw"
	"math"
)
type (
	Filler interface {
		rgba(x, y int) color.RGBA
	}
	FillerWithBound interface {
		Filler
		to(r image.Rectangle)
	}
)
type ColorFiller color.RGBA



var ColorFillerModel = color.ModelFunc(func(c color.Color) color.Color {
	r, g, b, a := c.RGBA()
	return ColorFiller{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)}
})
func NewColorFiller(r, g, b, a uint8) ColorFiller {
	return ColorFiller{r, g, b, a}
}
func (s ColorFiller) rgba(x, y int) color.RGBA {
	return color.RGBA(s)
}
func (s ColorFiller) RGBA() (r, g, b, a uint32) {
	return s.RGBA()
}

type ImageFillerMode uint8
const (
	ImageFillerFixed            ImageFillerMode = iota
	ImageFillerNearest          ImageFillerMode = iota
	ImageFillerNearestNeighbor  ImageFillerMode = iota
	ImageFillerGausian          ImageFillerMode = iota
	ImageFillerRepeat           ImageFillerMode = iota
	ImageFillerVerticalRepeat   ImageFillerMode = iota
	ImageFillerHorizontalRepeat ImageFillerMode = iota
)

func NewImageFiller(i image.Image, mode ImageFillerMode) Filler {
	temp := image.NewRGBA(i.Bounds())
	draw.Draw(temp, temp.Rect, i, i.Bounds().Min, draw.Src)
	switch mode {
	case ImageFillerFixed:
		return &_ImageFillerFixed{
			img: temp,
		}
	case ImageFillerNearest:
		return &_ImageFillerNearest{
			img: temp,
		}
	case ImageFillerNearestNeighbor:
		return &_ImageFillerNearestNeighbor{
			img: temp,
		}
	case ImageFillerGausian:
		return &_ImageFillerGaussian{
			img: temp,
		}
	case ImageFillerRepeat:
		return  &_ImageFillerRepeat{
			img:temp,
			width: temp.Rect.Dx(),
			height:temp.Rect.Dy(),
		}
	case ImageFillerVerticalRepeat:
		return  &_ImageFillerVerticalRepeat{
			img:temp,
			width:temp.Rect.Dx(),
			height:temp.Rect.Dy(),
		}
	case ImageFillerHorizontalRepeat:
		return  &_ImageFillerHorizontalRepeat{
			img:temp,
			width:temp.Rect.Dx(),
			height:temp.Rect.Dy(),
		}
	}
	panic("Unknown mode")
}

//===============================================================================================
type memorization struct {
	val        float32
	start, end int
}
func filterRange(v, length int, delta, radius float32) (val float32, start, end int) {
	val = float32(v) * delta
	if val < 0 {
		val = 0
	} else if val >= float32(length) {
		val = float32(length)
	}
	start, end = int(val-radius+0.5), int(val+radius)
	if start < 0 {
		start = 0
	}
	if end >= length {
		end = length
	}
	return
}
const support = 1.
//===============================================================================================
type _ImageFillerFixed struct {
	img *image.RGBA
}
func (s *_ImageFillerFixed) rgba(x, y int) color.RGBA {
	offset := s.img.PixOffset(x+s.img.Rect.Min.X, y+s.img.Rect.Min.X)
	if offset >= len(s.img.Pix){
		return color.RGBA{0,0,0,0}
	}
	return color.RGBA{
		R: s.img.Pix[offset+pixr],
		G: s.img.Pix[offset+pixg],
		B: s.img.Pix[offset+pixb],
		A: s.img.Pix[offset+pixa],
	}
}
//===============================================================================================
type _ImageFillerNearest struct {
	img *image.RGBA

	t      image.Rectangle
	mw     []memorization
	mh     []memorization
	scaleh float32
	scalev float32
}
func (s _ImageFillerNearest) fn(x float32) float32 {
	if x < 0 {
		x = -x
	}
	if x < 1.0 {
		return 1.0 - x
	}
	return 0
}
func (s *_ImageFillerNearest) rgba(x, y int) color.RGBA {
	var rr, rg, rb, ra float32
	var sum float32
	//
	if !(s.t.Min.X <= x && x < s.t.Max.X && s.t.Min.Y <= y && y < s.t.Max.Y) {
		return color.RGBA{0, 0, 0, 0}
	}
	h, v := s.mw[x-s.t.Min.X], s.mh[y-s.t.Min.Y]
	// H : pixel evaluate
	for kx := h.start; kx < h.end; kx++ {
		srcoffset := s.img.PixOffset(kx+s.img.Rect.Min.X, int(v.val)+s.img.Rect.Min.Y)
		normal := (float32(kx) - h.val) / s.scaleh
		res := s.fn(normal)
		// normalized pixr, pixg, pixb, pixa and sum
		rr += float32(s.img.Pix[srcoffset+pixr]) * res
		rg += float32(s.img.Pix[srcoffset+pixg]) * res
		rb += float32(s.img.Pix[srcoffset+pixb]) * res
		ra += float32(s.img.Pix[srcoffset+pixa]) * res
		sum += res
	}
	// V : pixel evaluate
	for ky := v.start; ky < v.end; ky++ {
		srcoffset := s.img.PixOffset(int(h.val)+s.img.Rect.Min.X, ky+s.img.Rect.Min.Y)
		normal := (float32(ky) - v.val) / s.scalev
		res := s.fn(normal)
		// normalized pixr, pixg, pixb, pixa and sum
		rr += float32(s.img.Pix[srcoffset+pixr]) * res
		rg += float32(s.img.Pix[srcoffset+pixg]) * res
		rb += float32(s.img.Pix[srcoffset+pixb]) * res
		ra += float32(s.img.Pix[srcoffset+pixa]) * res
		sum += res
	}
	return color.RGBA{
		R: uint8(Clamp((rr/sum)+0.5, 0, 255)),
		G: uint8(Clamp((rg/sum)+0.5, 0, 255)),
		B: uint8(Clamp((rb/sum)+0.5, 0, 255)),
		A: uint8(Clamp((ra/sum)+0.5, 0, 255)),
	}
}
func (s *_ImageFillerNearest) to(r image.Rectangle) {
	s.t = r
	s.mw = make([]memorization, s.t.Dx())
	s.mh = make([]memorization, s.t.Dy())

	// memorization
	var deltaH = float32(s.img.Rect.Dx()) / float32(r.Dx())
	s.scaleh = float32(math.Max(float64(deltaH), 1.0))
	var radiusH = float32(math.Ceil(float64(s.scaleh * support)))
	var deltaV = float32(s.img.Rect.Dy()) / float32(r.Dy())
	s.scalev = float32(math.Max(float64(deltaV), 1.0))
	var radiusV = float32(math.Ceil(float64(s.scalev * support)))
	//
	for x := 0; x < r.Dx(); x++ {
		xsrc, hstart, hend := filterRange(x, s.img.Rect.Dx(), deltaH, radiusH)
		s.mw[x] = memorization{
			val:   xsrc + float32(s.img.Rect.Min.X),
			start: hstart + s.img.Rect.Min.X,
			end:   hend + s.img.Rect.Min.X,
		}
	}
	for y := 0; y < r.Dy(); y++ {
		ysrc, vstart, vend := filterRange(y, s.img.Rect.Dx(), deltaV, radiusV)
		s.mh[y] = memorization{
			val:   ysrc + float32(s.img.Rect.Min.Y),
			start: vstart + s.img.Rect.Min.Y,
			end:   vend + s.img.Rect.Min.Y,
		}
	}
}

//===============================================================================================
type _ImageFillerNearestNeighbor struct {
	img    *image.RGBA
	t      image.Rectangle
	dx, dy float32
}
func (s *_ImageFillerNearestNeighbor) rgba(x, y int) color.RGBA {
	x = int((float32(x)+0.5)*s.dx) + s.img.Rect.Min.X
	y = int((float32(y)+0.5)*s.dy) + s.img.Rect.Min.Y
	if x >= s.img.Rect.Max.X {
		x = s.img.Rect.Max.X - 1
	}
	if y >= s.img.Rect.Max.Y {
		y = s.img.Rect.Max.Y - 1
	}
	offset := s.img.PixOffset(x, y)
	return color.RGBA{
		R: s.img.Pix[offset+pixr],
		G: s.img.Pix[offset+pixg],
		B: s.img.Pix[offset+pixb],
		A: s.img.Pix[offset+pixa],
	}
}
func (s *_ImageFillerNearestNeighbor) to(r image.Rectangle) {
	s.t = r
	s.dx = float32(s.img.Rect.Dx()) / float32(r.Dx())
	s.dy = float32(s.img.Rect.Dy()) / float32(r.Dy())
}

//===============================================================================================
type _ImageFillerGaussian struct {
	img *image.RGBA

	t      image.Rectangle
	mw     []memorization
	mh     []memorization
	scaleh float32
	scalev float32
}
func (s _ImageFillerGaussian) fn(x float32) float32 {
	var tempx = math.Abs(float64(x))
	if x < 1.0 {
		exp := 2.0
		x *= 2.0
		y := math.Pow(0.5, math.Pow(tempx, exp))
		base := math.Pow(0.5, math.Pow(2, exp))
		return float32((y - base) / (1 - base))
	}
	return 0
}
func (s *_ImageFillerGaussian) rgba(x, y int) color.RGBA {
	var rr, rg, rb, ra float32
	var sum float32
	//
	if !(s.t.Min.X <= x && x < s.t.Max.X && s.t.Min.Y <= y && y < s.t.Max.Y) {
		return color.RGBA{0, 0, 0, 0}
	}
	h, v := s.mw[x-s.t.Min.X], s.mh[y-s.t.Min.Y]
	// H : pixel evaluate
	for kx := h.start; kx < h.end; kx++ {
		srcoffset := s.img.PixOffset(kx+s.img.Rect.Min.X, int(v.val)+s.img.Rect.Min.Y)
		normal := (float32(kx) - h.val) / s.scaleh
		res := s.fn(normal)
		// normalized pixr, pixg, pixb, pixa and sum
		rr += float32(s.img.Pix[srcoffset+pixr]) * res
		rg += float32(s.img.Pix[srcoffset+pixg]) * res
		rb += float32(s.img.Pix[srcoffset+pixb]) * res
		ra += float32(s.img.Pix[srcoffset+pixa]) * res
		sum += res
	}
	// V : pixel evaluate
	for ky := v.start; ky < v.end; ky++ {
		srcoffset := s.img.PixOffset(int(h.val)+s.img.Rect.Min.X, ky+s.img.Rect.Min.Y)
		normal := (float32(ky) - v.val) / s.scalev
		res := s.fn(normal)
		// normalized pixr, pixg, pixb, pixa and sum
		rr += float32(s.img.Pix[srcoffset+pixr]) * res
		rg += float32(s.img.Pix[srcoffset+pixg]) * res
		rb += float32(s.img.Pix[srcoffset+pixb]) * res
		ra += float32(s.img.Pix[srcoffset+pixa]) * res
		sum += res
	}
	return color.RGBA{
		R: uint8(Clamp((rr/sum)+0.5, 0, 255)),
		G: uint8(Clamp((rg/sum)+0.5, 0, 255)),
		B: uint8(Clamp((rb/sum)+0.5, 0, 255)),
		A: uint8(Clamp((ra/sum)+0.5, 0, 255)),
	}
}
func (s *_ImageFillerGaussian) to(r image.Rectangle) {
	s.t = r
	s.mw = make([]memorization, s.t.Dx())
	s.mh = make([]memorization, s.t.Dy())

	// memorization
	var deltaH = float32(s.img.Rect.Dx()) / float32(r.Dx())
	s.scaleh = float32(math.Max(float64(deltaH), 1.0))
	var radiusH = float32(math.Ceil(float64(s.scaleh * support)))
	var deltaV = float32(s.img.Rect.Dy()) / float32(r.Dy())
	s.scalev = float32(math.Max(float64(deltaV), 1.0))
	var radiusV = float32(math.Ceil(float64(s.scalev * support)))
	//
	for x := 0; x < r.Dx(); x++ {
		xsrc, hstart, hend := filterRange(x, s.img.Rect.Dx(), deltaH, radiusH)
		s.mw[x] = memorization{
			val:   xsrc + float32(s.img.Rect.Min.X),
			start: hstart + s.img.Rect.Min.X,
			end:   hend + s.img.Rect.Min.X,
		}
	}
	for y := 0; y < r.Dy(); y++ {
		ysrc, vstart, vend := filterRange(y, s.img.Rect.Dx(), deltaV, radiusV)
		s.mh[y] = memorization{
			val:   ysrc + float32(s.img.Rect.Min.Y),
			start: vstart + s.img.Rect.Min.Y,
			end:   vend + s.img.Rect.Min.Y,
		}
	}
}
//===============================================================================================
type _ImageFillerRepeat struct {
	img *image.RGBA
	width, height int
}
func (s *_ImageFillerRepeat) rgba(x, y int) color.RGBA {
	x = x % s.width
	y = y % s.height
	offset := s.img.PixOffset(x+s.img.Rect.Min.X, y+s.img.Rect.Min.X)
	return color.RGBA{
		R: s.img.Pix[offset+pixr],
		G: s.img.Pix[offset+pixg],
		B: s.img.Pix[offset+pixb],
		A: s.img.Pix[offset+pixa],
	}
}
//===============================================================================================
type _ImageFillerVerticalRepeat struct {
	img *image.RGBA
	width, height int
}
func (s *_ImageFillerVerticalRepeat) rgba(x, y int) color.RGBA {
	if x / s.width > 0{
		return color.RGBA{0,0,0,0}
	}
	y = y % s.height
	offset := s.img.PixOffset(x+s.img.Rect.Min.X, y+s.img.Rect.Min.X)
	if offset >= len(s.img.Pix){
		return color.RGBA{0,0,0,0}
	}
	return color.RGBA{
		R: s.img.Pix[offset+pixr],
		G: s.img.Pix[offset+pixg],
		B: s.img.Pix[offset+pixb],
		A: s.img.Pix[offset+pixa],
	}
}
//===============================================================================================
type _ImageFillerHorizontalRepeat struct {
	img *image.RGBA
	width, height int
}
func (s *_ImageFillerHorizontalRepeat) rgba(x, y int) color.RGBA {
	x = x % s.width
	if y / s.height > 0{
		return color.RGBA{0,0,0,0}
	}
	offset := s.img.PixOffset(x+s.img.Rect.Min.X, y+s.img.Rect.Min.X)
	return color.RGBA{
		R: s.img.Pix[offset+pixr],
		G: s.img.Pix[offset+pixg],
		B: s.img.Pix[offset+pixb],
		A: s.img.Pix[offset+pixa],
	}
}