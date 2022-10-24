package app

import (
	"fmt"
	"image/color"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/t0l1k/sread/ui"
	"golang.org/x/image/font"
)

type RRLabel struct {
	text           string
	wordsPerMinute string
	rect           *ui.Rect
	Image          *ebiten.Image
	Dirty, Visible bool
	bg, fg, fg2    color.Color
	fontFace       font.Face
	glyphs         []text.Glyph
	y              int
}

func NewRRLabel(txt string, rect []int, bg, fg, fg2 color.Color, fontSize int) *RRLabel {
	r := ui.NewRect(rect)
	fontFace := ui.GetFonts().Get(fontSize)
	b := text.BoundString(fontFace, txt)
	y := r.CenterY() - b.Min.Y - b.Dy()/2
	return &RRLabel{
		text:           txt,
		rect:           r,
		Image:          nil,
		Dirty:          true,
		Visible:        true,
		bg:             bg,
		fg:             fg,
		fg2:            fg2,
		fontFace:       fontFace,
		y:              y,
		wordsPerMinute: strconv.Itoa(ui.GetPreferences().Get("default words per minute speed").(int)),
	}
}

func (l *RRLabel) SetWordsPerMinute(value string) {
	if l.wordsPerMinute == value {
		return
	}
	l.wordsPerMinute = value
	l.Dirty = true
}

func (l *RRLabel) SetText(value string) {
	if l.text == value {
		return
	}
	l.text = value
	l.Dirty = true
}

func (l *RRLabel) SetFont(value int) {
	l.fontFace = ui.GetFonts().Get(value)
	b := text.BoundString(l.fontFace, l.text)
	l.y = l.rect.CenterY() - b.Min.Y - b.Dy()/2
	l.Dirty = true
}

func (l *RRLabel) initGlyphs() {
	l.glyphs = nil
	l.glyphs = text.AppendGlyphs(l.glyphs, l.fontFace, l.text)
}

func (l *RRLabel) Layout() {
	w, h := l.rect.Size()
	if l.Image == nil {
		l.Image = ebiten.NewImage(w, h)
	} else {
		l.Image.Clear()
	}
	l.Image.Fill(l.bg)
	l.initGlyphs()
	wordCenter := l.getCenter(l.text)
	lblCenter := int(float64(l.rect.W) / 4)
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
			r = 64
		}
		op.ColorM.Scale(r, g, b, 255)
		l.Image.DrawImage(v.Image, op)
	}
	l.drawAttributes(w, h, lblCenter)
	l.Dirty = false
}

func (l *RRLabel) drawAttributes(w, h, center int) {
	ebitenutil.DrawLine(l.Image, 0, 0, float64(w), 0, l.fg)
	ebitenutil.DrawLine(l.Image, 0, float64(h)-1, float64(w), float64(h)-1, l.fg)
	margin := int(float64(l.rect.GetLowestSize()) * 0.1)
	x0 := center
	ebitenutil.DrawLine(l.Image,
		float64(x0),
		0,
		float64(x0),
		float64(margin), l.fg)
	ebitenutil.DrawLine(l.Image,
		float64(x0),
		float64(l.rect.Bottom()),
		float64(x0),
		float64(l.rect.H-margin), l.fg)
	h = int(float64(l.rect.H) * 0.2)
	w = h * 2
	x := l.rect.W - w
	y := l.rect.H - h - 1
	rect := []int{x, y, w, h}
	lbl := ui.NewLabel(l.wordsPerMinute, rect, l.bg, l.fg)
	defer lbl.Close()
	lbl.Draw(l.Image)
}

func (l *RRLabel) getCenter(word string) int {
	ln := len([]rune(word))
	if ln < 3 {
		return 1
	}
	return int(float64((ln)+6)) / 4
}

func (l *RRLabel) Update(dt int) {}

func (l *RRLabel) Draw(surface *ebiten.Image) {
	if l.Dirty {
		l.Layout()
	}
	if l.Visible {
		op := &ebiten.DrawImageOptions{}
		x, y := l.rect.Pos()
		op.GeoM.Translate(float64(x), float64(y))
		surface.DrawImage(l.Image, op)
	}
}

func (l *RRLabel) Resize(rect []int) {
	l.rect = ui.NewRect(rect)
	l.Dirty = true
	l.Image = nil
}

func (l RRLabel) String() string {
	return fmt.Sprintf("%v %v", l.text, l.rect)
}

func (l *RRLabel) Close() {
	l.Image.Dispose()
}
