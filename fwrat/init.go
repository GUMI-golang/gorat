package fwrat

import "github.com/go-gl/glfw/v3.2/glfw"

func init() {
	err := glfw.Init()
	if err != nil {
		panic(err)
	}
}
