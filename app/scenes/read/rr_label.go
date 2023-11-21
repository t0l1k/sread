package scene_read

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/t0l1k/eui"
	"github.com/t0l1k/sread/app"
	"golang.org/x/image/font"
)

type RRLabel struct {
	eui.DrawableBase
	image            *ebiten.Image
	text             string
	bg, fg, fgActive color.Color
	fontFace         font.Face
	glyphs           []text.Glyph
	y                int
}

func NewRRLabel() *RRLabel {
	theme := eui.GetUi().GetTheme()
	rr := &RRLabel{
		bg:       theme.Get(app.AppRRLabelBg),
		fg:       theme.Get(app.AppRRLabelFg),
		fgActive: theme.Get(app.AppRRLabelFgActive),
	}
	return rr
}

func (l *RRLabel) Draw(surface *ebiten.Image) {
	if l.Dirty {
		l.Layout()
	}
	op := &ebiten.DrawImageOptions{}
	x, y := l.Rect.Pos()
	op.GeoM.Translate(float64(x), float64(y))
	surface.DrawImage(l.image, op)
}

func (l *RRLabel) SetText(value string) {
	if l.text == value {
		return
	}
	l.text = value
	l.Dirty = true
}

func (l *RRLabel) SetFontSize(value int) {
	l.fontFace = eui.GetFonts().Get(value)
	b := text.BoundString(l.fontFace, l.text)
	l.y = l.Rect.CenterY() - b.Min.Y - b.Dy()/2
	l.Dirty = true
}

func (l *RRLabel) Layout() {
	w0, h0 := l.Rect.Size()
	if l.image == nil {
		l.image = ebiten.NewImage(w0, h0)
	} else {
		l.image.Clear()
	}
	l.image.Fill(l.bg)
	l.initGlyphs()
	wordCenter := l.getCenter(l.text)
	lblCenter := int(float64(l.Rect.W) / 4)
	var x, y int
	if wordCenter == 1 {
		x = lblCenter - int(l.glyphs[wordCenter-1].Image.Bounds().Dx()/2)
	} else {
		x = lblCenter - int(
			(l.glyphs[wordCenter].X+l.glyphs[wordCenter-1].X)/2)
	}
	y = l.y
	op := &ebiten.DrawImageOptions{}
	for i, v := range l.glyphs {
		op.GeoM.Reset()
		op.GeoM.Translate(float64(x), float64(y))
		op.GeoM.Translate(v.X, v.Y)
		op.ColorM.Reset()
		var r, g, b float64 = 0, 0, 0
		if wordCenter-1 == i {
			r = 255
		}
		op.ColorM.Scale(r, g, b, 255)
		l.image.DrawImage(v.Image, op)
	}
	l.drawAttributes(w0, h0, lblCenter)
	l.Dirty = false
}

func (l *RRLabel) initGlyphs() {
	l.glyphs = nil
	l.glyphs = text.AppendGlyphs(l.glyphs, l.fontFace, l.text)
}

func (l *RRLabel) drawAttributes(w, h, center int) {
	ebitenutil.DrawLine(l.image, 0, 0, float64(w), 0, l.fg)
	ebitenutil.DrawLine(l.image, 0, float64(h)-1, float64(w), float64(h)-1, l.fg)
	margin := int(float64(l.Rect.GetLowestSize()) * 0.15)
	x0 := center
	ebitenutil.DrawLine(l.image,
		float64(x0),
		0,
		float64(x0),
		float64(margin), l.fg)
	ebitenutil.DrawLine(l.image,
		float64(x0),
		float64(l.Rect.Bottom()),
		float64(x0),
		float64(l.Rect.H-margin), l.fg)
}

func (l *RRLabel) getCenter(word string) int {
	ln := len([]rune(word))
	if ln < 3 {
		return 1
	}
	return int(float64((ln)+6)) / 4
}

func (l *RRLabel) Resize(value []int) {
	l.Rect = eui.NewRect(value)
	l.image = nil
	l.Dirty = true
}
