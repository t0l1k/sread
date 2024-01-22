package scene_read

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/t0l1k/eui"
	"github.com/t0l1k/sread/app/data"
)

type RapidReadScene struct {
	eui.SceneBase
	topBar   *eui.TopBar
	rrLabel  *RRLabel
	rrPlayer *RRPlayer
	inGame   bool
	delay    int
	dt       int
	wpmVar   *eui.IntVar
	idxVar   *eui.StringVar
	book     *data.Book
}

func NewRapidReadScene() *RapidReadScene {
	s := &RapidReadScene{
		dt: 0,
	}
	s.topBar = eui.NewTopBar("Чтение", nil)
	s.Add(s.topBar)
	s.rrLabel = NewRRLabel()
	s.Add(s.rrLabel)
	s.wpmVar = eui.NewIntVar(300)
	s.idxVar = eui.NewStringVar("")
	s.rrPlayer = NewRRPlayer(s.playerButtonLogic, s.wpmVar, s.idxVar)
	s.Add(s.rrPlayer)
	return s
}

func (r *RapidReadScene) playerButtonLogic(b *eui.Button) {
	if b.GetText() == bReset {
		r.resetReadingBook()
	} else if b.GetText() == bPrev {
		r.wherePrevParagraph()
	} else if b.GetText() == bNext {
		r.whereNextParagraph()
	} else if b.GetText() == bPlay {
		r.toggleReading()
	}

}

func (r *RapidReadScene) SetDelay(delay int) int {
	return int(60 / float64(delay) * 1000)
}

func (r *RapidReadScene) LoadBookFromHistory(filename string) {
	r.book = data.LoadBookByFilename(filename)
	r.book.Setup()
	log.Printf("Звгружено:%v из истории с %v.", r.book.GetName(), r.book.GetIndex())
	if r.book.GetStatus() == data.Finished {
		r.book.SetIndex(0)
		r.book.SetStatus(data.InReading)
	}
	r.book.GetParagraph().SetIndex(r.book.GetIndex())
	r.getNextWord()
	r.wpmVar.SetValue(r.book.GetLastSpeed())
	r.delay = r.SetDelay(r.wpmVar.Value().(int))
}

func (r *RapidReadScene) LoadBookFromClipboard() {
	r.book = data.LoadBookFromClipboardAndSave()
	r.getNextWord()
	r.book.SetStatus(data.InReading)
	r.wpmVar.SetValue(r.book.GetLastSpeed())
	r.delay = r.SetDelay(r.wpmVar.Value().(int))
	log.Println("Читать из буфера обмена")
}

func (r *RapidReadScene) getNextWord() {
	if r.book.GetParagraph().NextWord() {
		word := r.book.GetParagraph().Value()
		r.rrLabel.SetText(word)
		str := fmt.Sprintf("(%v/%v)", r.book.GetParagraph().Index()+1, r.book.GetParagraph().Size())
		r.idxVar.SetValue(str)
	} else if r.book.GetParagraph().IsLastWorld() {
		r.book.SetStatus(data.Finished)
		r.inGame = false
	}
}

func (r *RapidReadScene) checkKeypress() {
	step := 60
	if inpututil.IsKeyJustReleased(ebiten.KeySpace) {
		r.toggleReading()
	} else if inpututil.IsKeyJustReleased(ebiten.KeyUp) {
		if r.wpmVar.Value().(int) < step*300 {
			wpm := r.wpmVar.Value().(int)
			wpm += step
			r.wpmVar.SetValue(wpm)
			r.delay = r.SetDelay(wpm)
		}
	} else if inpututil.IsKeyJustReleased(ebiten.KeyDown) {
		if r.wpmVar.Value().(int) > step {
			wpm := r.wpmVar.Value().(int)
			wpm -= step
			r.wpmVar.SetValue(wpm)
			r.delay = r.SetDelay(wpm)
		}
	} else if inpututil.IsKeyJustReleased(ebiten.KeyLeft) {
		r.wherePrevParagraph()
	} else if inpututil.IsKeyJustReleased(ebiten.KeyRight) {
		r.whereNextParagraph()
	} else if inpututil.IsKeyJustReleased(ebiten.KeyR) {
		r.resetReadingBook()
	}
}

func (r *RapidReadScene) resetReadingBook() {
	r.pauseReading()
	r.book.SetIndex(0)
	r.book.SetStatus(data.InReading)
	r.book.GetParagraph().SetIndex(r.book.GetIndex())
	r.getNextWord()
}

func (r *RapidReadScene) wherePrevParagraph() {
	r.pauseReading()
	i := 0
	for r.book.GetParagraph().PrevWord() {
		word := r.book.GetParagraph().Value()
		for _, v := range word {
			if i > 0 && v == '.' {
				r.getNextWord()
				return
			}
		}
		i++
		if r.book.GetParagraph().IsFirstWorld() {
			r.resetReadingBook()
			return
		}
	}
}

func (r *RapidReadScene) whereNextParagraph() {
	r.pauseReading()
	for r.book.GetParagraph().NextWord() {
		word := r.book.GetParagraph().Value()
		for _, v := range word {
			if v == '.' {
				r.getNextWord()
				return
			}
		}
	}
}

func (r *RapidReadScene) pauseReading() {
	r.inGame = false
}

func (r *RapidReadScene) toggleReading() {
	r.inGame = !r.inGame
}

func (r *RapidReadScene) Entered() {
	r.Resize()
}

func (r *RapidReadScene) Quit() {
	r.book.SetIndex(r.book.GetParagraph().Index())
	r.book.SetLastSpeed(r.wpmVar.Value().(int))
	data.GetDb().UpdateBook(r.book)
}

func (r *RapidReadScene) Update(dt int) {
	r.SceneBase.Update(dt)
	if r.inGame {
		r.dt += dt
		if r.dt > r.delay {
			r.dt -= r.delay
			r.getNextWord()
		}
	}
	r.checkKeypress()
}

func (r *RapidReadScene) Resize() {
	w0, h0 := eui.GetUi().Size()
	h1 := int(float64(h0) * 0.05)
	r.topBar.Resize([]int{0, 0, w0, h1})
	w2, h2 := int(float64(w0)*0.95), int(float64(h0)*0.2)
	x, y := (w0-w2)/2, (h0-h2-h1)/2
	r.rrLabel.Resize([]int{x, y, w2, h2})
	r.rrLabel.SetFontSize(h2 / 2)
	y += h2
	r.rrPlayer.Resize([]int{x, y, w2, h2 / 2})
}
