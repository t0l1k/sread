package app

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/t0l1k/sread/ui"
)

type RRLoadScene struct {
	rect                                     *ui.Rect
	container                                []ui.Drawable
	btnQuit, btnNewFile, btnNewClip, btnCont *ui.Button
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
	s.btnNewFile = ui.NewButton("Read from file", rect, bbg, bfg, func(b *ui.Button) {
		sc := NewRapidReadScene()
		name := "readme.md"
		sc.LoadBook(name)
		ui.GetUi().Push(sc)
		log.Printf("Begin new book read from %v", name)
	})
	s.Add(s.btnNewFile)
	s.btnNewClip = ui.NewButton("Read from clipboard", rect, bbg, bfg, func(b *ui.Button) {
		sc := NewRapidReadScene()
		name := "1.txt"
		sc.LoadBook(name)
		ui.GetUi().Push(sc)
		log.Printf("Begin new book read from clipboard %v", name)
	})
	s.Add(s.btnNewClip)
	s.btnCont = ui.NewButton("Continue read last", rect, bbg, bfg, func(b *ui.Button) {
		sc := NewRRContScene()
		ui.GetUi().Push(sc)
		log.Printf("Continue.")
	})
	s.Add(s.btnCont)
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
	w, h = int(float64(r.rect.GetLowestSize())*0.5), int(float64(r.rect.GetLowestSize())*0.1)
	x, y = (r.rect.W-w)/2, (r.rect.H-h)/2
	r.btnNewFile.Resize([]int{x, y, w, h})
	y -= int(float64(h) * 1.2)
	r.btnNewClip.Resize([]int{x, y, w, h})
	y -= int(float64(h) * 1.2)
	r.btnCont.Resize([]int{x, y, w, h})
}

func (r *RRLoadScene) Quit() {
	for _, v := range r.container {
		v.Close()
	}
}
