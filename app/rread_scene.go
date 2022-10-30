package app

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/t0l1k/sread/ui"
)

type RapidReadScene struct {
	rect           *ui.Rect
	container      []ui.Drawable
	btnQuit        *ui.Button
	rrLabel        *RRLabel
	inGame         bool
	delay          int
	dt             int
	wordsPerMinute int
	paragraphs     *paragraphs
	paragraph      *paragraph
	book           *Book
}

func NewRapidReadScene() *RapidReadScene {
	rect := []int{0, 0, 1, 1}
	bbg, bfg := ui.GetTheme().Get("button bg"), ui.GetTheme().Get("button fg")
	bg, fg, fg2 := ui.GetTheme().Get("bg"), ui.GetTheme().Get("fg"), ui.GetTheme().Get("fg2")
	s := &RapidReadScene{
		rect: ui.NewRect(rect),
		dt:   0,
	}
	s.btnQuit = ui.NewButton("<", rect, bbg, bfg, func(b *ui.Button) {
		ui.GetUi().Pop()
	})
	s.Add(s.btnQuit)
	s.rrLabel = NewRRLabel("", rect, bg, fg, fg2, 1)
	s.Add(s.rrLabel)
	return s
}

func (r *RapidReadScene) SetDelay(delay int) int {
	return int(60 / float64(delay) * 1000)
}

func (r *RapidReadScene) LoadBookFromHistory(filename string) {
	r.book = LoadBookByFilename(filename)
	log.Printf("Loaded file:%v from history at %v %v.", r.book.filename, r.book.idxA, r.book.idxB)
	if r.book.status == finished {
		r.book.idxA = 0
		r.book.idxB = 0
	}
	r.paragraphs = r.book.paragraps
	r.paragraphs.Set(r.book.idxA)
	r.paragraphs.NextParagraph()
	r.paragraph = newParagraph(r.paragraphs.Value())
	r.paragraph.SetWord(r.book.idxB)
	r.getNextWord()
	r.wordsPerMinute = r.book.lastSpeed
	ui.GetPreferences().Set("default words per minute speed", r.wordsPerMinute)
	r.delay = r.SetDelay(r.wordsPerMinute)
	r.rrLabel.SetWordsPerMinute(strconv.Itoa(r.wordsPerMinute))
	r.inGame = true
}

func (r *RapidReadScene) LoadBookFromClipboard() {
	r.book = LoadBookAndSaveFromClipboard()
	r.paragraphs = r.book.paragraps
	r.paragraphs.NextParagraph()
	r.paragraph = newParagraph(r.paragraphs.Value())
	r.getNextWord()
	r.book.status = inReading
	r.inGame = true
	log.Println("Loaded text from clipboard.")
}

func (r *RapidReadScene) Update(dt int) {
	if r.inGame {
		r.dt += dt
		if r.dt > r.delay {
			r.dt -= r.delay
			r.getNextWord()
		}
	}
	r.checkKeypress()
	for _, v := range r.container {
		v.Update(dt)
	}
}

func (r *RapidReadScene) getNextWord() {
	if r.paragraph.NextWord() {
		word := r.paragraph.Value()
		r.rrLabel.SetText(word)
	} else {
		if r.paragraphs.NextParagraph() {
			par := r.paragraphs.Value()
			r.paragraph = newParagraph(par)
		} else {
			r.book.status = finished
			r.inGame = false
		}
	}
}

func (r *RapidReadScene) checkKeypress() {
	step := ui.GetPreferences().Get("step").(int)
	if inpututil.IsKeyJustReleased(ebiten.KeySpace) {
		r.inGame = !r.inGame
		ui.GetUi().ShowNotification("toggle pause")
	} else if inpututil.IsKeyJustReleased(ebiten.KeyUp) {
		if r.wordsPerMinute < step*300 {
			r.wordsPerMinute += step
			ui.GetPreferences().Set("default words per minute speed", r.wordsPerMinute)
			r.delay = r.SetDelay(r.wordsPerMinute)
			r.rrLabel.SetWordsPerMinute(strconv.Itoa(r.wordsPerMinute))
		}
		ui.GetUi().ShowNotification(fmt.Sprintf("Speed up words per minute:%v", r.wordsPerMinute))
		log.Printf("now words per minute:%v delay is:%vms", r.wordsPerMinute, r.delay)
	} else if inpututil.IsKeyJustReleased(ebiten.KeyDown) {
		if r.wordsPerMinute > step {
			r.wordsPerMinute -= step
			ui.GetPreferences().Set("default words per minute speed", r.wordsPerMinute)
			r.delay = r.SetDelay(r.wordsPerMinute)
			r.rrLabel.SetWordsPerMinute(strconv.Itoa(r.wordsPerMinute))
		}
		ui.GetUi().ShowNotification(fmt.Sprintf("Speed down words per minute:%v", r.wordsPerMinute))
		log.Printf("now words per minute:%v delay is:%vms", r.wordsPerMinute, r.delay)
	}
}

func (r *RapidReadScene) Draw(surface *ebiten.Image) {
	for _, v := range r.container {
		v.Draw(surface)
	}
}

func (r *RapidReadScene) Entered() {
	r.wordsPerMinute = ui.GetPreferences().Get("default words per minute speed").(int)
	r.delay = r.SetDelay(r.wordsPerMinute)
	r.Resize()
}

func (r *RapidReadScene) Add(value ui.Drawable) {
	r.container = append(r.container, value)
}

func (r *RapidReadScene) Resize() {
	w, h := ui.GetUi().GetScreenSize()
	r.rect = ui.NewRect([]int{0, 0, w, h})
	x, y, w, h := 0, 0, int(float64(r.rect.GetLowestSize())*0.05), int(float64(r.rect.GetLowestSize())*0.05)
	r.btnQuit.Resize([]int{x, y, w, h})
	w, h = int(float64(r.rect.W)*0.95), int(float64(r.rect.H)*0.15)
	x, y = (r.rect.W-w)/2, (r.rect.H-h)/2
	r.rrLabel.Resize([]int{x, y, w, h})
	r.rrLabel.SetFont(h / 2)
}

func (r *RapidReadScene) Quit() {
	log.Printf("Quit reading at idxA:%v, idxB:%v with speed:%v", r.paragraphs.current, r.paragraph.current, r.wordsPerMinute)
	r.book.Update(r.paragraphs.current, r.paragraph.current, r.wordsPerMinute, r.book.status)
	GetDb().UpdateBook(r.book)
	for _, v := range r.container {
		v.Close()
	}
}
