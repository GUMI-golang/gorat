package gorat

import (
	"testing"
	"github.com/go-gl/mathgl/mgl32"
)

func TestTriangle_Rotation(t *testing.T) {
	cwtriangle := Triangle{
		A: mgl32.Vec2{-1, -1},
		B: mgl32.Vec2{0, -1},
		C: mgl32.Vec2{1, -1},
	}
	ccwtriangle := Triangle{
		A: mgl32.Vec2{-1, -1},
		B: mgl32.Vec2{1, -1},
		C: mgl32.Vec2{0, 1},
	}
	t.Log("Is CW  :", CW == cwtriangle.RotateDirection())
	if CW != +cwtriangle.RotateDirection(){
		t.Error(cwtriangle, " is CW Triangle")
	}
	t.Log("Is CCW :", CCW == cwtriangle.RotateDirection())
	if CCW != +ccwtriangle.RotateDirection(){
		t.Error(ccwtriangle, " is CCW Triangle")
	}
}
