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
	topBar         *eui.TopBar
	statusLine     *eui.Text
	rrLabel        *RRLabel
	inGame         bool
	delay          int
	dt             int
	wordsPerMinute int
	book           *data.Book
}

func NewRapidReadScene() *RapidReadScene {
	s := &RapidReadScene{
		dt: 0,
	}
	s.topBar = eui.NewTopBar("Чтение")
	s.Add(s.topBar)
	s.rrLabel = NewRRLabel()
	s.Add(s.rrLabel)
	s.statusLine = eui.NewText("")
	s.Add(s.statusLine)
	return s
}

func (r *RapidReadScene) SetDelay(delay int) int {
	return int(60 / float64(delay) * 1000)
}

func (r *RapidReadScene) LoadBookFromHistory(filename string) {
	r.book = data.LoadBookByFilename(filename)
	r.book.Setup()
	log.Printf("Звгружено:%v из истории с %v.", r.book.GetFileName(), r.book.GetIndex())
	if r.book.GetStatus() == data.Finished {
		r.book.SetIndex(0)
		r.book.SetStatus(data.InReading)
	}
	r.book.GetParagraph().SetWord(r.book.GetIndex())
	r.getNextWord()
	r.wordsPerMinute = r.book.GetLastSpeed()
	// ui.GetPreferences().Set("default words per minute speed", r.wordsPerMinute)
	r.delay = r.SetDelay(r.wordsPerMinute)
	r.inGame = true
}

func (r *RapidReadScene) LoadBookFromClipboard() {
	r.book = data.LoadBookFromClipboardAndSave()
	r.getNextWord()
	r.book.SetStatus(data.InReading)
	r.inGame = true
	log.Println("Читать из буфера обмена")
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

func (r *RapidReadScene) getNextWord() {
	if r.book.GetParagraph().NextWord() {
		word := r.book.GetParagraph().Value()
		r.rrLabel.SetText(word)
	} else {
		r.book.SetStatus(data.Finished)
		r.inGame = false
	}
}

func (r *RapidReadScene) checkKeypress() {
	step := 60
	if inpututil.IsKeyJustReleased(ebiten.KeySpace) {
		r.inGame = !r.inGame
	} else if inpututil.IsKeyJustReleased(ebiten.KeyUp) {
		if r.wordsPerMinute < step*300 {
			r.wordsPerMinute += step
			// ui.GetPreferences().Set("default words per minute speed", r.wordsPerMinute)
			r.delay = r.SetDelay(r.wordsPerMinute)
		}
		r.statusLine.SetText(fmt.Sprintf("скорость чтения: слов в минуту %v смена слова на скорости:%vms", r.wordsPerMinute, r.delay))
	} else if inpututil.IsKeyJustReleased(ebiten.KeyDown) {
		if r.wordsPerMinute > step {
			r.wordsPerMinute -= step
			// ui.GetPreferences().Set("default words per minute speed", r.wordsPerMinute)
			r.delay = r.SetDelay(r.wordsPerMinute)
		}
		// ui.GetUi().ShowNotification(fmt.Sprintf("Speed down words per minute:%v", r.wordsPerMinute))
		r.statusLine.SetText(fmt.Sprintf("скорость чтения: слов в минуту %v смена слова на скорости:%vms", r.wordsPerMinute, r.delay))
	}
}

func (r *RapidReadScene) Entered() {
	r.wordsPerMinute = 300
	r.delay = r.SetDelay(r.wordsPerMinute)
	r.Resize()
}

func (r *RapidReadScene) Resize() {
	w0, h0 := eui.GetUi().Size()
	h1 := int(float64(h0) * 0.05)
	r.topBar.Resize([]int{0, 0, w0, h1})
	r.statusLine.Resize([]int{0, h0 - h1, w0, h1})
	w2, h2 := int(float64(w0)*0.95), int(float64(h0)*0.2)
	x, y := (w0-w2)/2, (h0-h2-h1)/2
	r.rrLabel.Resize([]int{x, y, w2, h2})
	r.rrLabel.SetFontSize(h2 / 2)
}
