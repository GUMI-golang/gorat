package fwrat

import (
	"github.com/GUMI-golang/gorat"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type Offscreen struct {
	wnd *glfw.Window
}
func OffscreenContext(width, height int) (o *Offscreen) {
	//
	major, minor := gorat.Driver().Version()
	glfw.WindowHint(glfw.ContextVersionMajor, major)
	glfw.WindowHint(glfw.ContextVersionMinor, minor)
	//glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.Visible, glfw.False)
	wnd, err := glfw.CreateWindow(width, height, "", nil, nil)
	if err != nil {
		return
	}
	wnd.MakeContextCurrent()
	err = gorat.Driver().Init()
	return &Offscreen{
		wnd:wnd,
	}
}
func (s *Offscreen ) GetWindow() *glfw.Window {
	return s.wnd
}