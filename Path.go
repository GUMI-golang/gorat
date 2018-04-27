package gorat

import "golang.org/x/image/math/fixed"

//
//import (
//	"golang.org/vecx/image/math/fixed"
//	"image/Color"
//)
//
type Path struct {
	to VectorDrawer
	//
	backup Options
	//
	works []interface{}
}
//func (s *Path) Close() {
//	s.works = append(s.works, pClose{
//	})
//	s.isbegin = false
//}
//
type (
	pMoveTo struct {
		to fixed.Point52_12
	}
	pLineTo struct {
		to fixed.Point52_12
	}
	pQuadTo struct {
		pivot, to fixed.Point52_12
	}
	pCubeTo struct {
		pivot1, pivot2, to fixed.Point52_12
	}
	pClose struct {
	}
	pStroke struct {}
	pFill struct {}
	pClear struct {
		bound fixed.Rectangle52_12
	}
	//
	pFillStyle struct {
		fstyle Filler
	}

)
//
////func NewPath(dst SubRasterizer) VectorPath {
////	return &Path{
////		dst:dst,
////		backup:dst.Options(),
////		works:make([]interface{}, 0, 32),
////	}
////}
//func (s *Path) MoveTo(to fixed.Point52_12) {
//	s.works = append(s.works, pMoveTo{
//		to:to,
//	})
//	s.isbegin = true
//}
//func (s *Path) LineTo(to fixed.Point52_12) {
//	s.works = append(s.works, pLineTo{
//		to:to,
//	})
//	s.isbegin = true
//}
//func (s *Path) QuadTo(pivot fixed.Point52_12, to fixed.Point52_12) {
//	s.works = append(s.works, pQuadTo{
//		pivot:pivot,
//		to:to,
//	})
//	s.isbegin = true
//}
//func (s *Path) CubeTo(pivot1 fixed.Point52_12, pivot2 fixed.Point52_12, to fixed.Point52_12) {
//	s.works = append(s.works, pCubeTo{
//		pivot1:pivot1,
//		pivot2:pivot2,
//		to:to,
//	})
//	s.isbegin = true
//}
//func (s *Path) Stroke() {
//	s.works = append(s.works, pStroke{
//
//	})
//}
//func (s *Path) Fill() {
//	s.works = append(s.works, pFill{
//	})
//}
//func (s *Path) Clear(bound fixed.Rectangle52_12) {
//	s.works = append(s.works, pClear{
//		bound:bound,
//	})
//}
//
//func (s *Path) StrokeJoin(join StrokeJoin) {
//	s.works = append(s.works, pStrokeJoin{
//		join:join,
//	})
//	s.opt.strokeJoin = join
//}
//func (s *Path) StrokeCap(cap StrokeCap) {
//	s.works = append(s.works, pStrokeCap{
//		cap:cap,
//	})
//	s.opt.strokeCap = cap
//}
//func (s *Path) StrokeWidth(width float32) {
//	s.works = append(s.works, pStrokeWidth{
//		width:width,
//	})
//	s.opt.strokeWidth = width
//}
//func (s *Path) Color(Color Color.Color) {
//	s.works = append(s.works, pColor{
//		Color:Color,
//	})
//	s.opt.Color = Color
//}
//func (s *Path) FillStyle(fstyle Filler) {
//	s.works = append(s.works, pFillStyle{
//		fstyle:fstyle,
//	})
//	s.opt.filler = fstyle
//}
//
//func (s *Path) Options() Options {
//	return s.opt
//}
//func (s *Path) Revert(opts Options) {
//	s.works = append(s.works, pRevert{
//		opts: opts,
//	})
//	s.opt = opts
//}
//
//func (s *Path) Commit() {
//	//s.backup.Setup(s.dst)
//	for _, w := range s.works {
//		switch t := w.(type) {
//		case pMoveTo:
//			s.dst.MoveTo(t.to)
//		case pLineTo:
//			s.dst.LineTo(t.to)
//		case pQuadTo:
//			s.dst.QuadTo(t.pivot, t.to)
//		case pCubeTo:
//			s.dst.CubeTo(t.pivot1, t.pivot2, t.to)
//		case pStroke:
//			s.dst.Stroke()
//		case pFill:
//			s.dst.Fill()
//		case pClear:
//			s.dst.Clear(t.bound)
//		case pStrokeDash:
//			s.dst.StrokeDash(t.dashes...)
//		case pStrokeJoin:
//			s.dst.StrokeJoin(t.join)
//		case pStrokeCap:
//			s.dst.StrokeCap(t.cap)
//		case pStrokeWidth:
//			s.dst.StrokeWidth(t.width)
//		case pColor:
//			s.dst.Color(t.Color)
//		case pFillStyle:
//			s.dst.FillStyle(t.fstyle)
//		case pRevert:
//			s.dst.Revert(t.opts)
//		}
//	}
//	s.backup.Setup(s.dst)
//}
//
