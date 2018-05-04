package gorat

import (
	"bytes"
	"github.com/GUMI-golang/gorat/shaders"
	"github.com/GUMI-golang/gumi/gcore"
	"sync"
)

type Context interface {
	Driver() Driver

	LoadStroker() HardwareProgram
	LoadFiller() HardwareProgram
	LoadColor() HardwareProgram
	LoadFixed() HardwareProgram
	LoadGaussian() HardwareProgram
	LoadNearest() HardwareProgram
	LoadNearestNeighbor() HardwareProgram
	LoadRepeat() HardwareProgram
	LoadRepeatHorizontal() HardwareProgram
	LoadRepeatVertical() HardwareProgram

	UnloadStroker(p HardwareProgram)
	UnloadFiller(p HardwareProgram)
	UnloadColor(p HardwareProgram)
	UnloadFixed(p HardwareProgram)
	UnloadGaussian(p HardwareProgram)
	UnloadNearest(p HardwareProgram)
	UnloadNearestNeighbor(p HardwareProgram)
	UnloadRepeat(p HardwareProgram)
	UnloadRepeatHorizontal(p HardwareProgram)
	UnloadRepeatVertical(p HardwareProgram)
	//	utils
	UtilLoadMixing() HardwareProgram
	UtilUnloadMixing(p HardwareProgram)
	// delete
	Delete()
}

var (
	ctxlock = new(sync.RWMutex)
	ctx     Context
)

func Use(c Context) {
	ctxlock.Lock()
	defer ctxlock.Unlock()
	ctx = c
}
func call() Context {
	ctxlock.RLock()
	return ctx
}
func back() {
	ctxlock.RUnlock()
}

func CreateSimpleContext(driver Driver) Context {
	res := &_Context{
		drv:                driver,
		hwpStroker:         make(chan HardwareProgram, 1),
		hwpFiller:          make(chan HardwareProgram, 1),
		hwpColor:           make(chan HardwareProgram, 1),
		hwpFixed:           make(chan HardwareProgram, 1),
		hwpGaussian:        make(chan HardwareProgram, 1),
		hwpNearest:         make(chan HardwareProgram, 1),
		hwpNearestNeighbor: make(chan HardwareProgram, 1),
		hwpRepeat:          make(chan HardwareProgram, 1),
		hwpRepeat_h:        make(chan HardwareProgram, 1),
		hwpRepeat_v:        make(chan HardwareProgram, 1),
		hwpuMixing:        make(chan HardwareProgram, 1),
	}
	res.hwpStroker <- gcore.MustValue(driver.ComputeProgram(bytes.NewReader(shaders.MustAsset("glGorat-Stroker-type2.cs.glsl")))).(HardwareProgram)
	res.hwpFiller <- gcore.MustValue(driver.ComputeProgram(bytes.NewReader(shaders.MustAsset("glGorat-Filler.cs.glsl")))).(HardwareProgram)
	res.hwpColor <- gcore.MustValue(driver.ComputeProgram(bytes.NewReader(shaders.MustAsset("glGorat-Fill-Color.cs.glsl")))).(HardwareProgram)
	res.hwpFixed <- gcore.MustValue(driver.ComputeProgram(bytes.NewReader(shaders.MustAsset("glGorat-Fill-Fixed.cs.glsl")))).(HardwareProgram)
	res.hwpGaussian <- gcore.MustValue(driver.ComputeProgram(bytes.NewReader(shaders.MustAsset("glGorat-Fill-Gaussian.cs.glsl")))).(HardwareProgram)
	res.hwpNearest <- gcore.MustValue(driver.ComputeProgram(bytes.NewReader(shaders.MustAsset("glGorat-Fill-Nearest.cs.glsl")))).(HardwareProgram)
	res.hwpNearestNeighbor <- gcore.MustValue(driver.ComputeProgram(bytes.NewReader(shaders.MustAsset("glGorat-Fill-NearestNeighbor.cs.glsl")))).(HardwareProgram)
	res.hwpRepeat <- gcore.MustValue(driver.ComputeProgram(bytes.NewReader(shaders.MustAsset("glGorat-Fill-Repeat.cs.glsl")))).(HardwareProgram)
	res.hwpRepeat_h <- gcore.MustValue(driver.ComputeProgram(bytes.NewReader(shaders.MustAsset("glGorat-Fill-Repeat-horizontal.cs.glsl")))).(HardwareProgram)
	res.hwpRepeat_v <- gcore.MustValue(driver.ComputeProgram(bytes.NewReader(shaders.MustAsset("glGorat-Fill-Repeat-vertical.cs.glsl")))).(HardwareProgram)
	res.hwpuMixing <- gcore.MustValue(driver.ComputeProgram(bytes.NewReader(shaders.MustAsset("glGorat-Util-Mixing.cs.glsl")))).(HardwareProgram)
	return res
}

type _Context struct {
	drv                Driver
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
	//
	hwpuMixing chan HardwareProgram
}

func (s *_Context) Delete() {
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
}

func (s *_Context) UtilLoadMixing() HardwareProgram {
	return <- s.hwpuMixing
}
func (s *_Context) UtilUnloadMixing(p HardwareProgram) {
	s.hwpuMixing <- p
}

func (s *_Context) Driver() Driver {
	return s.drv
}
func (s *_Context) LoadStroker() HardwareProgram {
	return <-s.hwpStroker
}
func (s *_Context) LoadFiller() HardwareProgram {
	return <-s.hwpFiller
}
func (s *_Context) LoadColor() HardwareProgram {
	return <-s.hwpColor
}
func (s *_Context) LoadFixed() HardwareProgram {
	return <-s.hwpFixed
}
func (s *_Context) LoadGaussian() HardwareProgram {
	return <-s.hwpGaussian
}
func (s *_Context) LoadNearest() HardwareProgram {
	return <-s.hwpNearest
}
func (s *_Context) LoadNearestNeighbor() HardwareProgram {
	return <-s.hwpNearestNeighbor
}
func (s *_Context) LoadRepeat() HardwareProgram {
	return <-s.hwpRepeat
}
func (s *_Context) LoadRepeatHorizontal() HardwareProgram {
	return <-s.hwpRepeat_h
}
func (s *_Context) LoadRepeatVertical() HardwareProgram {
	return <-s.hwpRepeat_v
}
func (s *_Context) UnloadStroker(p HardwareProgram) {
	s.hwpStroker <- p
}
func (s *_Context) UnloadFiller(p HardwareProgram) {
	s.hwpFiller <- p
}
func (s *_Context) UnloadColor(p HardwareProgram) {
	s.hwpColor <- p
}
func (s *_Context) UnloadFixed(p HardwareProgram) {
	s.hwpFixed <- p
}
func (s *_Context) UnloadGaussian(p HardwareProgram) {
	s.hwpGaussian <- p
}
func (s *_Context) UnloadNearest(p HardwareProgram) {
	s.hwpNearest <- p
}
func (s *_Context) UnloadNearestNeighbor(p HardwareProgram) {
	s.hwpNearestNeighbor <- p
}
func (s *_Context) UnloadRepeat(p HardwareProgram) {
	s.hwpRepeat <- p
}
func (s *_Context) UnloadRepeatHorizontal(p HardwareProgram) {
	s.hwpRepeat_h <- p
}
func (s *_Context) UnloadRepeatVertical(p HardwareProgram) {
	s.hwpRepeat_v <- p
}
