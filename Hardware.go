package gorat

import (
	"github.com/go-gl/mathgl/mgl32"
	"image"
	"image/color"
	"math"
	"runtime"
	"log"
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
	w, h := to.Size()
	res := new(Hardware)
	res.DefaultOption()
	res.root = res
	res.bound = mgl32.Vec4{0, 0, float32(w), float32(h)}
	res.ws = call().Driver().WorkSpace(w, h)
	defer back()
	res.rs = to
	return res
}
func _SubHardwareCloser(h *SubHardware) {
	h.ws.Delete()
}
// Rasterizer
func (s *Hardware) Setup(w, h int) {
	s.rs.Resize(w, h)
	s.ws.Resize(w, h)
	s.working = make([]mgl32.Vec2, 0, HARDWAREPOINTMINCAPACITY)
	//s.rs = hwDriver.Result(w, h)
}
// SubRasterizer
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
	temp.ws = call().Driver().WorkSpace(r.Dx(), r.Dy())
	defer back()
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
			{to[vecx] , to[vecy] },
		}
	} else {
		s.CloseTo()
		s.working = append(s.working, mgl32.Vec2{to[vecx] , to[vecy] })
	}
}
func (s *SubHardware) LineTo(to mgl32.Vec2) {

	s.working = append(s.working, mgl32.Vec2{to[vecx] , to[vecy] })
}
func (s *SubHardware) QuadTo(pivot, to mgl32.Vec2) {
	from := s.Point()
	pivot = mgl32.Vec2{pivot[vecx] , pivot[vecy] }
	to = mgl32.Vec2{to[vecx] , to[vecy] }
	for _, to := range quadFromTo(from, pivot, to) {
		s.LineTo(to)
	}
}
func (s *SubHardware) CubeTo(pivot1, pivot2, to mgl32.Vec2) {
	from := s.Point()
	pivot1 = mgl32.Vec2{pivot1[vecx] , pivot1[vecy] }
	pivot2 = mgl32.Vec2{pivot2[vecx] , pivot2[vecy] }
	to = mgl32.Vec2{to[vecx] , to[vecy] }
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
		log.Println("Too little point")
		return
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
		log.Println("Too little point")
		return
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
func (s *SubHardware) fill(defineFiller Filler) {
	var prog HardwareProgram
	var c = call()
	defer back()
	//=======================================================
	// stroking
	prog = c.LoadStroker()
	prog.Use()
	// points context
	ctx := c.Driver().Context()
	defer ctx.Delete()
	// set points data
	ctx.Set(s.working...)
	prog.BindWorkspace(0, s.ws)
	prog.BindContext(1, ctx)
	prog.Compute(len(s.working)-1, 1, 1)
	c.UnloadStroker(prog)
	//=======================================================
	// filling
	prog = c.LoadFiller()
	prog.Use()
	// bound
	fillingBound := c.Driver().Bound()
	defer fillingBound.Delete()
	fillingBound.Set(image.Rectangle{
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
	prog.BindBound(1, fillingBound)
	prog.Compute(int(s.bound[3]-s.bound[1]), 1, 1)
	c.UnloadFiller(prog)
	//=======================================================
	// fillstyling
	resultBound := c.Driver().Bound()
	defer resultBound.Delete()
	resultBound.Set(s.Bound())
	switch f := defineFiller.(type) {
	case ColorFiller:
		clr := c.Driver().Color()
		defer clr.Delete()
		clr.Set(color.RGBA(f))
		prog = c.LoadColor()
		prog.Use()
		prog.BindWorkspace(0, s.ws)
		prog.BindResult(1, s.root.rs)
		prog.BindBound(2, resultBound)
		prog.BindColor(3, clr)
		prog.Compute(int(s.bound[2]-s.bound[0]), int(s.bound[3]-s.bound[1]), 1)
		c.UnloadColor(prog)
	case *_ImageFillerFixed:
		filler := c.Driver().Filler(f.img)
		defer filler.Delete()
		prog = c.LoadFixed()
		prog.Use()
		prog.BindWorkspace(0, s.ws)
		prog.BindResult(1, s.root.rs)
		prog.BindBound(2, resultBound)
		prog.BindFiller(3, filler)
		prog.BindBound(4, fillingBound)
		prog.Compute(int(s.bound[2]-s.bound[0]), int(s.bound[3]-s.bound[1]), 1)
		c.UnloadFixed(prog)
	case *_ImageFillerGaussian:
		filler := c.Driver().Filler(f.img)
		defer filler.Delete()
		prog = c.LoadGaussian()
		prog.Use()
		prog.BindWorkspace(0, s.ws)
		prog.BindResult(1, s.root.rs)
		prog.BindBound(2, resultBound)
		prog.BindFiller(3, filler)
		prog.BindBound(4, fillingBound)
		prog.Compute(int(s.bound[2]-s.bound[0]), int(s.bound[3]-s.bound[1]), 1)
		c.UnloadGaussian(prog)
	case *_ImageFillerNearest:
		filler := c.Driver().Filler(f.img)
		defer filler.Delete()
		prog = c.LoadNearest()
		prog.Use()
		prog.BindWorkspace(0, s.ws)
		prog.BindResult(1, s.root.rs)
		prog.BindBound(2, resultBound)
		prog.BindFiller(3, filler)
		prog.BindBound(4, fillingBound)
		prog.Compute(int(s.bound[2]-s.bound[0]), int(s.bound[3]-s.bound[1]), 1)
		c.UnloadNearest(prog)
	case *_ImageFillerNearestNeighbor:
		filler := c.Driver().Filler(f.img)
		defer filler.Delete()
		prog = c.LoadNearestNeighbor()
		prog.Use()
		prog.BindWorkspace(0, s.ws)
		prog.BindResult(1, s.root.rs)
		prog.BindBound(2, resultBound)
		prog.BindFiller(3, filler)
		prog.BindBound(4, fillingBound)
		prog.Compute(int(s.bound[2]-s.bound[0]), int(s.bound[3]-s.bound[1]), 1)
		c.UnloadNearestNeighbor(prog)
	case *_ImageFillerRepeat:
		filler := c.Driver().Filler(f.img)
		defer filler.Delete()
		prog = c.LoadRepeat()
		prog.Use()
		prog.BindWorkspace(0, s.ws)
		prog.BindResult(1, s.root.rs)
		prog.BindBound(2, resultBound)
		prog.BindFiller(3, filler)
		prog.BindBound(4, fillingBound)
		prog.Compute(int(s.bound[2]-s.bound[0]), int(s.bound[3]-s.bound[1]), 1)
		c.UnloadRepeat(prog)
	case *_ImageFillerHorizontalRepeat:
		filler := c.Driver().Filler(f.img)
		defer filler.Delete()
		prog = c.LoadRepeatHorizontal()
		prog.Use()
		prog.BindWorkspace(0, s.ws)
		prog.BindResult(1, s.root.rs)
		prog.BindBound(2, resultBound)
		prog.BindFiller(3, filler)
		prog.BindBound(4, fillingBound)
		prog.Compute(int(s.bound[2]-s.bound[0]), int(s.bound[3]-s.bound[1]), 1)
		c.UnloadRepeatHorizontal(prog)
	case *_ImageFillerVerticalRepeat:
		filler := c.Driver().Filler(f.img)
		defer filler.Delete()
		prog = c.LoadRepeatVertical()
		prog.Use()
		prog.BindWorkspace(0, s.ws)
		prog.BindResult(1, s.root.rs)
		prog.BindBound(2, resultBound)
		prog.BindFiller(3, filler)
		prog.BindBound(4, fillingBound)
		prog.Compute(int(s.bound[2]-s.bound[0]), int(s.bound[3]-s.bound[1]), 1)
		c.UnloadRepeatVertical(prog)
	}
}