package app

import (
	"fmt"
	"strings"
)

type book struct {
	current int
	data    []string
}

func newBook() *book {
	b := &book{
		current: -1,
	}
	return b
}

func (b *book) SetParagraph(idx int) {
	b.current = idx - 1
}

func (b *book) Add(value string) {
	b.data = append(b.data, value)
}

func (b *book) Value() string {
	return b.data[b.current]
}

func (b *book) NextParagraph() bool {
	b.current++
	return len(b.data) > b.current

}

func (b *book) String() string {
	return fmt.Sprintf("%v %v", b.current, b.data[b.current])
}

type paragraph struct {
	current int
	data    []string
}

func newParagraph(value string) *paragraph {
	p := &paragraph{
		current: -1,
		data:    strings.Fields(value),
	}
	return p
}

func (p *paragraph) SetWord(idx int) {
	p.current = idx - 1
}

func (p *paragraph) Value() string {
	return p.data[p.current]
}

func (p *paragraph) NextWord() bool {
	p.current++
	return len(p.data) > p.current

}

func (p *paragraph) String() string {
	return fmt.Sprintf("%v %v", p.current, p.data[p.current])
}
