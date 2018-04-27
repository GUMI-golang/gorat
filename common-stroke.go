package gorat

import (
	"github.com/go-gl/mathgl/mgl32"
	"math"
)

func splitStroke(from []mgl32.Vec2) (to [][]mgl32.Vec2) {
	var temp []mgl32.Vec2
	for _, v := range from {
		if math.IsNaN(float64(v[vecx])) {
			to = append(to, temp[:len(temp)-1])
			temp = nil
			continue
		}
		temp = append(temp, v)
	}
	return
}
func stroke(points []mgl32.Vec2, opt Options) (res []mgl32.Vec2) {
	if len(points) < 2 {
		return nil
	}
	var b0, b1 []mgl32.Vec2
	rot90 := mgl32.Rotate2D(-math.Pi / 2)
	rad := opt.width / 2
	//
	var v0, v1, v2 mgl32.Vec2
	var l01, l12 mgl32.Vec2
	var vv01, vv12 mgl32.Vec2
	// start
	v0 = points[0]
	v1 = points[1]
	l01 = v1.Sub(v0)
	vv01 = rot90.Mul2x1(l01).Normalize().Mul(rad)
	// cap
	switch opt.cap {
	case StrokeCapButt:
	case StrokeCapSqaure:
		r01n := l01.Normalize().Mul(rad)
		b0 = append(b0, v0.Sub(r01n).Sub(vv01), v0.Sub(r01n).Add(vv01))

	case StrokeCapRound:
		r01n := l01.Normalize().Mul(rad)
		b0 = append(b0, quadFromTo(v0.Sub(vv01),v0.Sub(r01n).Sub(vv01), v0.Sub(r01n))...)
		b0 = append(b0, quadFromTo(v0.Sub(r01n), v0.Sub(r01n).Add(vv01), v0.Add(vv01))...)
	}
	b0 = append(b0, v0.Add(vv01), )
	b1 = append(b1, v0.Sub(vv01), )
	// while
	for i := 2; i < len(points); i++ {
		v2 = points[i]
		l12 = v2.Sub(v1)
		vv12 = rot90.Mul2x1(l12).Normalize().Mul(rad)
		//
		cos := float64(l01.Mul(-1).Dot(l12) / l01.Len() / l12.Len())
		dir := vv01.Add(vv12).Normalize().Mul(rad / float32(math.Sqrt((1 - cos)/2)))
		if l01.Dot(dir) < 0 {
			// b0 inside
			b0 = append(b0, v1.Add(dir))
			switch opt.join {
			case StrokeJoinBevel:
				b1 = append(b1,v1.Sub(vv01), v1.Sub(vv12))

			case StrokeJoinMiter:
				b1 = append(b1,v1.Sub(vv01), v1.Sub(dir), v1.Sub(vv12))
			case StrokeJoinRound:

				b1 = append(b1,quadFromTo(v1.Sub(vv01), v1.Sub(dir), v1.Sub(vv12))...)
			}

		} else {
			switch opt.join {
			case StrokeJoinBevel:
				b0 = append(b0, v1.Add(vv01), v1.Add(vv12))
			case StrokeJoinMiter:
				b0 = append(b0, v1.Add(vv01), v1.Add(dir), v1.Add(vv12))
			case StrokeJoinRound:
				b0 = append(b0, quadFromTo(v1.Add(vv01), v1.Add(dir), v1.Add(vv12))...)
			}

			b1 = append(b1, v1.Sub(dir))
		}
		//
		v0 = v1
		v1 = v2
		l01 = l12
		vv01 = vv12
	}
	b0 = append(b0, v1.Add(vv01), )
	b1 = append(b1, v1.Sub(vv01), )
	// last cap
	switch opt.cap {
	case StrokeCapButt:
	case StrokeCapSqaure:
		r01n := l01.Normalize().Mul(rad)
		b1 = append(b1, v1.Add(r01n).Sub(vv01), v1.Add(r01n).Add(vv01))
	case StrokeCapRound:
		r01n := l01.Normalize().Mul(rad)
		b0 = append(b0, quadFromTo(v1.Add(vv01), v1.Add(r01n).Add(vv01), v1.Add(r01n))...)
		b0 = append(b0, quadFromTo(v1.Add(r01n), v1.Add(r01n).Sub(vv01), v1.Sub(vv01))...)
	}
	for i := len(b1) -1; i >=0; i--{
		b0 = append(b0, b1[i])
	}
	return b0
}

