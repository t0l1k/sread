package scene_history

import (
	"os"

	"github.com/t0l1k/eui"
	"github.com/t0l1k/sread/app/data"
	scene_read "github.com/t0l1k/sread/app/scenes/read"
)

type RRHistoryScene struct {
	eui.SceneBase
	topBar             *eui.TopBar
	listView           *eui.ListView
	btnViewClearDialog *eui.Button
	clearDialog        *DialogClearDB
}

func NewRRHistoryScene() *RRHistoryScene {
	s := &RRHistoryScene{}
	s.topBar = eui.NewTopBar("Загрузить из истории чтения", nil)
	s.Add(s.topBar)
	s.listView = eui.NewListView()
	s.Add(s.listView)
	s.btnViewClearDialog = eui.NewButton("Вызвать диалог удаления", func(b *eui.Button) {
		s.btnViewClearDialog.Disable()
		s.clearDialog.Setup()
		s.clearDialog.Visible(true)
		s.listView.Visible(false)
	})
	s.Add(s.btnViewClearDialog)
	s.clearDialog = NewDialogClearDB(func(b *eui.Button) {
		var result []string
		for _, v := range s.clearDialog.list.GetCheckBoxes() {
			if v.IsChecked() {
				result = append(result, v.GetText())
			}
		}
		data.GetDb().DeleteRecords(result)
		s.clearDialog.Visible(false)
		s.listView.Reset()
		s.setupHistory()
		s.listView.Visible(true)
		s.btnViewClearDialog.Enable()
	})
	s.Add(s.clearDialog)
	s.clearDialog.Visible(false)
	return s
}

func (s *RRHistoryScene) setupHistory() {
	if _, err := os.Stat("texts"); os.IsNotExist(err) {
		return
	}
	for _, v := range data.GetDb().GetNames() {
		btn := eui.NewButton(v, s.loadBook)
		s.listView.Add(btn)
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
	r.btnViewClearDialog.Resize([]int{w0 - h*5, 0, h * 5, h})
	h1 := h0 - h
	w2, h2 := int(float64(w0)*0.9), h1
	x := (w0 - w2) / 2
	y := h
	r.listView.Resize([]int{x, y, w2, h2})
	r.listView.Itemsize(h)
	r.clearDialog.Resize([]int{x, y, w2, h2})
}
