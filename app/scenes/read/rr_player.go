package scene_read

import (
	"image/color"

	"github.com/t0l1k/eui"
	"github.com/t0l1k/sread/app"
)

type RRPlayer struct {
	eui.View
	bg, fg                              color.Color
	btnReset, btnPrev, btnPlay, btnNext *eui.Button
	lblWpm                              *eui.Text
}

func NewRRPlayer(fReset, fPrev, fPlay, fNext func(b *eui.Button), wpmVar *eui.IntVar) *RRPlayer {
	theme := eui.GetUi().GetTheme()
	rr := &RRPlayer{
		bg: theme.Get(app.AppRRLabelBg),
		fg: theme.Get(app.AppRRLabelFg),
	}
	rr.SetupView()
	rr.SetHorizontal()
	rr.btnReset = eui.NewButton("|<<", fReset)
	rr.Add(rr.btnReset)
	rr.btnPrev = eui.NewButton("|<", fPrev)
	rr.Add(rr.btnPrev)
	rr.btnPlay = eui.NewButton(">", fPlay)
	rr.Add(rr.btnPlay)
	rr.btnNext = eui.NewButton(">|", fNext)
	rr.Add(rr.btnNext)
	rr.lblWpm = eui.NewText("")
	rr.Add(rr.lblWpm)
	wpmVar.Attach(rr.lblWpm)
	return rr
}

func (l *RRPlayer) Resize(value []int) {
	l.Rect(value)
	l.BoxLayout.Resize(value)
	l.Dirty(true)
}
