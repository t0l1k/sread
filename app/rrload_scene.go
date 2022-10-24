package app

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/t0l1k/sread/ui"
)

type RRLoadScene struct {
	rect                            *ui.Rect
	container                       []ui.Drawable
	btnQuit, btnNewClip, btnHistory *ui.Button
	lblName                         *ui.Label
}

func NewRRLoadScene() *RRLoadScene {
	rect := []int{0, 0, 1, 1}
	bbg, bfg := ui.GetTheme().Get("button bg"), ui.GetTheme().Get("button fg")
	s := &RRLoadScene{
		rect: ui.NewRect(rect),
	}
	s.btnQuit = ui.NewButton("<", rect, bbg, bfg, func(b *ui.Button) {
		ui.GetUi().Pop()
	})
	s.Add(s.btnQuit)
	s.lblName = ui.NewLabel("Rapid reading assistant", rect, bbg, bfg)
	s.Add(s.lblName)
	s.btnNewClip = ui.NewButton("Read from clipboard", rect, bbg, bfg, func(b *ui.Button) {
		sc := NewRapidReadScene()
		sc.LoadBookFromClipboard()
		ui.GetUi().Push(sc)
		log.Printf("Begin new book read from clipboard.")
	})
	s.Add(s.btnNewClip)
	s.btnHistory = ui.NewButton("Reading history", rect, bbg, bfg, func(b *ui.Button) {
		sc := NewRRHistoryScene()
		ui.GetUi().Push(sc)
		log.Printf("Continue.")
	})
	s.Add(s.btnHistory)
	return s
}

func (r *RRLoadScene) Update(dt int) {
	for _, v := range r.container {
		v.Update(dt)
	}
}

func (r *RRLoadScene) Draw(surface *ebiten.Image) {
	for _, v := range r.container {
		v.Draw(surface)
	}
}

func (r *RRLoadScene) Entered() {
	r.Resize()
}

func (r *RRLoadScene) Add(value ui.Drawable) {
	r.container = append(r.container, value)
}

func (r *RRLoadScene) Resize() {
	w, h := ui.GetUi().GetScreenSize()
	r.rect = ui.NewRect([]int{0, 0, w, h})
	x, y, w, h := 0, 0, int(float64(r.rect.GetLowestSize())*0.05), int(float64(r.rect.GetLowestSize())*0.05)
	r.btnQuit.Resize([]int{x, y, w, h})
	w = int(float64(r.rect.GetLowestSize()) * 0.5)
	x = h
	r.lblName.Resize([]int{x, y, w, h})
	h = int(float64(r.rect.GetLowestSize()) * 0.1)
	x, y = (r.rect.W-w)/2, (r.rect.H-h)/2-h
	r.btnNewClip.Resize([]int{x, y, w, h})
	y += int(float64(h) * 1.2)
	r.btnHistory.Resize([]int{x, y, w, h})
}

func (r *RRLoadScene) Quit() {
	for _, v := range r.container {
		v.Close()
	}
}
