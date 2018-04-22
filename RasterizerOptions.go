package gorat

import (
	"image/color"
	"image"
)

type StrokeJoin uint8

const (
	StrokeJoinBevel StrokeJoin = iota
	StrokeJoinRound StrokeJoin = iota
	StrokeJoinMiter StrokeJoin = iota
)

type StrokeCap uint8

const (
	StrokeCapButt   StrokeCap = iota
	StrokeCapRound  StrokeCap = iota
	StrokeCapSqaure StrokeCap = iota
)

type Overlap uint8

const (
	OverlapVacate     Overlap = iota
	OverlapCWMinusCCW Overlap = iota
	OverlapCCWMinusCW Overlap = iota
	OverlapPlus       Overlap = iota
	OverlapCW         Overlap = iota
	OverlapCCW        Overlap = iota
	OverlapAverage    Overlap = iota
)

type Filler interface {
	RGBA(x, y int) color.RGBA
}
type FillerWithBound interface {
	Filler
	To(r image.Rectangle)
}

type AntiAliasing uint8

const (
	AntiAliasing1x = 0
	AntiAliasing2x = 2
	AntiAliasing4x = 4
	AntiAliasing8x = 8
)