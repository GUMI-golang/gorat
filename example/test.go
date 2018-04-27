package main

import (
	"github.com/GUMI-golang/gumi/gcore"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/gl/v4.3-core/gl"
)

func main() {
	gcore.Must(glfw.Init())
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	//glfw.WindowHint(glfw.)
	wnd0 :=  gcore.MustValue(glfw.CreateWindow(400, 300, "window 0", nil, nil)).(*glfw.Window)
	wnd1 :=  gcore.MustValue(glfw.CreateWindow(400, 300, "window 1", nil, nil)).(*glfw.Window)
	wnd0.MakeContextCurrent()
	gcore.Must(gl.Init())

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var points = []float32{
		0,1,0,
		1,0,0,
		-1,0,0,
	}
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(points), gl.Ptr(&points[0]), gl.STATIC_DRAW);


	for !wnd0.ShouldClose(){
		gl.EnableVertexAttribArray(0);
		gl.BindBuffer(gl.ARRAY_BUFFER, vbo);
		gl.VertexAttribPointer(
			0,                  // 0번째 속성(attribute). 0 이 될 특별한 이유는 없지만, 쉐이더의 레이아웃(layout)와 반드시 맞추어야 합니다.
			3,                  // 크기(size)
			gl.FLOAT,           // 타입(type)
			false,           // 정규화(normalized)?
			0,                  // 다음 요소 까지 간격(stride)
			gl.PtrOffset(0),            // 배열 버퍼의 오프셋(offset; 옮기는 값)
		)
		// 삼각형 그리기!
		gl.DrawArrays(gl.TRIANGLES, 0, 3); // 버텍스 0에서 시작해서; 총 3개의 버텍스로 -> 하나의 삼각형
		gl.DisableVertexAttribArray(0);


		wnd0.SwapBuffers()

		wnd1.SwapBuffers()
		glfw.PollEvents()
	}
}
