package data

import (
	"fmt"
)

type paragraph struct {
	current int // word in paragraph
	data    []string
}

func newParagraph() *paragraph {
	return &paragraph{current: -1}
}

func (p *paragraph) Add(value string) {
	p.data = append(p.data, value)
}

func (p *paragraph) SetWord(idx int) {
	p.current = idx - 1
}

func (p *paragraph) Value() string {
	return p.data[p.current]
}

func (p *paragraph) Index() int {
	return p.current
}

func (p *paragraph) IsFirstWorld() bool {
	return p.current == 0
}

func (p *paragraph) IsLastWorld() bool {
	return p.current >= len(p.data)
}

func (p *paragraph) NextWord() bool {
	if !p.IsLastWorld() {
		p.current++
	}
	return len(p.data) > p.current
}

func (p *paragraph) PrevWord() bool {
	if p.current >= 0 {
		p.current--
	}
	return p.current >= 0
}

func (p *paragraph) Size() int {
	return len(p.data)
}

func (p *paragraph) String() string {
	return fmt.Sprintf("%v %v", p.current, p.data[p.current])
}
