package GLRasterizer

import (
	"image"
)

func Simple(rgba *image.RGBA) {
	//gcore.Capture("out", test)
	//glfw.WindowHint(glfw.ContextVersionMajor, 4)
	//glfw.WindowHint(glfw.ContextVersionMinor, 2)
	//glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	//wnd, err := glfw.CreateWindow(500, 400, "test", nil, nil)
	//if err != nil {
	//	panic(err)
	//}
	//b := []byte{0x00, 0x01, 0x02}
	//wnd.MakeContextCurrent()
	//err = gl.Init()
	//if err != nil {
	//	panic(err)
	//}
	////
	//prog := NewProgram()
	//gcore.Must(prog.compile(Load("example/GLRasterizer/glGorat.vs.glsl"), gl.VERTEX_SHADER))
	//gcore.Must(prog.compile(Load("example/GLRasterizer/glGorat.fs.glsl"), gl.FRAGMENT_SHADER))
	//gcore.Must(prog.link())
	////
	//tex := NewResult()
	//tex.Set(test)
	////
	//prog.Use()
	//gl.Enable(gl.STENCIL_TEST)
	//gl.ClearColor(0,0,0,1)
	//gl.BindFragDataLocation(uint32(prog), 0, gl.Str("pix\x00"))
	////
	//gl.CullFace(gl.CCW)
	//gl.CullFace(gl.CW)
	////
	//var vao uint32
	//gl.GenVertexArrays(1, &vao)
	//gl.BindVertexArray(vao)
	////
	//var vbo uint32
	//gl.GenBuffers(1, &vbo)
	//gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	//gl.NamedBufferData(gl.ARRAY_BUFFER, FULLRECTUVSIZE, gl.Ptr(FULLRECTUV), gl.STATIC_DRAW)
	////
	//vertAttrib := uint32(gl.GetAttribLocation(uint32(prog), gl.Str("vertCoord\x00")))
	//gl.EnableVertexAttribArray(vertAttrib)
	//gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))
	////
	//texCoordAttrib := uint32(gl.GetAttribLocation(uint32(prog), gl.Str("vertUV\x00")))
	//gl.EnableVertexAttribArray(texCoordAttrib)
	//gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))
	////
	//for !wnd.ShouldClose() {
	//	gl.Clear(gl.COLOR_BUFFER_BIT | gl.STENCIL_BUFFER_BIT)
	//	// Render
	//	gl.UseProgram(uint32(prog))
	//
	//	gl.BindVertexArray(vao)
	//
	//	gl.ActiveTexture(gl.TEXTURE0)
	//	gl.BindTexture(gl.TEXTURE_2D, uint32(tex))
	//
	//	gl.DrawArrays(gl.TRIANGLES, 0, 6*2*3)
	//
	//	wnd.SwapBuffers()
	//	glfw.PollEvents()
	//}

}
