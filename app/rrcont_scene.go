package app

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/t0l1k/sread/ui"
)

type RRContScene struct {
	rect      *ui.Rect
	container []ui.Drawable
	btnQuit   *ui.Button
	txtLst    *ui.List
}

func NewRRContScene() *RRContScene {
	rect := []int{0, 0, 1, 1}
	bbg, bfg := ui.GetTheme().Get("button bg"), ui.GetTheme().Get("button fg")
	s := &RRContScene{
		rect: ui.NewRect(rect),
	}
	s.btnQuit = ui.NewButton("<", rect, bbg, bfg, func(b *ui.Button) {
		ui.GetUi().Pop()
	})
	s.Add(s.btnQuit)
	s.txtLst = ui.NewList(nil, nil, rect, bbg, bfg, 1)
	s.Add(s.txtLst)
	return s
}

func (r *RRContScene) Update(dt int) {
	for _, v := range r.container {
		v.Update(dt)
	}
}

func (r *RRContScene) Draw(surface *ebiten.Image) {
	for _, v := range r.container {
		v.Draw(surface)
	}
}

func (r *RRContScene) Entered() {
	r.Resize()
}

func (r *RRContScene) Add(value ui.Drawable) {
	r.container = append(r.container, value)
}

func (r *RRContScene) Resize() {
	w, h := ui.GetUi().GetScreenSize()
	r.rect = ui.NewRect([]int{0, 0, w, h})
	x, y, w, h := 0, 0, int(float64(r.rect.GetLowestSize())*0.05), int(float64(r.rect.GetLowestSize())*0.05)
	r.btnQuit.Resize([]int{x, y, w, h})
}

func (r *RRContScene) Quit() {
	for _, v := range r.container {
		v.Close()
	}
}
