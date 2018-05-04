package fwrat

import (
	"github.com/GUMI-golang/gorat"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/GUMI-golang/gumi/gcore"
	"bytes"
	"github.com/GUMI-golang/gorat/shaders"
	"github.com/pkg/errors"
)

type Offscreen struct {
	wnd *glfw.Window
	drv gorat.Driver
	//
	hwpStroker         chan gorat.HardwareProgram
	hwpFiller          chan gorat.HardwareProgram
	hwpColor           chan gorat.HardwareProgram
	hwpFixed           chan gorat.HardwareProgram
	hwpGaussian        chan gorat.HardwareProgram
	hwpNearest         chan gorat.HardwareProgram
	hwpNearestNeighbor chan gorat.HardwareProgram
	hwpRepeat          chan gorat.HardwareProgram
	hwpRepeat_h        chan gorat.HardwareProgram
	hwpRepeat_v        chan gorat.HardwareProgram
	//
	hwpuMixing         chan gorat.HardwareProgram
}

func (s *Offscreen) UtilLoadMixing() gorat.HardwareProgram {
	return <- s.hwpuMixing
}

func (s *Offscreen) UtilUnloadMixing(p gorat.HardwareProgram) {
	s.hwpuMixing <- p
}

func CreateOffscreenContext(driver gorat.Driver) (*Offscreen, error) {
	if driver == nil{
		return nil, errors.New("Driver can't be nil")
	}
	//
	major, minor := driver.Version()
	glfw.WindowHint(glfw.ContextVersionMajor, major)
	glfw.WindowHint(glfw.ContextVersionMinor, minor)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.Visible, glfw.False)
	wnd, err := glfw.CreateWindow(10, 10, "", nil, nil)
	if err != nil {
		return nil, err
	}
	wnd.MakeContextCurrent()
	//
	err = driver.Init()
	if err != nil {
		return nil, err
	}

	res := &Offscreen{
		wnd:                wnd,
		drv: driver,
		hwpStroker:         make(chan gorat.HardwareProgram, 1),
		hwpFiller:          make(chan gorat.HardwareProgram, 1),
		hwpColor:           make(chan gorat.HardwareProgram, 1),
		hwpFixed:           make(chan gorat.HardwareProgram, 1),
		hwpGaussian:        make(chan gorat.HardwareProgram, 1),
		hwpNearest:         make(chan gorat.HardwareProgram, 1),
		hwpNearestNeighbor: make(chan gorat.HardwareProgram, 1),
		hwpRepeat:          make(chan gorat.HardwareProgram, 1),
		hwpRepeat_h:        make(chan gorat.HardwareProgram, 1),
		hwpRepeat_v:        make(chan gorat.HardwareProgram, 1),
		hwpuMixing:        make(chan gorat.HardwareProgram, 1),
	}
	res.hwpStroker <- gcore.MustValue(driver.ComputeProgram(bytes.NewReader(shaders.MustAsset("glGorat-Stroker-type2.cs.glsl")))).(gorat.HardwareProgram)
	res.hwpFiller <- gcore.MustValue(driver.ComputeProgram(bytes.NewReader(shaders.MustAsset("glGorat-Filler.cs.glsl")))).(gorat.HardwareProgram)
	res.hwpColor <- gcore.MustValue(driver.ComputeProgram(bytes.NewReader(shaders.MustAsset("glGorat-Fill-Color.cs.glsl")))).(gorat.HardwareProgram)
	res.hwpFixed <- gcore.MustValue(driver.ComputeProgram(bytes.NewReader(shaders.MustAsset("glGorat-Fill-Fixed.cs.glsl")))).(gorat.HardwareProgram)
	res.hwpGaussian <- gcore.MustValue(driver.ComputeProgram(bytes.NewReader(shaders.MustAsset("glGorat-Fill-Gaussian.cs.glsl")))).(gorat.HardwareProgram)
	res.hwpNearest <- gcore.MustValue(driver.ComputeProgram(bytes.NewReader(shaders.MustAsset("glGorat-Fill-Nearest.cs.glsl")))).(gorat.HardwareProgram)
	res.hwpNearestNeighbor <- gcore.MustValue(driver.ComputeProgram(bytes.NewReader(shaders.MustAsset("glGorat-Fill-NearestNeighbor.cs.glsl")))).(gorat.HardwareProgram)
	res.hwpRepeat <- gcore.MustValue(driver.ComputeProgram(bytes.NewReader(shaders.MustAsset("glGorat-Fill-Repeat.cs.glsl")))).(gorat.HardwareProgram)
	res.hwpRepeat_h <- gcore.MustValue(driver.ComputeProgram(bytes.NewReader(shaders.MustAsset("glGorat-Fill-Repeat-horizontal.cs.glsl")))).(gorat.HardwareProgram)
	res.hwpRepeat_v <- gcore.MustValue(driver.ComputeProgram(bytes.NewReader(shaders.MustAsset("glGorat-Fill-Repeat-vertical.cs.glsl")))).(gorat.HardwareProgram)
	res.hwpuMixing <- gcore.MustValue(driver.ComputeProgram(bytes.NewReader(shaders.MustAsset("glGorat-Util-Mixing.cs.glsl")))).(gorat.HardwareProgram)
	return res, nil
}

func (s *Offscreen) Driver() gorat.Driver {
	return s.drv
}

func (s *Offscreen) LoadStroker() gorat.HardwareProgram {
	return <-s.hwpStroker
}
func (s *Offscreen) LoadFiller() gorat.HardwareProgram {
	return <-s.hwpFiller
}
func (s *Offscreen) LoadColor() gorat.HardwareProgram {
	return <-s.hwpColor
}
func (s *Offscreen) LoadFixed() gorat.HardwareProgram {
	return <-s.hwpFixed
}
func (s *Offscreen) LoadGaussian() gorat.HardwareProgram {
	return <-s.hwpGaussian
}
func (s *Offscreen) LoadNearest() gorat.HardwareProgram {
	return <-s.hwpNearest
}
func (s *Offscreen) LoadNearestNeighbor() gorat.HardwareProgram {
	return <-s.hwpNearestNeighbor
}
func (s *Offscreen) LoadRepeat() gorat.HardwareProgram {
	return <-s.hwpRepeat
}
func (s *Offscreen) LoadRepeatHorizontal() gorat.HardwareProgram {
	return <-s.hwpRepeat_h
}
func (s *Offscreen) LoadRepeatVertical() gorat.HardwareProgram {
	return <-s.hwpRepeat_v
}

func (s *Offscreen) UnloadStroker(p gorat.HardwareProgram) {
	s.hwpStroker <- p
}
func (s *Offscreen) UnloadFiller(p gorat.HardwareProgram) {
	s.hwpFiller <- p
}
func (s *Offscreen) UnloadColor(p gorat.HardwareProgram) {
	s.hwpColor <- p
}
func (s *Offscreen) UnloadFixed(p gorat.HardwareProgram) {
	s.hwpFixed <- p
}
func (s *Offscreen) UnloadGaussian(p gorat.HardwareProgram) {
	s.hwpGaussian <- p
}
func (s *Offscreen) UnloadNearest(p gorat.HardwareProgram) {
	s.hwpNearest <- p
}
func (s *Offscreen) UnloadNearestNeighbor(p gorat.HardwareProgram) {
	s.hwpNearestNeighbor <- p
}
func (s *Offscreen) UnloadRepeat(p gorat.HardwareProgram) {
	s.hwpRepeat <- p
}
func (s *Offscreen) UnloadRepeatHorizontal(p gorat.HardwareProgram) {
	s.hwpRepeat_h <- p
}
func (s *Offscreen) UnloadRepeatVertical(p gorat.HardwareProgram) {
	s.hwpRepeat_v <- p
}

func (s *Offscreen) GetWindow() *glfw.Window {
	return s.wnd
}
func (s *Offscreen) Delete() {
	s.LoadStroker().Delete()
	s.LoadFiller().Delete()
	s.LoadColor().Delete()
	s.LoadFixed().Delete()
	s.LoadGaussian().Delete()
	s.LoadNearest().Delete()
	s.LoadNearestNeighbor().Delete()
	s.LoadRepeat().Delete()
	s.LoadRepeatHorizontal().Delete()
	s.LoadRepeatVertical().Delete()
	s.UtilLoadMixing().Delete()
	//
	close(s.hwpStroker)
	close(s.hwpFiller)
	close(s.hwpColor)
	close(s.hwpFixed)
	close(s.hwpGaussian)
	close(s.hwpNearest)
	close(s.hwpNearestNeighbor)
	close(s.hwpRepeat)
	close(s.hwpRepeat_h)
	close(s.hwpRepeat_v)
	close(s.hwpuMixing)
	s.wnd.Destroy()
}