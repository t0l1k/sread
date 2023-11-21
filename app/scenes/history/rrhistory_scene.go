package scene_history

import (
	"github.com/t0l1k/eui"
	"github.com/t0l1k/sread/app/data"
	scene_read "github.com/t0l1k/sread/app/scenes/read"
)

type RRHistoryScene struct {
	eui.SceneBase
	topBar    *eui.TopBar
	filesList *eui.ListView
}

func NewRRHistoryScene() *RRHistoryScene {
	s := &RRHistoryScene{}
	s.topBar = eui.NewTopBar("Загрузить из истории чтения")
	s.Add(s.topBar)
	s.filesList = eui.NewListView()
	s.Add(s.filesList)
	return s
}

func (s *RRHistoryScene) setupHistory() {
	for _, v := range data.GetDb().GetNames() {
		btn := eui.NewButton(v, s.loadBook)
		s.filesList.Add(btn)
	}
}

func (r *RRHistoryScene) loadBook(b *eui.Button) {
	eui.GetUi().Pop()
	sc := scene_read.NewRapidReadScene()
	sc.LoadBookFromHistory(b.GetText())
	eui.GetUi().Push(sc)
}

func (r *RRHistoryScene) Entered() {
	r.setupHistory()
	r.Resize()
}

func (r *RRHistoryScene) Resize() {
	w0, h0 := eui.GetUi().Size()
	h := int(float64(h0) * 0.1)
	r.topBar.Resize([]int{0, 0, w0, h})
	h1 := h0 - h
	w2, h2 := int(float64(w0)*0.9), h1
	x := (w0 - w2) / 2
	y := h
	r.filesList.Resize([]int{x, y, w2, h2})
	r.filesList.Itemsize(h)
}
