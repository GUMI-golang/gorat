package gorat

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Triangle struct {
	A, B, C mgl32.Vec2
}

func (s Triangle) Expand(l float32) Triangle {
	xsum := s.A[0] + s.B[0] + s.C[0]
	ysum := s.A[1] + s.B[1] + s.C[1]
	center := Vec2(xsum/3, ysum/3)
	cntA := s.A.Sub(center)
	cntB := s.B.Sub(center)
	cntC := s.C.Sub(center)
	return Triangle{
		s.A.Add(cntA.Normalize().Mul(l)),
		s.B.Add(cntB.Normalize().Mul(l)),
		s.C.Add(cntC.Normalize().Mul(l)),
	}
}

type AABB struct {
	Min mgl32.Vec2
	Max mgl32.Vec2
}

func (s AABB) In(point mgl32.Vec2) bool {
	return (s.Min[0] <= point[0] && point[0] < s.Max[0]) && (s.Min[1] <= point[1] && point[1] < s.Max[1])
}
func NewAABB(points ...mgl32.Vec2) AABB {
	tempxs := make([]float32, len(points))
	tempys := make([]float32, len(points))
	for i, p := range points {
		tempxs[i] = p[0]
		tempys[i] = p[1]
	}
	return AABB{
		Min: [2]float32{
			minimal(tempxs...),
			minimal(tempys...),
		},
		Max: [2]float32{
			maximum(tempxs...),
			maximum(tempys...),
		},
	}
}
func (s Triangle) AABB() (aabb AABB) {
	return NewAABB(s.A, s.B, s.C)
}
func Cross(a, b mgl32.Vec2) float32 {
	return a[0]*b[1] - a[1]*b[0]
}

func (s Triangle) UVW(p mgl32.Vec2) (inside bool, u, v, w float32) {

	pa := s.A.Sub(p)
	pb := s.B.Sub(p)
	pc := s.C.Sub(p)
	//p radius(not exact)
	//const radius = 0.4
	//pa = pa.Sub(pa.Normalize().Mul(radius))
	//pb = pb.Sub(pb.Normalize().Mul(radius))
	//pc = pc.Sub(pc.Normalize().Mul(radius))
	u = Cross(pb, pc)
	v = Cross(pc, pa)
	w = Cross(pa, pb)
	if s.RotateDirection() == CCW {
		u = -u
		v = -v
		w = -w
	}
	if u < 0 || v < 0 || w < 0 {
		return false, 0, 0, 0
	}
	total := u + v + w
	u, v, w = u/total, v/total, w/total
	return true, u, v, w

	//ab := s.B.Sub(s.A)
	//bc := s.C.Sub(s.B)
	//ac := s.C.Sub(s.A)
	//area := Cross(ab, ac)
	//rot := float32(1)
	//if area < 0 {
	//	area = -area
	//	rot = float32(-1)
	//}
	//// too small triangle ignore
	//const AREATH = 1
	//if area < AREATH {
	//	return false, 0, 0, 0
	//}
	//ap := p.Sub(s.A)
	//bp := p.Sub(s.B)
	//cp := p.Sub(s.C)
	//u = Cross(bc, bp) * rot
	//v = -Cross(ac, cp) * rot
	//w = Cross(ab, ap) * rot
	//const TH = -0.25
	//if u < TH{
	//	return false, 0, 0, 0
	//}
	//if v < TH{
	//	return false, 0, 0, 0
	//}
	//if w < TH{
	//	return false, 0, 0, 0
	//}
	//
	//if u/area + v / area + w /area < 0.9999{
	//	fmt.Println(p,u,v ,w, area)
	//}
	//const r = 0.51
	//p0 := Vec2(p[0]-r, p[1]-r)
	//p1 := Vec2(p[0]+r, p[1]-r)
	//p2 := Vec2(p[0]-r, p[1]+r)
	//p3 := Vec2(p[0]+r, p[1]+r)
	//pi, pu, pv, pw := uvwHelper(p0, s.A, s.B, s.C, ab, bc, ac, rot)
	//inside = inside || pi
	//u += pu
	//v += pv
	//w += pw
	//pi, pu, pv, pw = uvwHelper(p1, s.A, s.B, s.C, ab, bc, ac, rot)
	//inside = inside || pi
	//u += pu
	//v += pv
	//w += pw
	//pi, pu, pv, pw = uvwHelper(p2, s.A, s.B, s.C, ab, bc, ac, rot)
	//inside = inside || pi
	//u += pu
	//v += pv
	//w += pw
	//pi, pu, pv, pw = uvwHelper(p3, s.A, s.B, s.C, ab, bc, ac, rot)
	//inside = inside || pi
	//u += pu
	//v += pv
	//w += pw
	//return true, u / area , v / area , w / area
}
func (s Triangle) AAUVW(p mgl32.Vec2, aa AntiAliasing) (inside bool, u, v, w float32) {
	const aaradius = 0.25
	//
	var samples []mgl32.Vec2
	switch aa {
	case AntiAliasing1x:
		return s.UVW(p)
	case AntiAliasing2x:
		samples = []mgl32.Vec2{
			{p[0] - aaradius, p[1]},
			{p[0] + aaradius, p[1]},
		}
	case AntiAliasing4x:
		samples = []mgl32.Vec2{
			{p[0] - aaradius, p[1] - aaradius},
			{p[0] + aaradius, p[1] - aaradius},
			{p[0] - aaradius, p[1] + aaradius},
			{p[0] + aaradius, p[1] + aaradius},
		}
	case AntiAliasing8x:
		samples = []mgl32.Vec2{
			{p[0] - aaradius, p[1] - aaradius},
			{p[0] - aaradius, p[1]},
			{p[0] - aaradius, p[1] + aaradius},
			{p[0]			, p[1] + aaradius},
			{p[0] + aaradius, p[1] + aaradius},
			{p[0] + aaradius, p[1] },
			{p[0] + aaradius, p[1] - aaradius},
			{p[0]			, p[1] - aaradius},
		}
	default:
		panic("Undefined AA scale")
	}
	count := float32(aa)
	for _, sample := range samples {
		ti, tu, tv, tw := s.UVW(sample)
		inside = inside || ti
		u += tu
		v += tv
		w += tw
	}
	u /= count
	v /= count
	w /= count
	return inside, u, v, w
}

type RotateDirection uint8

func (s RotateDirection) String() string {
	if s == CW {
		return "CW"
	}
	if s == CCW {
		return "CCW"
	}
	return "Unknown"
}

const (
	CW  RotateDirection = iota
	CCW RotateDirection = iota
)

func (s Triangle) RotateDirection() RotateDirection {
	ab := s.B.Sub(s.A)
	ac := s.C.Sub(s.A)
	if Cross(ab, ac) > 0 {
		return CW
	}
	return CCW
}
