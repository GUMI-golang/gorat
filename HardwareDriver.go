package gorat

import (
	"bytes"
	"github.com/GUMI-golang/gorat/shaders"
	"github.com/GUMI-golang/gumi/gcore"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/pkg/errors"
	"image"
	"image/color"
	"io"
	"sync"
)

type HardwareDriver interface {
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
var (
	hwDriver HardwareDriver
	hwdLock  = new(sync.Mutex)
	//
	hwpOnce            = new(sync.Once)
	hwpStroker         chan HardwareProgram
	hwpFiller          chan HardwareProgram
	hwpColor           chan HardwareProgram
	hwpFixed           chan HardwareProgram
	hwpGaussian        chan HardwareProgram
	hwpNearest         chan HardwareProgram
	hwpNearestNeighbor chan HardwareProgram
	hwpRepeat          chan HardwareProgram
	hwpRepeat_h        chan HardwareProgram
	hwpRepeat_v        chan HardwareProgram
)

func SetupDriver(hwd HardwareDriver) (err error) {
	if hwd == nil {
		return errors.New("Driver not nilable")
	}
	if CheckDriver() {
		return errors.New("Already have driver")
	}
	hwdLock.Lock()
	defer hwdLock.Unlock()
	// DriverSetuo Must run once
	hwDriver = hwd
	//
	hwpStroker = make(chan HardwareProgram, 1)
	hwpFiller = make(chan HardwareProgram, 1)
	hwpColor = make(chan HardwareProgram, 1)
	hwpFixed = make(chan HardwareProgram, 1)
	hwpGaussian = make(chan HardwareProgram, 1)
	hwpNearest = make(chan HardwareProgram, 1)
	hwpNearestNeighbor = make(chan HardwareProgram, 1)
	hwpRepeat = make(chan HardwareProgram, 1)
	hwpRepeat_h = make(chan HardwareProgram, 1)
	hwpRepeat_v = make(chan HardwareProgram, 1)
	//
	return nil
}
func Driver() HardwareDriver {
	return hwDriver
}
func CheckDriver() bool {
	hwdLock.Lock()
	defer hwdLock.Unlock()
	return hwDriver != nil
}
func checkDriverOrPanic() {
	hwpOnce.Do(func() {
		hwpStroker <- gcore.MustValue(hwDriver.ComputeProgram(bytes.NewReader(shaders.MustAsset("glGorat-Stroker-type2.cs.glsl")))).(HardwareProgram)
		hwpFiller <- gcore.MustValue(hwDriver.ComputeProgram(bytes.NewReader(shaders.MustAsset("glGorat-Filler.cs.glsl")))).(HardwareProgram)
		hwpColor <- gcore.MustValue(hwDriver.ComputeProgram(bytes.NewReader(shaders.MustAsset("glGorat-Fill-Color.cs.glsl")))).(HardwareProgram)
		hwpFixed <- gcore.MustValue(hwDriver.ComputeProgram(bytes.NewReader(shaders.MustAsset("glGorat-Fill-Fixed.cs.glsl")))).(HardwareProgram)
		hwpGaussian <- gcore.MustValue(hwDriver.ComputeProgram(bytes.NewReader(shaders.MustAsset("glGorat-Fill-Gaussian.cs.glsl")))).(HardwareProgram)
		hwpNearest <- gcore.MustValue(hwDriver.ComputeProgram(bytes.NewReader(shaders.MustAsset("glGorat-Fill-Nearest.cs.glsl")))).(HardwareProgram)
		hwpNearestNeighbor <- gcore.MustValue(hwDriver.ComputeProgram(bytes.NewReader(shaders.MustAsset("glGorat-Fill-NearestNeighbor.cs.glsl")))).(HardwareProgram)
		hwpRepeat <- gcore.MustValue(hwDriver.ComputeProgram(bytes.NewReader(shaders.MustAsset("glGorat-Fill-Repeat.cs.glsl")))).(HardwareProgram)
		hwpRepeat_h <- gcore.MustValue(hwDriver.ComputeProgram(bytes.NewReader(shaders.MustAsset("glGorat-Fill-Repeat-horizontal.cs.glsl")))).(HardwareProgram)
		hwpRepeat_v <- gcore.MustValue(hwDriver.ComputeProgram(bytes.NewReader(shaders.MustAsset("glGorat-Fill-Repeat-vertical.cs.glsl")))).(HardwareProgram)
	})
	if hwDriver == nil {
		panic("SubHardware driver not ready")
	}
}
