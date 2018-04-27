package GLRasterizer

import (
	"github.com/go-gl/glfw/v3.2/glfw"
	"io"
	"io/ioutil"
	"bytes"
)

func init() {
	glfw.Init()
}

func Load(path string) io.Reader {
	bts, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return bytes.NewReader(bts)
}