package gorat
//
//import (
//	"github.com/fogleman/gg"
//	"golang.org/vecx/image/math/fixed"
//	"image"
//	"image/Color"
//	"image/draw"
//)
//
//// TODO : subSoftwareRasterizer use gg.Context, make it ownself
//type SoftwareRasterizer struct {
//	subSoftwareRasterizer
//}
//func (s *SoftwareRasterizer) Resize(w, h int) {
//	s.ctx = gg.NewContext(w, h)
//	s.bound = image.Rect(0,0,w,h)
//	s.opt = Options{
//		StrokeDashes: []fixed.Int52_12{I(0)},
//		strokeJoin:   StrokeJoinBevel,
//		strokeCap:    StrokeCapButt,
//		strokeWidth:  1,
//		Color:        Color.Black,
//		filler:       nil,
//	}
//	s.opt.Setup(s)
//}
//
//func NewSoftwareRasterzier(w, h int) *SoftwareRasterizer {
//	temp := new(SoftwareRasterizer)
//	temp.Resize(w, h)
//	return temp
//
//}
////
//type subSoftwareRasterizer struct {
//	ctx *gg.Context
//	opt Options
//	bound image.Rectangle
//	isbegin bool
//
//}
//func (s *subSoftwareRasterizer) Draw(img draw.Image, op draw.Op) {
//	src := s.ctx.Image().(*image.RGBA)
//	src.Rect = s.bound
//	draw.Draw(img, img.Bounds(), src, src.Rect.Min, op)
//}
//func (s *subSoftwareRasterizer) SubRasterizer(subbound fixed.Rectangle52_12) SubRasterizer {
//	sbd := FixedRectToRect(subbound)
//	src := s.ctx.Image().(*image.RGBA)
//	src.Rect = s.bound
//	//
//	subsrc := src.SubImage(sbd).(*image.RGBA)
//	bound := subsrc.Rect
//	subsrc.Rect = image.Rect(0,0,bound.Dx(), bound.Dy())
//	return &subSoftwareRasterizer{
//		ctx: gg.NewContextForRGBA(subsrc),
//		bound: subsrc.Rect,
//		opt:s.opt,
//	}
//}
//
//func (s *subSoftwareRasterizer) NewPath() VectorPath {
//	return NewPath(s)
//}
//
////
//func (s *subSoftwareRasterizer) Size() fixed.Rectangle52_12 {
//
//	return RectToFixedRect(s.bound.Sub(s.bound.Min))
//}
////
//func (s *subSoftwareRasterizer) IsBegin() bool {
//	return s.isbegin
//}
//func (s *subSoftwareRasterizer) Close() {
//	s.ctx.ClosePath()
//	s.isbegin = false
//}
//func (s *subSoftwareRasterizer) MoveTo(to fixed.Point52_12) {
//	s.ctx.MoveTo(float64(Fixed64ToFloat64(to.X)), float64(Fixed64ToFloat64(to.Y)))
//	s.isbegin = true
//}
//func (s *subSoftwareRasterizer) LineTo(to fixed.Point52_12) {
//	s.ctx.LineTo(float64(Fixed64ToFloat64(to.X)), float64(Fixed64ToFloat64(to.Y)))
//	s.isbegin = true
//}
//func (s *subSoftwareRasterizer) QuadTo(pivot, to fixed.Point52_12) {
//	s.ctx.QuadraticTo(
//		float64(Fixed64ToFloat64(pivot.X)), float64(Fixed64ToFloat64(pivot.Y)),
//		float64(Fixed64ToFloat64(to.X)), float64(Fixed64ToFloat64(to.Y)),
//	)
//	s.isbegin = true
//}
//func (s *subSoftwareRasterizer) CubeTo(pivot1, pivot2, to fixed.Point52_12) {
//	s.ctx.CubicTo(
//		float64(Fixed64ToFloat64(pivot1.X)), float64(Fixed64ToFloat64(pivot1.Y)),
//		float64(Fixed64ToFloat64(pivot2.X)), float64(Fixed64ToFloat64(pivot2.Y)),
//		float64(Fixed64ToFloat64(to.X)), float64(Fixed64ToFloat64(to.Y)),
//	)
//	s.isbegin = true
//}
//func (s *subSoftwareRasterizer) Stroke() {
//	s.Close()
//	s.ctx.Stroke()
//
//}
//func (s *subSoftwareRasterizer) Fill() {
//	s.Close()
//	s.ctx.Fill()
//}
//func (s *subSoftwareRasterizer) Clear(rectangle fixed.Rectangle52_12) {
//	s.ctx.Push()
//	defer s.ctx.Pop()
//	s.ctx.SetColor(Color.RGBA{0, 0, 0, 0})
//	s.ctx.DrawRectangle(
//		float64(rectangle.Min.X),
//		float64(rectangle.Min.X),
//		float64(rectangle.Min.X),
//		float64(rectangle.Min.X),
//	)
//}
////
//func (s *subSoftwareRasterizer) StrokeDash(dashes ...fixed.Int52_12) {
//	s.ctx.SetDash(fixed64sToFloat64s(dashes...)...)
//	s.opt.StrokeDashes = dashes
//
//}
//func (s *subSoftwareRasterizer) StrokeJoin(join StrokeJoin) {
//	switch join {
//	case StrokeJoinBevel:
//		s.ctx.SetLineJoin(gg.LineJoinBevel)
//	case StrokeJoinRound:
//		s.ctx.SetLineJoin(gg.LineJoinRound)
//	case StrokeJoinMiter:
//		// TODO : StrokeJoinMiter not support gg
//		s.ctx.SetLineJoin(gg.LineJoinRound)
//	}
//	s.opt.strokeJoin = join
//}
//func (s *subSoftwareRasterizer) StrokeCap(cap StrokeCap) {
//	switch cap {
//	case StrokeCapButt:
//		s.ctx.SetLineCap(gg.LineCapButt)
//	case StrokeCapRound:
//		s.ctx.SetLineCap(gg.LineCapRound)
//	case StrokeCapSqaure:
//		s.ctx.SetLineCap(gg.LineCapSquare)
//	}
//	s.opt.strokeCap = cap
//}
//func (s *subSoftwareRasterizer) StrokeWidth(width fixed.Int52_12) {
//	s.ctx.SetLineWidth(float64(Fixed64ToFloat64(width)))
//	s.opt.strokeWidth = width
//}
//func (s *subSoftwareRasterizer) Color(Color Color.Color) {
//	s.ctx.SetColor(Color)
//	s.opt.Color = Color
//}
//func (s *subSoftwareRasterizer) FillStyle(fstyle Filler) {
//	// s.ctx.SetFillStyle(fstyle)
//	s.opt.filler = fstyle
//}
//
//func (s *subSoftwareRasterizer) Options() Options {
//	return s.opt
//}
//
//func (s *subSoftwareRasterizer) Revert(opts Options) {
//	opts.Setup(s)
//}
////
////func (s *subSoftwareRasterizer) DrawString(text string, vecx, vecy float64) {
////	s.ctx.DrawString(text, vecx, vecy)
////}
////func (s *subSoftwareRasterizer) MeasureString(text string) (vecx, vecy float64) {
////	return s.ctx.MeasureString(text)
////}
