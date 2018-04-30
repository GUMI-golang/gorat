package gorat

import (
	"github.com/go-gl/mathgl/mgl32"
	"image"
	"image/color"
	"math"
	"runtime"
	"image/draw"
)

const (
	HARDWAREPOINTMINCAPACITY = 64
)

type Hardware struct {
	rs HardwareResult
	SubHardware
}
type SubHardware struct {
	root *Hardware
	//
	ws HardwareWorkspace
	//
	bound   mgl32.Vec4
	working []mgl32.Vec2
	Options
}

func NewHardware(to HardwareResult) *Hardware {
	checkDriverOrPanic()
	w, h := to.Size()
	res := new(Hardware)
	res.DefaultOption()
	res.root = res
	res.bound = mgl32.Vec4{0, 0, float32(w), float32(h)}
	res.ws = hwDriver.WorkSpace(w, h)
	res.rs = to
	return res
}
func _SubHardwareCloser(h *SubHardware) {
	h.ws.Delete()
}
func (s *Hardware) Setup(w, h int) {
	s.rs.Resize(w, h)
	s.ws.Resize(w, h)
	s.working = make([]mgl32.Vec2, 0, HARDWAREPOINTMINCAPACITY)
	//s.rs = hwDriver.Result(w, h)
}

//
func (s *SubHardware) Root() Rasterizer {
	return s.root
}
func (s *SubHardware) Bound() image.Rectangle {
	return image.Rectangle{
		Min: image.Point{
			X: int(s.bound[0]),
			Y: int(s.bound[1]),
		},
		Max: image.Point{
			X: int(s.bound[2]),
			Y: int(s.bound[3]),
		},
	}
}
func (s *SubHardware) SubRasterizer(r image.Rectangle) SubRasterizer {
	temp := &SubHardware{
		root:  s.root,
		bound: Vec4(float32(r.Min.X), float32(r.Min.Y), float32(r.Max.X), float32(r.Max.Y)),
	}
	temp.DefaultOption()
	temp.ws = hwDriver.WorkSpace(r.Dx(), r.Dy())
	runtime.SetFinalizer(temp, _SubHardwareCloser)
	return temp
}

// vector drawer
func (s *SubHardware) Size() (w, h float32) {
	return s.bound[2] - s.bound[0], s.bound[3] - s.bound[1]
}
func (s *SubHardware) Reset() {
	s.ws.Clear()
	s.working = make([]mgl32.Vec2, 0, HARDWAREPOINTMINCAPACITY)
}
func (s *SubHardware) MoveTo(to mgl32.Vec2) {
	if len(s.working) == 0 {
		s.working = []mgl32.Vec2{
			{to[vecx] + s.bound[vecx], to[vecy] + s.bound[vecy]},
		}
	} else {
		s.CloseTo()
		s.working = append(s.working, mgl32.Vec2{to[vecx] + s.bound[vecx], to[vecy] + s.bound[vecy]})
	}
}
func (s *SubHardware) LineTo(to mgl32.Vec2) {

	s.working = append(s.working, mgl32.Vec2{to[vecx] + s.bound[vecx], to[vecy] + s.bound[vecy]})
}
func (s *SubHardware) QuadTo(pivot, to mgl32.Vec2) {
	from := s.Point()
	pivot = mgl32.Vec2{pivot[vecx] + s.bound[vecx], pivot[vecy] + s.bound[vecy]}
	to = mgl32.Vec2{to[vecx] + s.bound[vecx], to[vecy] + s.bound[vecy]}
	for _, to := range quadFromTo(from, pivot, to) {
		s.LineTo(to)
	}
}
func (s *SubHardware) CubeTo(pivot1, pivot2, to mgl32.Vec2) {
	from := s.Point()
	pivot1 = mgl32.Vec2{pivot1[vecx] + s.bound[vecx], pivot1[vecy] + s.bound[vecy]}
	pivot2 = mgl32.Vec2{pivot2[vecx] + s.bound[vecx], pivot2[vecy] + s.bound[vecy]}
	to = mgl32.Vec2{to[vecx] + s.bound[vecx], to[vecy] + s.bound[vecy]}
	for _, to := range cubeFromTo(from, pivot1, pivot2, to) {
		s.LineTo(to)
	}
}
func (s *SubHardware) CloseTo() {
	if math.IsNaN(float64(s.Point()[vecx])) {
		return
	}
	s.working = append(s.working, s.PreviousPoint())
	s.working = append(s.working, nanvec2)
}
func (s *SubHardware) PreviousPoint() mgl32.Vec2 {
	var i int
	for i = len(s.working) - 1; i >= 0; i-- {
		if math.IsNaN(float64(s.working[i][0])) {
			if i == len(s.working)-1 {
				i = -2
			}
			break
		}
	}
	if i < -1 {
		return nanvec2
	}
	return s.working[i+1]

}
func (s *SubHardware) Point() mgl32.Vec2 {
	if len(s.working) == 0 {
		return nanvec2
	}
	return s.working[len(s.working)-1]
}
func (s *SubHardware) Clear() {
	s.root.rs.RectClear(s.Bound())
}
func (s *SubHardware) Stroke() {
	if len(s.working) < 2{
		panic("Too little point")
	}
	s.CloseTo()
	defer s.Reset()
	splited := splitStroke(s.working)
	for _, l := range splited {
		for _, v := range stroke(l, s.Options) {
			s.LineTo(v)
		}
		s.CloseTo()
	}
	s.fill(s.color)
}
func (s *SubHardware) Fill() {
	if len(s.working) < 3{
		panic("Too little point")
	}
	s.CloseTo()
	defer s.Reset()
	s.fill(s.filler)
}
func (s *SubHardware) FillStroke() {
	if len(s.working) < 3{
		panic("Too little point")
	}
	panic("implement me")
}
func (s *SubHardware) StrokeFill() {
	if len(s.working) < 3{
		panic("Too little point")
	}
	panic("implement me")
}

func (s *SubHardware) debugStroke(stroking, filling draw.Image) {
	if len(s.working) < 2{
		panic("Too little point")
	}
	s.CloseTo()
	defer s.Reset()
	splited := splitStroke(s.working)
	for _, l := range splited {
		for _, v := range stroke(l, s.Options) {
			s.LineTo(v)
		}
		s.CloseTo()
	}
	s.fillDebug(s.color, stroking, filling)
}
func (s *SubHardware) debugFill(stroking, filling draw.Image) {
	if len(s.working) < 3{
		panic("Too little point")
	}
	s.CloseTo()
	defer s.Reset()
	s.fillDebug(s.filler, stroking, filling)
}

func (s *SubHardware) fill(defineFiller Filler) {
	var prog HardwareProgram
	//=======================================================
	// stroking
	prog = <-hwpStroker
	prog.Use()
	// points context
	ctx := hwDriver.Context()
	defer ctx.Delete()
	// bound
	bd := hwDriver.Bound()
	defer bd.Delete()
	// set points data
	ctx.Set(s.working...)
	bd.Set(image.Rect(int(s.bound[0]), int(s.bound[1]), int(s.bound[2]), int(s.bound[3])))
	//
	prog.BindWorkspace(0, s.ws)
	prog.BindContext(1, ctx)
	prog.BindBound(2, bd)
	prog.Compute(len(s.working)-1, 1, 1)
	hwpStroker <- prog
	//=======================================================
	// filling
	prog = <-hwpFiller
	prog.Use()

	bd.Set(image.Rectangle{
		Min: image.Point{
			X: math.MaxInt32,
			Y: math.MaxInt32,
		},
		Max: image.Point{
			X: math.MinInt32,
			Y: math.MinInt32,
		},
	})
	prog.BindWorkspace(0, s.ws)
	prog.BindBound(1, bd)
	prog.Compute(int(s.bound[3]-s.bound[1]), 1, 1)
	hwpFiller <- prog
	//=======================================================
	// fillstyling
	switch f := defineFiller.(type) {
	case ColorFiller:
		c := hwDriver.Color()
		defer c.Delete()
		c.Set(color.RGBA(f))
		prog = <-hwpColor
		prog.Use()
		prog.BindWorkspace(0, s.ws)
		prog.BindResult(1, s.root.rs)
		prog.BindColor(2, c)

		prog.Compute(int(s.bound[2]-s.bound[0]), int(s.bound[3]-s.bound[1]), 1)
		hwpColor <- prog
	case *_ImageFillerFixed:
		filler := hwDriver.Filler(f.img)
		defer filler.Delete()
		prog = <-hwpFixed
		prog.Use()
		prog.BindWorkspace(0, s.ws)
		prog.BindResult(1, s.root.rs)
		prog.BindBound(2, bd)
		prog.BindFiller(3, filler)
		prog.Compute(int(s.bound[2]-s.bound[0]), int(s.bound[3]-s.bound[1]), 1)
		hwpFixed <- prog
	case *_ImageFillerGaussian:
		filler := hwDriver.Filler(f.img)
		defer filler.Delete()
		prog = <-hwpGaussian
		prog.Use()
		prog.BindWorkspace(0, s.ws)
		prog.BindResult(1, s.root.rs)
		prog.BindBound(2, bd)
		prog.BindFiller(3, filler)
		prog.Compute(int(s.bound[2]-s.bound[0]), int(s.bound[3]-s.bound[1]), 1)
		hwpGaussian <- prog
	case *_ImageFillerNearest:
		filler := hwDriver.Filler(f.img)
		defer filler.Delete()
		prog = <-hwpNearest
		prog.Use()
		prog.BindWorkspace(0, s.ws)
		prog.BindResult(1, s.root.rs)
		prog.BindBound(2, bd)
		prog.BindFiller(3, filler)
		prog.Compute(int(s.bound[2]-s.bound[0]), int(s.bound[3]-s.bound[1]), 1)
		hwpNearest <- prog
	case *_ImageFillerNearestNeighbor:
		filler := hwDriver.Filler(f.img)
		defer filler.Delete()
		prog = <-hwpNearestNeighbor
		prog.Use()
		prog.BindWorkspace(0, s.ws)
		prog.BindResult(1, s.root.rs)
		prog.BindBound(2, bd)
		prog.BindFiller(3, filler)
		prog.Compute(int(s.bound[2]-s.bound[0]), int(s.bound[3]-s.bound[1]), 1)
		hwpNearestNeighbor <- prog
	case *_ImageFillerRepeat:
		filler := hwDriver.Filler(f.img)
		defer filler.Delete()
		prog = <-hwpRepeat
		prog.Use()
		prog.BindWorkspace(0, s.ws)
		prog.BindResult(1, s.root.rs)
		prog.BindBound(2, bd)
		prog.BindFiller(3, filler)
		prog.Compute(int(s.bound[2]-s.bound[0]), int(s.bound[3]-s.bound[1]), 1)
		hwpRepeat <- prog
	case *_ImageFillerHorizontalRepeat:
		filler := hwDriver.Filler(f.img)
		defer filler.Delete()
		prog = <-hwpRepeat_h
		prog.Use()
		prog.BindWorkspace(0, s.ws)
		prog.BindResult(1, s.root.rs)
		prog.BindBound(2, bd)
		prog.BindFiller(3, filler)
		prog.Compute(int(s.bound[2]-s.bound[0]), int(s.bound[3]-s.bound[1]), 1)
		hwpRepeat_h <- prog
	case *_ImageFillerVerticalRepeat:
		filler := hwDriver.Filler(f.img)
		defer filler.Delete()
		prog = <-hwpRepeat_v
		prog.Use()
		prog.BindWorkspace(0, s.ws)
		prog.BindResult(1, s.root.rs)
		prog.BindBound(2, bd)
		prog.BindFiller(3, filler)
		prog.Compute(int(s.bound[2]-s.bound[0]), int(s.bound[3]-s.bound[1]), 1)
		hwpRepeat_v <- prog
	}
}
func (s *SubHardware) fillDebug(defineFiller Filler, stroking, filling draw.Image) {
	var prog HardwareProgram
	//=======================================================
	// stroking
	prog = <-hwpStroker
	prog.Use()
	// points context
	ctx := hwDriver.Context()
	defer ctx.Delete()
	// bound
	bd := hwDriver.Bound()
	defer bd.Delete()
	// set points data
	ctx.Set(s.working...)
	bd.Set(image.Rect(int(s.bound[0]), int(s.bound[1]), int(s.bound[2]), int(s.bound[3])))
	//
	prog.BindWorkspace(0, s.ws)
	prog.BindContext(1, ctx)
	prog.BindBound(2, bd)
	prog.Compute(len(s.working)-1, 1, 1)
	hwpStroker <- prog
	draw.Draw(stroking, stroking.Bounds(), s.ws.Visualize(), image.ZP, draw.Src)
	//=======================================================
	// filling
	prog = <-hwpFiller
	prog.Use()

	bd.Set(image.Rectangle{
		Min: image.Point{
			X: math.MaxInt32,
			Y: math.MaxInt32,
		},
		Max: image.Point{
			X: math.MinInt32,
			Y: math.MinInt32,
		},
	})
	prog.BindWorkspace(0, s.ws)
	prog.BindBound(1, bd)
	prog.Compute(int(s.bound[3]-s.bound[1]), 1, 1)
	draw.Draw(filling, filling.Bounds(), s.ws.Visualize(), image.ZP, draw.Src)
	hwpFiller <- prog
	//=======================================================
	// fillstyling
	switch f := defineFiller.(type) {
	case ColorFiller:
		c := hwDriver.Color()
		defer c.Delete()
		c.Set(color.RGBA(f))
		prog = <-hwpColor
		prog.Use()
		prog.BindWorkspace(0, s.ws)
		prog.BindResult(1, s.root.rs)
		prog.BindColor(2, c)

		prog.Compute(int(s.bound[2]-s.bound[0]), int(s.bound[3]-s.bound[1]), 1)
		hwpColor <- prog
	case *_ImageFillerFixed:
		filler := hwDriver.Filler(f.img)
		defer filler.Delete()
		prog = <-hwpFixed
		prog.Use()
		prog.BindWorkspace(0, s.ws)
		prog.BindResult(1, s.root.rs)
		prog.BindBound(2, bd)
		prog.BindFiller(3, filler)
		prog.Compute(int(s.bound[2]-s.bound[0]), int(s.bound[3]-s.bound[1]), 1)
		hwpFixed <- prog
	case *_ImageFillerGaussian:
		filler := hwDriver.Filler(f.img)
		defer filler.Delete()
		prog = <-hwpGaussian
		prog.Use()
		prog.BindWorkspace(0, s.ws)
		prog.BindResult(1, s.root.rs)
		prog.BindBound(2, bd)
		prog.BindFiller(3, filler)
		prog.Compute(int(s.bound[2]-s.bound[0]), int(s.bound[3]-s.bound[1]), 1)
		hwpGaussian <- prog
	case *_ImageFillerNearest:
		filler := hwDriver.Filler(f.img)
		defer filler.Delete()
		prog = <-hwpNearest
		prog.Use()
		prog.BindWorkspace(0, s.ws)
		prog.BindResult(1, s.root.rs)
		prog.BindBound(2, bd)
		prog.BindFiller(3, filler)
		prog.Compute(int(s.bound[2]-s.bound[0]), int(s.bound[3]-s.bound[1]), 1)
		hwpNearest <- prog
	case *_ImageFillerNearestNeighbor:
		filler := hwDriver.Filler(f.img)
		defer filler.Delete()
		prog = <-hwpNearestNeighbor
		prog.Use()
		prog.BindWorkspace(0, s.ws)
		prog.BindResult(1, s.root.rs)
		prog.BindBound(2, bd)
		prog.BindFiller(3, filler)
		prog.Compute(int(s.bound[2]-s.bound[0]), int(s.bound[3]-s.bound[1]), 1)
		hwpNearestNeighbor <- prog
	case *_ImageFillerRepeat:
		filler := hwDriver.Filler(f.img)
		defer filler.Delete()
		prog = <-hwpRepeat
		prog.Use()
		prog.BindWorkspace(0, s.ws)
		prog.BindResult(1, s.root.rs)
		prog.BindBound(2, bd)
		prog.BindFiller(3, filler)
		prog.Compute(int(s.bound[2]-s.bound[0]), int(s.bound[3]-s.bound[1]), 1)
		hwpRepeat <- prog
	case *_ImageFillerHorizontalRepeat:
		filler := hwDriver.Filler(f.img)
		defer filler.Delete()
		prog = <-hwpRepeat_h
		prog.Use()
		prog.BindWorkspace(0, s.ws)
		prog.BindResult(1, s.root.rs)
		prog.BindBound(2, bd)
		prog.BindFiller(3, filler)
		prog.Compute(int(s.bound[2]-s.bound[0]), int(s.bound[3]-s.bound[1]), 1)
		hwpRepeat_h <- prog
	case *_ImageFillerVerticalRepeat:
		filler := hwDriver.Filler(f.img)
		defer filler.Delete()
		prog = <-hwpRepeat_v
		prog.Use()
		prog.BindWorkspace(0, s.ws)
		prog.BindResult(1, s.root.rs)
		prog.BindBound(2, bd)
		prog.BindFiller(3, filler)
		prog.Compute(int(s.bound[2]-s.bound[0]), int(s.bound[3]-s.bound[1]), 1)
		hwpRepeat_v <- prog
	}
}
