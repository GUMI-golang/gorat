package figurat

import (
	"math"
	"golang.org/x/image/math/fixed"
)

type Angle fixed.Int52_12
type AngleType uint8

const (
	Degree AngleType = iota
	Radian AngleType = iota
)
func BuildAngle(value float64, kind AngleType) Angle {
	switch kind {
	default:
		fallthrough
	case Degree:
		return Angle(value * math.Pi / 180)
	case Radian:
		return Angle(value)
	}
}
const AngleZero Angle = 0
const AnglePi Angle = 12867