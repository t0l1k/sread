package scene_read

import (
	"image/color"

	"github.com/t0l1k/eui"
	"github.com/t0l1k/sread/app"
)

const (
	bReset = "|<<"
	bPrev  = "|<"
	bPlay  = ">"
	bNext  = ">|"
)

type RRPlayer struct {
	eui.View
	bg, fg                              color.Color
	btnReset, btnPrev, btnPlay, btnNext *eui.Button
	lblWpm, lblIndex                    *eui.Text
}

func NewRRPlayer(btnLogic func(b *eui.Button), wpmVar *eui.IntVar, indexVar *eui.StringVar) *RRPlayer {
	theme := eui.GetUi().GetTheme()
	rr := &RRPlayer{
		bg: theme.Get(app.AppRRLabelBg),
		fg: theme.Get(app.AppRRLabelFg),
	}
	rr.SetupView()
	rr.SetHorizontal()
	rr.btnReset = eui.NewButton(bReset, btnLogic)
	rr.Add(rr.btnReset)
	rr.btnPrev = eui.NewButton(bPrev, btnLogic)
	rr.Add(rr.btnPrev)
	rr.btnPlay = eui.NewButton(bPlay, btnLogic)
	rr.Add(rr.btnPlay)
	rr.btnNext = eui.NewButton(bNext, btnLogic)
	rr.Add(rr.btnNext)
	rr.lblWpm = eui.NewText("")
	rr.Add(rr.lblWpm)
	wpmVar.Attach(rr.lblWpm)
	rr.lblIndex = eui.NewText("")
	rr.Add(rr.lblIndex)
	indexVar.Attach(rr.lblIndex)
	return rr
}
