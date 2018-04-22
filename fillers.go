package gorat

import (
	"image"
	"image/color"
	"image/draw"
)

type ColorFiller color.RGBA

func NewColorFiller(r, g, b, a uint8) ColorFiller {
	return ColorFiller{r, g, b, a}
}
func (s ColorFiller) RGBA(x, y int) color.RGBA {
	return color.RGBA(s)
}

type ImageFillerMode uint8

const (
	ImageFillerFixed           ImageFillerMode = iota
	ImageFillerNearest         ImageFillerMode = iota
	ImageFillerNearestNeighbor ImageFillerMode = iota
	ImageFillerGausian         ImageFillerMode = iota
)

func NewImageFiller(i image.Image, mode ImageFillerMode) Filler {
	temp := image.NewRGBA(i.Bounds())
	draw.Draw(temp, temp.Rect, i, i.Bounds().Min, draw.Src)
	switch mode {
	case ImageFillerFixed:
		return &_ImageFillerFixed{
			img: temp,
		}
	case ImageFillerGausian:
		return &_ImageFillerFixed{
			img: temp,
		}
	}
	panic("Unknown mode")
}

type _ImageFillerFixed struct {
	img *image.RGBA
}

func (s *_ImageFillerFixed) RGBA(x, y int) color.RGBA {
	offset := s.img.PixOffset(x, y)
	return color.RGBA{
		R: s.img.Pix[offset+r],
		G: s.img.Pix[offset+g],
		B: s.img.Pix[offset+b],
		A: s.img.Pix[offset+a],
	}
}
