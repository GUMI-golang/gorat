package gorat

//
//func NewSoftware(size image.Rectangle) *SoftwareRoot {
//	res := new(SoftwareRoot)
//	res.img = image.NewRGBA(size)
//	res.DefaultOption()
//	res.bound = res.img.Rect
//	res.root = res
//	res.Reset()
//	return res
//}
//type SoftwareRoot struct {
//	// result after fill
//	img *image.RGBA
//	SoftwareSub
//}
//
//
//
//func (s *SoftwareRoot) Setup(w, h int) {
//	s.img = image.NewRGBA(image.Rect(0, 0, w, h))
//}
//
//
//func (s *SoftwareRoot) Image() image.Image {
//	return s.img
//}