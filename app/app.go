package app

import "github.com/t0l1k/eui"

func NewApp() *eui.Ui {
	a := eui.GetUi()
	a.SetTitle("Rapid Read")
	k := 2
	w, h := 200*k, 200*k
	a.SetSize(w, h)
	return a
}
