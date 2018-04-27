package v43

import (
	"fmt"
	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"image"
	"io"
	"io/ioutil"
	"strings"
	"math"
	"image/color"
)

type GLProgram uint32

func NewProgram() GLProgram {
	return GLProgram(gl.CreateProgram())
}
func (s GLProgram) Compile(source io.Reader, shaderType uint32) error {
	shader := gl.CreateShader(shaderType)
	bts, err := ioutil.ReadAll(source)
	if err != nil {
		return err
	}
	csources, free := gl.Strs(string(bts) + "\x00")
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	// check compile success
	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return fmt.Errorf("failed to compile %v: %v", source, log)
	}
	gl.AttachShader(uint32(s), shader)
	gl.DeleteShader(shader)
	return nil
}
func (s GLProgram) Link() error {
	gl.LinkProgram(uint32(s))
	var status int32
	gl.GetProgramiv(uint32(s), gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(uint32(s), gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(uint32(s), logLength, nil, gl.Str(log))

		return fmt.Errorf("failed to link program: %v", log)
	}
	return nil
}
func (s GLProgram) Use() {
	gl.UseProgram(uint32(s))
}

//
func (s GLProgram) BindWorkspace(i int, texf32 GLWorkspace) {
	gl.BindImageTexture(uint32(i), uint32(texf32), 0, false, 0, gl.READ_WRITE, gl.R32I)
}
func (s GLProgram) BindTexture(i int, texf32 GLResult) {
	gl.BindImageTexture(uint32(i), uint32(texf32), 0, false, 0, gl.WRITE_ONLY, gl.RGBA32F)
}
func (s GLProgram) BindContext(i int, ctx GLContext) {
	gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, uint32(i), uint32(ctx))
}
func (s GLProgram) BindBound(i int, ctx GLBound) {
	gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, uint32(i), uint32(ctx))
}
func (s GLProgram) BindColor(i int, ctx GLColor) {
	gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, uint32(i), uint32(ctx))
}
func (s GLProgram) BindFiller(i int, filler GLFiller) {
	gl.BindImageTexture(uint32(i), uint32(filler), 0, false, 0, gl.READ_ONLY, gl.RGBA32F)
}
func (s GLProgram) Compute(x, y, z int) {
	gl.DispatchCompute(uint32(x), uint32(y), uint32(z))
}

//
//func (s GLProgram) SetBound(name string, value image.Rectangle ){
//	loc := gl.GetUniformLocation(uint32(s), gl.Str(name + "\n"))
//	gl.Uniform4i(loc, int32(value.Min.X),int32(value.Min.Y),int32(value.Max.X),int32(value.Max.Y))
//}
//func (s GLProgram) GetBound(name string) image.Rectangle {
//	loc := gl.GetUniformLocation(uint32(s), gl.Str(name + "\n"))
//	var bound [4]int32
//	gl.GetnUniformiv(uint32(s), loc, 4 * 4, &bound[0])
//	return image.Rectangle{
//		Min:image.Point{
//			X:int(bound[0]),
//			Y:int(bound[1]),
//		},
//		Max:image.Point{
//			X:int(bound[2]),
//			Y:int(bound[3]),
//		},
//	}
//}

type GLResult uint32

func NewResult(w, h int) GLResult {
	var texture uint32
	gl.GenTextures(1, &texture)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA32F,
		int32(w),
		int32(h),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(nil))
	return GLResult(texture)
}
func (s GLResult) Get() *image.RGBA {
	var w, h = s.Size()
	res := image.NewRGBA(image.Rect(0, 0, w, h))
	gl.BindTexture(gl.TEXTURE_2D, uint32(s))
	gl.GetTexImage(gl.TEXTURE_2D, 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(res.Pix))
	return res
}
func (s GLResult) Size() (w, h int) {
	var rw, rh int32
	gl.BindTexture(gl.TEXTURE_2D, uint32(s))
	gl.GetTexLevelParameteriv(gl.TEXTURE_2D, 0, gl.TEXTURE_WIDTH, &rw)
	gl.GetTexLevelParameteriv(gl.TEXTURE_2D, 0, gl.TEXTURE_HEIGHT, &rh)
	return int(rw), int(rh)
}

type GLFiller uint32

func NewFiller(rgba *image.RGBA) GLFiller {
	var texture uint32
	gl.GenTextures(1, &texture)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA32F,
		int32(rgba.Rect.Dx()),
		int32(rgba.Rect.Dy()),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(&rgba.Pix[0]))
	return GLFiller(texture)
}


type GLContext uint32

func NewContext() GLContext {
	var buf uint32
	gl.GenBuffers(1, &buf)
	return GLContext(buf)
}
func (s GLContext) Set(p ...mgl32.Vec2) {
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, uint32(s))
	gl.BufferData(gl.SHADER_STORAGE_BUFFER, len(p)*4*2, gl.Ptr(p), gl.DYNAMIC_DRAW)
}

type GLWorkspace uint32

func NewWorkspace(w, h int) GLWorkspace {
	var texture uint32
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.R32I,
		int32(w),
		int32(h),
		0,
		gl.RED_INTEGER,
		gl.INT,
		gl.Ptr(nil))
	return GLWorkspace(texture)
}
func (s GLWorkspace) Print() {
	var w, h = s.Size()
	var temp = make([]int32, w*h)
	gl.GetTexImage(gl.TEXTURE_2D, 0, gl.RED_INTEGER, gl.INT, gl.Ptr(temp))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			fmt.Printf("%6d ", temp[x+w*y])
		}
		fmt.Println()
	}
}
func (s GLWorkspace) Visualize() *image.RGBA {
	var w, h = s.Size()
	res := image.NewRGBA(image.Rect(0, 0, w, h))
	var temp = make([]int32, w*h)
	gl.GetTexImage(gl.TEXTURE_2D, 0, gl.RED_INTEGER, gl.INT, gl.Ptr(temp))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			off := res.PixOffset(x, y)
			valtemp := float32(temp[x+w*y]) / math.MaxUint16 * math.MaxUint8
			if valtemp < 0{
				val := uint8(-valtemp)
				res.Pix[off + 1] = val
				res.Pix[off + 3] = val
			}else {
				val := uint8(valtemp)
				res.Pix[off + 0] = val
				res.Pix[off + 3] = val
			}
		}
	}
	return res
}
func (s GLWorkspace) Size() (w, h int) {
	var rw, rh int32
	gl.BindTexture(gl.TEXTURE_2D, uint32(s))
	gl.GetTexLevelParameteriv(gl.TEXTURE_2D, 0, gl.TEXTURE_WIDTH, &rw)
	gl.GetTexLevelParameteriv(gl.TEXTURE_2D, 0, gl.TEXTURE_HEIGHT, &rh)
	return int(rw), int(rh)
}
func (s GLWorkspace) Active(i int) {
	gl.ActiveTexture(uint32(i))
	gl.BindTexture(gl.TEXTURE_2D, uint32(s))
}

type GLBound uint32

func NewBound() GLBound {
	var buf uint32
	gl.GenBuffers(1, &buf)
	return GLBound(buf)
}
func (s GLBound) Set(bound image.Rectangle) {
	p := []int32{int32(bound.Min.X), int32(bound.Min.Y), int32(bound.Max.X), int32(bound.Max.Y)}
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, uint32(s))
	gl.BufferData(gl.SHADER_STORAGE_BUFFER, 4*4, gl.Ptr(p), gl.DYNAMIC_COPY)
}
func (s GLBound) Get() image.Rectangle {
	var bound [4]int32
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, uint32(s))
	gl.GetBufferSubData(gl.SHADER_STORAGE_BUFFER, 0, 4 * 4, gl.Ptr(&bound[0]))
	return image.Rectangle{
		Min: image.Point{
			X: int(bound[0]),
			Y: int(bound[1]),
		},
		Max: image.Point{
			X: int(bound[2]),
			Y: int(bound[3]),
		},
	}
}

type GLColor uint32

func NewColor() GLColor {
	var buf uint32
	gl.GenBuffers(1, &buf)
	return GLColor(buf)
}
func (s GLColor) Set(c color.Color) {
	r,g, b, a := c.RGBA()

	p := []float32{
		float32(r) / math.MaxUint16,
		float32(g) / math.MaxUint16,
		float32(b) / math.MaxUint16,
		float32(a) / math.MaxUint16,
	}
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, uint32(s))
	gl.BufferData(gl.SHADER_STORAGE_BUFFER, 4*4, gl.Ptr(p), gl.DYNAMIC_COPY)
}
func (s GLColor) Get() color.Color {
	var c [4]float32
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, uint32(s))
	gl.GetBufferSubData(gl.SHADER_STORAGE_BUFFER, 0, 4 * 4, gl.Ptr(&c[0]))
	return color.RGBA{
		R:uint8(c[0] * math.MaxUint8),
		G:uint8(c[1] * math.MaxUint8),
		B:uint8(c[2] * math.MaxUint8),
		A:uint8(c[3] * math.MaxUint8),
	}
}
