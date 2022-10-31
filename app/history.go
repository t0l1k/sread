package app

import (
	"log"
)

type History struct {
	count int
	books map[int]*Book
}

func newHistory() *History {
	return &History{
		books: make(map[int]*Book),
	}
}

func (h *History) New(id int, value *Book) {
	h.count = id
	h.books[h.count] = value
	log.Println("Added book", h.books[h.count])
}
