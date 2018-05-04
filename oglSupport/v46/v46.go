// TODO : Add DMA support
//
package v46

import (
	"fmt"
	"github.com/GUMI-golang/gorat"
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"image"
	"image/color"
	"io"
	"io/ioutil"
	"math"
	"strings"
)

func Driver() gorat.Driver {
	return driver{}
}

type driver struct{}

func (driver) Init() error {
	return gl.Init()
}
func (driver) Version() (major int, minor int) {
	return 4, 3
}
func (driver) ComputeProgram(source io.Reader) (gorat.HardwareProgram, error) {
	prog := GLProgram(gl.CreateProgram())
	err := prog.compile(source)
	if err != nil {
		return nil, err
	}
	err = prog.link()
	if err != nil {
		return nil, err
	}
	return prog, nil
}
func (driver) WorkSpace(w, h int) gorat.HardwareWorkspace {
	var ws uint32
	gl.GenTextures(1, &ws)
	gl.BindTexture(gl.TEXTURE_2D, ws)
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
	return GLWorkspace(ws)
}
func (driver) Result(w, h int) gorat.HardwareResult {
	var rs uint32
	gl.GenTextures(1, &rs)
	gl.BindTexture(gl.TEXTURE_2D, rs)
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
	return GLResult(rs)
}
func (driver) Bound() gorat.HardwareBound {
	var buf uint32
	gl.GenBuffers(1, &buf)
	return GLBound(buf)
}
func (driver) Color() gorat.HardwareColor {
	var buf uint32
	gl.GenBuffers(1, &buf)
	return GLColor(buf)
}
func (driver) Context() gorat.HardwareContext {
	var buf uint32
	gl.GenBuffers(1, &buf)
	return GLContext(buf)
}
func (driver) Filler(rgba *image.RGBA) gorat.HardwareFiller {
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

type GLProgram uint32

func (s GLProgram) compile(source io.Reader) error {
	shader := gl.CreateShader(gl.COMPUTE_SHADER)
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
func (s GLProgram) link() error {
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
func (s GLProgram) BindWorkspace(i int, o gorat.HardwareWorkspace) {
	gl.BindImageTexture(uint32(i), uint32(o.(GLWorkspace)), 0, false, 0, gl.READ_WRITE, gl.R32I)
}
func (s GLProgram) BindResult(i int, o gorat.HardwareResult) {
	gl.BindImageTexture(uint32(i), uint32(o.(GLResult)), 0, false, 0, gl.WRITE_ONLY, gl.RGBA32F)
}
func (s GLProgram) BindContext(i int, o gorat.HardwareContext) {
	gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, uint32(i), uint32(o.(GLContext)))
}
func (s GLProgram) BindBound(i int, o gorat.HardwareBound) {
	gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, uint32(i), uint32(o.(GLBound)))
}
func (s GLProgram) BindColor(i int, o gorat.HardwareColor) {
	gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, uint32(i), uint32(o.(GLColor)))
}
func (s GLProgram) BindFiller(i int, o gorat.HardwareFiller) {
	gl.BindImageTexture(uint32(i), uint32(o.(GLFiller)), 0, false, 0, gl.READ_ONLY, gl.RGBA32F)
}
func (s GLProgram) Use() {
	gl.UseProgram(uint32(s))
}
func (s GLProgram) Compute(x, y, z int) {
	gl.DispatchCompute(uint32(x), uint32(y), uint32(z))
}

type GLWorkspace uint32

func (s GLWorkspace) Resize(w, h int) {
	gl.BindTexture(gl.TEXTURE_2D, uint32(s))
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.R32I, int32(w), int32(h), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(nil))
}
func (s GLWorkspace) Clear() {
	var w, h = s.Size()
	gl.BindTexture(gl.TEXTURE_2D, uint32(s))
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.R32I, int32(w), int32(h), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(nil))
}
func (s GLWorkspace) Delete() {
	temp := uint32(s)
	gl.DeleteTextures(1, &temp)
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
			if valtemp < 0 {
				val := uint8(-valtemp)
				res.Pix[off+1] = val
				res.Pix[off+3] = val
			} else {
				val := uint8(valtemp)
				res.Pix[off+0] = val
				res.Pix[off+3] = val
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

type GLResult uint32

func (s GLResult) Resize(w, h int) {
	gl.BindTexture(gl.TEXTURE_2D, uint32(s))
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA32F, int32(w), int32(h), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(nil))
}
func (s GLResult) RectClear(r image.Rectangle) {
	gl.BindTexture(gl.TEXTURE_2D, uint32(s))
	gl.TexSubImage2D(gl.TEXTURE_2D, 0, int32(r.Min.X), int32(r.Min.Y), int32(r.Dx()), int32(r.Dy()), gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(nil))
}
func (s GLResult) Clear() {
	var w, h = s.Size()
	gl.BindTexture(gl.TEXTURE_2D, uint32(s))
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA32F, int32(w), int32(h), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(nil))
}
func (s GLResult) Pointer() uint32 {
	return uint32(s)
}
func (s GLResult) Delete() {
	temp := uint32(s)
	gl.DeleteTextures(1, &temp)
}
func (s GLResult) Get() *image.RGBA {
	var w, h = s.Size()
	res := image.NewRGBA(image.Rect(0, 0, w, h))
	var temp = make([]uint8, w*h*4)
	gl.BindTexture(gl.TEXTURE_2D, uint32(s))
	gl.GetTexImage(gl.TEXTURE_2D, 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(temp))
	for y := 0; y < h; y++ {
		copy(res.Pix[res.Stride*y:res.Stride*(y+1)], temp[res.Stride*(h-y-1):res.Stride*(h-y)])
	}
	return res
}
func (s GLResult) Size() (w, h int) {
	var rw, rh int32
	gl.BindTexture(gl.TEXTURE_2D, uint32(s))
	gl.GetTexLevelParameteriv(gl.TEXTURE_2D, 0, gl.TEXTURE_WIDTH, &rw)
	gl.GetTexLevelParameteriv(gl.TEXTURE_2D, 0, gl.TEXTURE_HEIGHT, &rh)
	return int(rw), int(rh)
}

type GLContext uint32

func (s GLContext) Delete() {
	temp := uint32(s)
	gl.DeleteBuffers(1, &temp)
}
func (s GLContext) Set(p ...mgl32.Vec2) {
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, uint32(s))
	gl.BufferData(gl.SHADER_STORAGE_BUFFER, len(p)*4*2, gl.Ptr(p), gl.DYNAMIC_DRAW)
}

type GLBound uint32

func (s GLBound) Delete() {
	temp := uint32(s)
	gl.DeleteBuffers(1, &temp)
}
func (s GLBound) Set(bound image.Rectangle) {
	p := []int32{int32(bound.Min.X), int32(bound.Min.Y), int32(bound.Max.X), int32(bound.Max.Y)}
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, uint32(s))
	gl.BufferData(gl.SHADER_STORAGE_BUFFER, 4*4, gl.Ptr(p), gl.DYNAMIC_COPY)
}
func (s GLBound) Get() image.Rectangle {
	var bound [4]int32
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, uint32(s))
	gl.GetBufferSubData(gl.SHADER_STORAGE_BUFFER, 0, 4*4, gl.Ptr(&bound[0]))
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

func (s GLColor) Delete() {
	temp := uint32(s)
	gl.DeleteBuffers(1, &temp)
}
func (s GLColor) Set(c color.Color) {
	r, g, b, a := c.RGBA()

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
	gl.GetBufferSubData(gl.SHADER_STORAGE_BUFFER, 0, 4*4, gl.Ptr(&c[0]))
	return color.RGBA{
		R: uint8(c[0] * math.MaxUint8),
		G: uint8(c[1] * math.MaxUint8),
		B: uint8(c[2] * math.MaxUint8),
		A: uint8(c[3] * math.MaxUint8),
	}
}

type GLFiller uint32

func (s GLFiller) Delete() {
	temp := uint32(s)
	gl.DeleteTextures(1, &temp)
}
