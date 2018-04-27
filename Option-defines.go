package gorat

type StrokeJoin uint8
type StrokeCap uint8

const (
	StrokeJoinBevel StrokeJoin = iota
	StrokeJoinRound StrokeJoin = iota
	StrokeJoinMiter StrokeJoin = iota
)
const (
	StrokeCapButt   StrokeCap = iota
	StrokeCapRound  StrokeCap = iota
	StrokeCapSqaure StrokeCap = iota
)