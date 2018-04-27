package textrat

import (
	"golang.org/x/image/math/fixed"
	"golang.org/x/image/font"
	"github.com/GUMI-golang/gorat"
	"github.com/GUMI-golang/gumi/gcore"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
	"github.com/go-gl/mathgl/mgl32"
	"image"
)

var Default Font = NewVectorFont(gcore.MustValue(truetype.Parse(goregular.TTF)).(*truetype.Font), 16, font.HintingFull)

type VectorFont struct {
	f    *truetype.Font
	n    FontName
	size fixed.Int26_6
	hint font.Hinting
}

func NewVectorFont(f *truetype.Font, size int, hint font.Hinting) *VectorFont {
	return &VectorFont{
		n: FontName{
			Name:     f.Name(truetype.NameIDFontSubfamily),
			Family:   f.Name(truetype.NameIDFontFamily),
			FullName: f.Name(truetype.NameIDFontFullName),
		},
		f:    f,
		size: fixed.I(size),
		hint: hint,
	}
}
func (s *VectorFont) Name() FontName {
	return s.n
}
func (s *VectorFont) Size() int {
	return s.size.Round()
}
func (s *VectorFont) SetSize(size int) {
	s.size = fixed.I(size)
}
func (s *VectorFont) Hint() font.Hinting {
	return s.hint
}
func (s *VectorFont) SetHint(hint font.Hinting) {
	s.hint = hint
}
//
func (s *VectorFont) Text(ctx gorat.VectorDrawer, text string, point mgl32.Vec2, align gcore.Align) {
	s.PathText(ctx, text, point, align)
	ctx.Fill()
}
func (s *VectorFont) TextInRect(ctx gorat.VectorDrawer, text string, rect image.Rectangle, align gcore.Align) {
	s.PathTextInRect(ctx, text, rect, align)
	ctx.Fill()
}
func (s *VectorFont) MeasureText(text string) (res mgl32.Vec2) {

	var previdx truetype.Index = 0
	for i, r := range []rune(text) {
		idx := s.f.Index(r)
		if i > 0 {
			res[0] += Fint32ToFloat32(s.f.Kern(s.size, previdx, idx))
		}
		hmat := s.f.HMetric(s.size, idx)
		res[0] += Fint32ToFloat32(hmat.AdvanceWidth + hmat.LeftSideBearing)
		previdx = idx
	}
	res[1] = Fint32ToFloat32(s.size)
	return res
}
func (s *VectorFont) MeasureHeight() (float32) {
	return Fint32ToFloat32(s.size)
}
func (s *VectorFont) PathText(ctx gorat.VectorDrawer, text string, point mgl32.Vec2, align gcore.Align) {
	s.drawText(
		ctx,
		text,
		point.Add(alignHelp(align, s.MeasureText(text))),
	)
}
func (s *VectorFont) PathTextInRect(ctx gorat.VectorDrawer, text string, rect image.Rectangle, align gcore.Align) {
	//sz := s.MeasureText(text)
	v, h := gcore.SplitAlign(align)
	var pt mgl32.Vec2
	switch v {
	case gcore.AlignTop:

	case gcore.AlignVertical:
		pt[1] = float32(rect.Max.Y)/ 2
	case gcore.AlignBottom:
		pt[1] = float32(rect.Max.Y) - Fint32ToFloat32(s.size)

	}
	switch h {
	case gcore.AlignLeft:

	case gcore.AlignHorizontal:
		pt[0] = float32(rect.Max.X)/ 2
	case gcore.AlignRight:
		pt[0] = float32(rect.Max.X)
	}
	s.PathText(ctx, text, pt, align)
}
//
func (s *VectorFont) drawText(ctx gorat.VectorDrawer, text string, point mgl32.Vec2) {


	var prevIdx truetype.Index = 0
	for i, r := range []rune(text) {
		idx := s.f.Index(r)
		b, err := s.load(idx, s.size)
		if err != nil {
			continue
		}
		if i > 0 {
			point[0] += Fint32ToFloat32(s.f.Kern(s.size, prevIdx, idx))
		}
		temp := point
		temp[1] += Fint32ToFloat32(b.Bounds.Min.Y)
		raster(ctx, b, temp)
		point[0] += Fint32ToFloat32(b.AdvanceWidth)
		prevIdx = idx
	}
}
func (s *VectorFont) load(i truetype.Index, scale fixed.Int26_6) (*truetype.GlyphBuf, error) {
	// TODO : Cache
	buf := &truetype.GlyphBuf{}
	err := buf.Load(s.f, scale, i, s.hint)
	if err != nil {
		return nil, err
	}

	return buf, nil
}
