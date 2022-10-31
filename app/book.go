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
	size               int
	idx                int // index in book, int paragraph
	data               *paragraph
	status             readStatus
}

func newBook() *Book {
	speed := ui.GetPreferences().Get("default words per minute speed").(int)
	return &Book{
		count:     1,
		lastSpeed: speed,
		idx:       0,
		status:    start,
	}
}

func (t *Book) Setup() {
	info, err := os.Stat(t.filename)
	if err != nil {
		panic(err)
	}
	t.dt = info.ModTime().Format("2006-01-02 15:04:05")
	var tmp1 string
	t.data, tmp1 = loadBook(t.filename)
	tmp := []rune(tmp1)
	if len(tmp) >= 80 {
		tmp = tmp[0:50]
	}
	t.name = string(tmp)
	t.size = t.data.Size()
	log.Println("Setup book:", len(t.data.data), t.name)
}

func (t *Book) Update(idx, lastspeed int, status readStatus) {
	t.count += 1
	if idx > 0 {
		idx -= 1
	}
	t.idx = idx
	t.lastSpeed = lastspeed
	t.status = status
}

func (t *Book) GetName() string {
	return t.name
}

func (t *Book) GetFileName() string {
	return t.filename
}

func (t *Book) String() string {
	s := fmt.Sprintf("Book:%v, read %v times, last read %v times, at speed:%v, from:%v, size:%v status:%v idx:%v", t.name, t.count, t.dt, t.lastSpeed, t.filename, t.size, t.status, t.idx)
	return s
}
