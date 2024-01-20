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
	return &Book{}
}

func (t *Book) Setup() {
	t.data, t.name = loadBook(t.content)
	t.size = t.data.Size()
	t.count++
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

func (t *Book) SetLastSpeed(value int) {
	t.lastSpeed = value
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
	s := fmt.Sprintf("Книга:%v, время создания %v, прочитана %v раз, скорость чтения:%v, размер:%v статус:%v закончено на:%v", t.name, t.dt, t.count, t.lastSpeed, t.size, t.status, t.idx)
	return s
}
