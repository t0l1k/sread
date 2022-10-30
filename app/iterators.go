package app

import (
	"fmt"
	"strings"
	"unicode"
)

type paragraphs struct {
	current int // current paragraph
	data    []string
}

func newParagraphs() *paragraphs {
	b := &paragraphs{
		current: -1,
	}
	return b
}

func (b *paragraphs) Set(idx int) {
	b.current = idx - 1
}

func (b *paragraphs) Add(value string) {
	if len(split(value)) == 0 {
		return
	}
	b.data = append(b.data, value)
}

func (b *paragraphs) Value() string {
	return b.data[b.current]
}

func (b *paragraphs) NextParagraph() bool {
	b.current++
	return len(b.data) > b.current

}

func (b *paragraphs) String() string {
	return fmt.Sprintf("%v:%v", b.current, b.data[b.current])
}

type paragraph struct {
	current int // word in paragraph
	data    []string
}

func newParagraph(value string) *paragraph {
	p := &paragraph{
		current: -1,
		data:    split(value),
	}
	return p
}

func (p *paragraph) SetWord(idx int) {
	if idx > len(p.data) {
		p.current = len(p.data) - 1
	} else {
		p.current = idx - 1
	}
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

func split(value string) []string {
	return strings.FieldsFunc(value, func(r rune) bool {
		if unicode.IsSpace(r) {
			return true
		} else if r == '-' {
			return true
		}
		return false
	})
}
