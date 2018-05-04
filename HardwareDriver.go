package gorat

import (
	"github.com/go-gl/mathgl/mgl32"
	"image"
	"image/color"
	"io"
)

type Driver interface {
	Init() error
	Version() (major int, minor int)
	ComputeProgram(source io.Reader) (HardwareProgram, error)
	WorkSpace(w, h int) HardwareWorkspace
	Result(w, h int) HardwareResult
	Bound() HardwareBound
	Color() HardwareColor
	Context() HardwareContext
	Filler(rgba *image.RGBA) HardwareFiller
}
type (
	HardwareProgram interface {
		Use()
		BindWorkspace(i int, o HardwareWorkspace)
		BindResult(i int, o HardwareResult)
		BindContext(i int, o HardwareContext)
		BindBound(i int, o HardwareBound)
		BindColor(i int, o HardwareColor)
		BindFiller(i int, o HardwareFiller)
		Compute(x, y, z int)
		Delete()
	}
	HardwareWorkspace interface {
		Delete()
		Visualize() *image.RGBA
		Clear()
		Size() (w, h int)
		Resize(w, h int)
	}
	HardwareResult interface {
		Pointer() uint32
		Get() *image.RGBA
		Size() (w, h int)
		Delete()
		Clear()
		RectClear(r image.Rectangle)
		Resize(w, h int)
	}
	HardwareContext interface {
		Set(p ...mgl32.Vec2)
		Delete()
	}
	HardwareBound interface {
		Set(bound image.Rectangle)

		Get() image.Rectangle
		Delete()
	}
	HardwareColor interface {
		Set(c color.Color)
		Delete()
	}
	HardwareFiller interface {
		Delete()
	}
)
