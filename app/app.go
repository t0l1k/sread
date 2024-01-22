package app

import "github.com/t0l1k/eui"

func NewApp() *eui.Ui {
	a := eui.GetUi()
	a.SetTitle("Помошник чтения")
	k := 2
	w, h := 200*k, 200*k
	a.SetSize(w, h)
	setAppTheme()
	return a
}

const (
	AppRRLabelBg eui.ThemeValue = iota + 100
	AppRRLabelFg
	AppRRLabelFgActive
)

func setAppTheme() {
	theme := eui.GetUi().GetTheme()
	theme.Set(eui.SceneBg, eui.Navy)
	theme.Set(eui.SceneFg, eui.YellowGreen)
	theme.Set(eui.ButtonBg, eui.Silver)
	theme.Set(eui.ButtonFg, eui.Black)
	theme.Set(eui.TextBg, eui.GreenYellow)
	theme.Set(eui.TextFg, eui.Black)
	theme.Set(AppRRLabelBg, eui.Silver)
	theme.Set(AppRRLabelFg, eui.Black)
	theme.Set(AppRRLabelFgActive, eui.Red)
}
