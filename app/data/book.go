package data

import (
	"fmt"
	"log"
)

type readStatus int

const (
	Start readStatus = iota
	InReading
	Finished
)

type Book struct {
	dt, name         string //create time or last access time, filename in drive, name in list
	count, lastSpeed int    // read count, last read speed
	size             int
	idx              int // index in book, int paragraph
	data             *paragraph
	status           readStatus
	content          []byte
}

func newBook() *Book {
	speed := 300
	return &Book{
		count:     1,
		lastSpeed: speed,
		idx:       0,
		status:    Start,
	}
}

func (t *Book) Setup() {
	t.data, t.name = loadBook(t.content)
	t.size = t.data.Size()
	log.Println("Setup book:", len(t.data.data), t.name)
}

func (t *Book) GetName() string {
	return t.name
}

func (t *Book) GetStatus() readStatus {
	return t.status
}

func (t *Book) SetStatus(value readStatus) {
	t.status = value
}

func (t *Book) GetLastSpeed() int {
	return t.lastSpeed
}

func (t *Book) GetIndex() int {
	return t.idx
}

func (t *Book) SetIndex(value int) {
	t.idx = value
}

func (t *Book) GetParagraph() *paragraph {
	return t.data
}

func (t *Book) String() string {
	s := fmt.Sprintf("Book:%v, read %v times, last read %v times, at speed:%v, size:%v status:%v idx:%v", t.name, t.count, t.dt, t.lastSpeed, t.size, t.status, t.idx)
	return s
}
