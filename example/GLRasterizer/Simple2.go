package GLRasterizer

import (
	"github.com/GUMI-golang/gumi/gcore"
	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"image"
	"math"
	"fmt"
	"os"
	"image/draw"
	"github.com/GUMI-golang/gorat/shaders"
	"bytes"
)

func Simple2(rgba *image.RGBA) {

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	wnd, err := glfw.CreateWindow(500, 400, "test", nil, nil)
	if err != nil {
		panic(err)
	}
	wnd.MakeContextCurrent()
	err = gl.Init()
	if err != nil {
		panic(err)
	}
	//
	cube := gcore.MustValue(os.Open("example/cubes_64.png")).(*os.File)
	defer cube.Close()
	img, _, err := image.Decode(cube)
	if err != nil {
		panic(err)
	}
	cubeimg := image.NewRGBA(img.Bounds())
	draw.Draw(cubeimg, cubeimg.Rect, img, img.Bounds().Min, draw.Src)
	//
	var width, height = 1024, 1024
	progStroker := NewProgram()

	gcore.Must(progStroker.Compile(bytes.NewReader(shaders.MustAsset("glGorat-Stroker.cs.glsl")), gl.COMPUTE_SHADER))
	gcore.Must(progStroker.Link())
	//
	progFiller := NewProgram()
	gcore.Must(progFiller.Compile(bytes.NewReader(shaders.MustAsset("glGorat-Filler.cs.glsl")), gl.COMPUTE_SHADER))
	gcore.Must(progFiller.Link())
	//
	progFill := NewProgram()
	gcore.Must(progFill.Compile(bytes.NewReader(shaders.MustAsset("glGorat-Fill-Gaussian.cs.glsl")), gl.COMPUTE_SHADER))
	gcore.Must(progFill.Link())
	//
	workspace := NewWorkspace(width, height)
	bd := NewBound()
	result := NewResult(width, height)
	filler := NewFiller(cubeimg)
	func(){
		progStroker.Use()

		progStroker.BindWorkspace(0, workspace)
		//
		ctx := NewContext()
		points := []mgl32.Vec2{
			{0,0},
			{0, float32(height)},
			{float32(width), float32(height)},
			{0,0},
		}
		ctx.Set(
			points...
		)

		progStroker.BindContext(1, ctx)
		//
		progStroker.Compute(len(points)-1, 1, 1)
	}()
	fmt.Println(workspace.Size())
	gcore.Capture("astroke", workspace.Visualize())
	func(){
		progFiller.Use()
		progFiller.BindWorkspace(0, workspace)
		bd.Set(image.Rectangle{
			Min: image.Point{
				X: math.MaxInt32,
				Y: math.MaxInt32,
			},
			Max: image.Point{
				X: math.MinInt32,
				Y: math.MinInt32,
			},
		})
		progFiller.BindBound(1, bd)
		progFiller.Compute(height, 1, 1)
	}()
	gcore.Capture("afill", workspace.Visualize())
	func(){

		progFill.Use()
		progFill.BindWorkspace(0, workspace)
		progFill.BindTexture(1, result)
		progFill.BindBound(2, bd)
		progFill.BindFiller(3, filler)
		progFill.Compute(width, height, 1)
	}()
	gcore.Capture("aresult", result.Get())
}
