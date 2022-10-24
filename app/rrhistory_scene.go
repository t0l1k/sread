package app

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/t0l1k/sread/ui"
)

type RRHistoryScene struct {
	rect      *ui.Rect
	container []ui.Drawable
	btnQuit   *ui.Button
	txtsLst   []*ui.Button
}

func NewRRHistoryScene() *RRHistoryScene {
	rect := []int{0, 0, 1, 1}
	bbg, bfg := ui.GetTheme().Get("button bg"), ui.GetTheme().Get("button fg")
	s := &RRHistoryScene{
		rect: ui.NewRect(rect),
	}
	s.btnQuit = ui.NewButton("<", rect, bbg, bfg, func(b *ui.Button) {
		ui.GetUi().Pop()
	})
	s.Add(s.btnQuit)

	for _, v := range GetHistory().GetList() {
		btn := ui.NewButton(v, rect, bbg, bfg, s.loadBook)
		s.txtsLst = append(s.txtsLst, btn)
		s.Add(btn)
	}

	return s
}

func (r *RRHistoryScene) loadBook(b *ui.Button) {
	fmt.Println(b.GetText())
	sc := NewRapidReadScene()
	sc.LoadBookFromHistory(b.GetText())
	ui.GetUi().Push(sc)
}

func (r *RRHistoryScene) Entered() {
	r.Resize()
}

func (r *RRHistoryScene) Update(dt int) {
	for _, v := range r.container {
		v.Update(dt)
	}
}

func (r *RRHistoryScene) Draw(surface *ebiten.Image) {
	for _, v := range r.container {
		v.Draw(surface)
	}
}

func (r *RRHistoryScene) Add(value ui.Drawable) {
	r.container = append(r.container, value)
}

func (r *RRHistoryScene) Resize() {
	w, h := ui.GetUi().GetScreenSize()
	r.rect = ui.NewRect([]int{0, 0, w, h})
	x, y, w1, h1 := 0, 0, int(float64(r.rect.GetLowestSize())*0.05), int(float64(r.rect.GetLowestSize())*0.05)
	r.btnQuit.Resize([]int{x, y, w1, h1})
	w2, h2 := int(float64(r.rect.W)*0.9), h1
	x = (r.rect.W - w2) / 2

	for _, v := range r.txtsLst {
		v.Resize([]int{x, y, w2, h2})
		y += int(float64(h1) * 1.05)
	}
}

func (r *RRHistoryScene) Quit() {
	for _, v := range r.container {
		v.Close()
	}
}
