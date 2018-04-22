package textrat

import (
	"golang.org/x/image/math/fixed"
	"golang.org/x/image/font"
	"github.com/GUMI-golang/gorat"
	"github.com/GUMI-golang/gumi/gcore"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
)

var Default Font = NewVectorFont(gcore.MustValue(truetype.Parse(goregular.TTF)).(*truetype.Font), 16, font.HintingFull)

type VectorFont struct {
	f    *truetype.Font
	n    FontName
	size fixed.Int52_12
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
		size: gorat.I(size),
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
	s.size = fixed.Int52_12(size) << 12
}
func (s *VectorFont) Hint() font.Hinting {
	return s.hint
}
func (s *VectorFont) SetHint(hint font.Hinting) {
	s.hint = hint
}
//
func (s *VectorFont) Text(ctx gorat.VectorDrawer, text string, point fixed.Point52_12, align gcore.Align) {
	s.PathText(ctx, text, point, align)
	ctx.Fill()
}
func (s *VectorFont) TextInRect(ctx gorat.VectorDrawer, text string, rect fixed.Rectangle52_12, align gcore.Align) {
	s.PathTextInRect(ctx, text, rect, align)
	ctx.Fill()
}
func (s *VectorFont) MeasureText(text string) (fixed.Point52_12) {
	var res fixed.Point52_12
	var previdx truetype.Index = 0
	var scale = gorat.Fixed64ToFixed32(s.size)
	for i, r := range []rune(text) {
		idx := s.f.Index(r)
		if i > 0 {
			res.X += gorat.Fixed32ToFixed64(s.f.Kern(scale, previdx, idx))
		}
		hmat := s.f.HMetric(scale, idx)
		res.X += gorat.Fixed32ToFixed64(hmat.AdvanceWidth + hmat.LeftSideBearing)
		previdx = idx
	}
	res.Y = s.size
	return res
}
func (s *VectorFont) MeasureHeight() (fixed.Int52_12) {
	return s.size
}
func (s *VectorFont) PathText(ctx gorat.VectorDrawer, text string, point fixed.Point52_12, align gcore.Align) {
	s.drawText(
		ctx,
		text,
		point.Add(alignHelp(align, s.MeasureText(text))),
	)
}
func (s *VectorFont) PathTextInRect(ctx gorat.VectorDrawer, text string, rect fixed.Rectangle52_12, align gcore.Align) {
	//sz := s.MeasureText(text)
	v, h := gcore.SplitAlign(align)
	var pt fixed.Point52_12
	switch v {
	case gcore.AlignTop:

	case gcore.AlignVertical:
		pt.Y = rect.Max.Y * 4096 / gorat.I(2)
	case gcore.AlignBottom:
		pt.Y = rect.Max.Y - s.size

	}
	switch h {
	case gcore.AlignLeft:

	case gcore.AlignHorizontal:
		pt.X = rect.Max.X * 4096 / gorat.I(2)
	case gcore.AlignRight:
		pt.X = rect.Max.X
	}
	s.PathText(ctx, text, pt, align)
}
//
func (s *VectorFont) drawText(ctx gorat.VectorDrawer, text string, point fixed.Point52_12) {

	var scale = gorat.Fixed64ToFixed32(s.size)

	var prevIdx truetype.Index = 0
	for i, r := range []rune(text) {
		idx := s.f.Index(r)
		b, err := s.load(idx, scale)
		if err != nil {
			continue
		}
		if i > 0 {
			point.X += gorat.Fixed32ToFixed64(s.f.Kern(scale, prevIdx, idx))
		}
		temp := point
		temp.Y += gorat.Fixed32ToFixed64(b.Bounds.Min.Y)
		raster(ctx, b, temp)
		point.X += gorat.Fixed32ToFixed64(b.AdvanceWidth)
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
