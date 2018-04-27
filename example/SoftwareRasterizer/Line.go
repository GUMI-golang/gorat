package SoftwareRasterizer

import (
	"image"
	"github.com/GUMI-golang/gorat"
)

func LineOne(to *image.RGBA)  {
	// SoftwareSub rasterizer support rgba Direct and rgba SubImage
	rat := gorat.NewSoftwareRasterizerRGBA(to)
	//
	rat.SetStrokeWidth(20)
	rat.SetStrokeJoin(gorat.StrokeJoinRound)
	//
	rat.MoveTo([2]float32{50,50}, nil)
	rat.LineTo([2]float32{100,150}, nil)
	//rat.LineTo([2]float32{150,150}, nil)
	rat.LineTo([2]float32{150,50}, nil)
	//rat.LineTo([2]float32{250,50}, nil)
	// rat.Close()
	// .Close() call auto
	rat.Stroke()
} 
