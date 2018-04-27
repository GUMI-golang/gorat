#GoRat
Gorat is vector rasterizer for go

##Recommand
- Need Software(CPU) Rasterizer need for making vector image
- Need Hareware(GPU) accelated rasterizer writen pure go(opengl(glsl) used, but except that there is no dependency)
- Need easy to use path rendering libray

##Limitation
- Hardware accelation need GLSL 4.3+(cause using [imageSize](https://www.khronos.org/opengl/wiki/Image_Load_Store#Image_size), etc)
- GoRat currently working in progress, all api is unstable
- There is a lot improvement
- Hardware(not Software Rasterizer) is pool to multithreading or goroutine it can only used on single thread

##Internal Project
###FWRat
It is annoying to prepare the hardware rasterizer yourself.
So we offer pre-prepared rasterizer.
This requires additional [github.com/go-gl/glfw/v3.2/glfw](https://github.com/go-gl/glfw/v3.2/glfw).
And Must call oglSupport before load FWRat

###TextRat
TextRat is Vector Font path drawer using it use [github.com/golang/freetype/truetype](https://github.com/golang/freetype/truetype)
for parse ttf font file

TextRat not use prerasterized font instead vector graphic

###oglSupport
GoratHardwareRasterizer uses the interface for flexibility.
Originally, you need to create the interface yourself, 
but it provides a predefined interface for [github.com/go-gl/gl/v4.3-core/gl](https://github.com/go-gl/gl/v4.3-core/gl) that you use the most.
This is done to support various versions of the gl version. If you observe only this interface, you can use gorat in vulkan.

##Dependancy
- [github.com/go-gl/mathgl/mgl32](https://github.com/go-gl/mathgl/mgl32)
- [github.com/GUMI-golang/gumi/gcore](https://github.com/GUMI-golang/gumi/gcore) : Gorat is a subproject of the GUMI project.
- [golang.org/x/image](https://golang.org/x/image)

###FWRat
- [github.com/go-gl/glfw/v3.2/glfw](https://github.com/go-gl/glfw/v3.2/glfw)

###TextRat
- [github.com/golang/freetype](https://github.com/golang/freetype)

###oglSupport
- [github.com/go-gl/gl/v4.3-core/gl](https://github.com/go-gl/gl/v4.3-core/gl)

###Build only
- [github.com/a-urth/go-bindata](https://github.com/a-urth/go-bindata) : For GLSL packing. user don't need this

######GoRat was heavily influenced by the golang.org/x/image/vector.