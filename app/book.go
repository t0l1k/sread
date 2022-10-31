package app

import (
	"fmt"
	"log"
	"os"

	"github.com/t0l1k/sread/ui"
)

type readStatus int

const (
	start readStatus = iota
	inReading
	finished
)

type Book struct {
	dt, filename, name string //create time or last access time, filename in drive, name in list
	count, lastSpeed   int    // read count, last read speed
	size               int64
	idxA, idxB         int // index in book, int paragraph
	paragraps          *paragraphs
	status             readStatus
}

func newBook() *Book {
	speed := ui.GetPreferences().Get("default words per minute speed").(int)
	return &Book{
		count:     1,
		lastSpeed: speed,
		idxA:      0,
		idxB:      0,
		status:    start,
	}
}

func (t *Book) Setup() {
	info, err := os.Stat(t.filename)
	if err != nil {
		panic(err)
	}
	t.dt = info.ModTime().Format("2006-01-02 15:04:05")
	t.size = info.Size()
	t.paragraps = loadBook(t.filename)
	tmp := []rune(t.paragraps.data[0])
	if len(tmp) >= 80 {
		tmp = tmp[0:50]
	}
	t.name = string(tmp)
	log.Println("Setup book:", len(t.paragraps.data), t.name)
}

func (t *Book) Update(idxa, idxb, lastspeed int, status readStatus) {
	t.count += 1
	t.idxA = idxa
	if idxb > 0 {
		idxb -= 1
	}
	t.idxB = idxb
	t.lastSpeed = lastspeed
	t.status = status
}

func (t *Book) GetBook() *paragraphs {
	return t.paragraps
}

func (t *Book) GetName() string {
	return t.name
}

func (t *Book) GetFileName() string {
	return t.filename
}

func (t *Book) String() string {
	s := fmt.Sprintf("Book:%v, read %v times, last read %v times, at speed:%v, from:%v, size:%v status:%v", t.name, t.count, t.dt, t.lastSpeed, t.filename, t.size, t.status)
	return s
}
