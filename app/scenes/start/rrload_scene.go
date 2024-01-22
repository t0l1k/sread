package scene_main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/t0l1k/eui"
	scene_history "github.com/t0l1k/sread/app/scenes/history"
	scene_read "github.com/t0l1k/sread/app/scenes/read"
)

type RRStartScene struct {
	eui.SceneBase
	layout                 *eui.BoxLayout
	topBar                 *eui.TopBar
	btnNewClip, btnHistory *eui.Button
}

func NewRRStartScene() *RRStartScene {
	s := &RRStartScene{}
	s.topBar = eui.NewTopBar("Помощник по быстрому чтению", nil)
	s.Add(s.topBar)
	s.topBar.SetShowStopwatch()
	s.layout = eui.NewVLayout()
	s.btnNewClip = eui.NewButton("Читать из буфера обмена", func(b *eui.Button) {
		sc := scene_read.NewRapidReadScene()
		sc.LoadBookFromClipboard()
		eui.GetUi().Push(sc)
		log.Printf("Читать из буфера обмена")
	})
	s.layout.Add(s.btnNewClip)
	s.btnHistory = eui.NewButton("История чтения", func(b *eui.Button) {
		sc := scene_history.NewRRHistoryScene()
		eui.GetUi().Push(sc)
		log.Printf("История чтения")
	})
	s.layout.Add(s.btnHistory)
	return s
}

func (r *RRStartScene) Update(dt int) {
	r.SceneBase.Update(dt)
	for _, v := range r.layout.Container {
		v.Update(dt)
	}
}

func (r *RRStartScene) Draw(surface *ebiten.Image) {
	r.SceneBase.Draw(surface)
	for _, v := range r.layout.Container {
		v.Draw(surface)
	}
}

func (r *RRStartScene) Entered() {
	r.Resize()
}

func (r *RRStartScene) Resize() {
	w0, h0 := eui.GetUi().Size()
	h1 := int(float64(h0) * 0.1)
	x, y := 0, 0
	r.topBar.Resize([]int{x, y, w0, h1})
	h2 := int(float64(h0) * 0.15)
	r.layout.Resize([]int{x + h2/2, h2, w0 - h2, h0 - h1*2})
}
