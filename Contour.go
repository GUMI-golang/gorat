package gorat

import (
	"golang.org/x/image/math/fixed"
	"encoding/binary"
)

// TODO
// Contour exist to reduce overhead
type Contour []byte
type ContourBuilder struct {
	bts [] byte
	points []fixed.Point52_12
}
func NewContourBuilder() *ContourBuilder {
	return &ContourBuilder{
		bts: nil,
	}
}
func (s *ContourBuilder ) Add(points ...  fixed.Point52_12)  {
	s.points = append(s.points, points...)
}
func (s *ContourBuilder ) Build() []Contour {
	return nil
}

func _PackBinary(point fixed.Point52_12) []byte {
	var pt [16]byte
	binary.LittleEndian.PutUint32(pt[:8], uint32(point.X))
	binary.LittleEndian.PutUint32(pt[8:16], uint32(point.Y))
	return pt[:]
}
// Careful! this func is not safe error
func _UnpackBinary(data []byte) (res fixed.Point52_12) {
	res.X = fixed.Int52_12(binary.LittleEndian.Uint64(data[:8]))
	res.Y = fixed.Int52_12(binary.LittleEndian.Uint64(data[8:16]))
	return
}
