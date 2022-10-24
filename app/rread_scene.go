package app

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/t0l1k/sread/ui"
)

type RapidReadScene struct {
	rect           *ui.Rect
	container      []ui.Drawable
	btnQuit        *ui.Button
	rrLabel        *RRLabel
	paragraph      *paragraph
	inGame         bool
	delay          int
	dt             int
	wordsPerMinute int
	book           *book
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
	return int(time.Duration(1.0/float64(delay)*float64(time.Minute)) / 1e6)
}

func (r *RapidReadScene) LoadBookFromHistory(name string) {
	r.book = GetHistory().LoadBookByFilename(name)
	r.book.NextParagraph()
	r.paragraph = newParagraph(r.book.Value())
	r.getNextWord()
	r.inGame = true
	log.Println("Loaded text from clipboard.")
}

func (r *RapidReadScene) LoadBookFromClipboard() {
	r.book = LoadBookAndSaveFromClipboard()
	r.book.NextParagraph()
	r.paragraph = newParagraph(r.book.Value())
	r.getNextWord()
	r.inGame = true
	log.Println("Loaded text from clipboard.")
}

// func (r *RapidReadScene) LoadBookFrom(filename string, paragraph, word int) {
// 	r.loadTextFile(filename)
// 	r.book.SetParagraph(paragraph)
// 	r.paragraph = newParagraph(r.book.Value())
// 	r.paragraph.SetWord(word)
// 	r.getNextWord()
// 	r.inGame = true
// }

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
		if r.book.NextParagraph() {
			par := r.book.Value()
			r.paragraph = newParagraph(par)
		} else {
			r.inGame = false
		}
	}
}

func (r *RapidReadScene) checkKeypress() {
	if inpututil.IsKeyJustReleased(ebiten.KeySpace) {
		r.inGame = !r.inGame
		ui.GetUi().ShowNotification("toggle pause")
	} else if inpututil.IsKeyJustReleased(ebiten.KeyUp) {
		if r.wordsPerMinute < 20000 {
			r.wordsPerMinute += 50
			ui.GetPreferences().Set("default words per minute speed", r.wordsPerMinute)
			r.delay = r.SetDelay(r.wordsPerMinute)
			r.rrLabel.SetWordsPerMinute(strconv.Itoa(r.wordsPerMinute))
		}
		ui.GetUi().ShowNotification(fmt.Sprintf("Speed up words per minute:%v", r.wordsPerMinute))
		log.Printf("now words per minute:%v delay is:%vms", r.wordsPerMinute, r.delay)
	} else if inpututil.IsKeyJustReleased(ebiten.KeyDown) {
		if r.wordsPerMinute > 50 {
			r.wordsPerMinute -= 50
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
	fmt.Println("Quit at index:", r.book.current, r.paragraph.current)
	for _, v := range r.container {
		v.Close()
	}
}
