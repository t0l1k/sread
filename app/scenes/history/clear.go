package scene_history

import (
	"github.com/t0l1k/eui"
	"github.com/t0l1k/sread/app/data"
)

type DialogClearDB struct {
	eui.View
	btnClear *eui.Button
	list     *eui.ListView
	dialFunc func(d *eui.Button)
}

func NewDialogClearDB(f func(d *eui.Button)) *DialogClearDB {
	t := &DialogClearDB{}
	t.SetupView()
	t.dialFunc = f
	t.list = eui.NewListView()
	t.Add(t.list)
	t.btnClear = eui.NewButton("Удалить выбранное", f)
	t.Add(t.btnClear)
	return t
}

func (t *DialogClearDB) Setup() {
	bg := eui.Navy
	fg := eui.Yellow
	t.list.SetupListViewCheckBoxs(data.GetDb().GetNames(), 30, 1, bg, fg, func(b *eui.Checkbox) {
	})
}

func (t *DialogClearDB) Resize(rect []int) {
	t.View.Resize(rect)
	x, y := t.GetRect().Pos()
	hTop := int(float64(t.GetRect().GetLowestSize()) * 0.1)
	w, h := t.GetRect().W, t.GetRect().H-hTop
	t.list.Resize([]int{x, y, w, h})
	t.list.Itemsize(hTop)
	y += h
	t.btnClear.Resize([]int{x, t.GetRect().H - hTop, w, hTop})
	t.Dirty(true)
}
