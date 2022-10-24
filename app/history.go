package app

import (
	"fmt"
	"path/filepath"
)

var historyInstance *History = nil

func init() {
	historyInstance = GetHistory()
}

func GetHistory() (h *History) {
	if historyInstance == nil {
		h = newHistory()
	} else {
		h = historyInstance
	}
	return h
}

type History struct {
	count  int
	values map[int]*Txt
}

func newHistory() *History { return &History{values: map[int]*Txt{}} }

func (h *History) LoadBookByFilename(value string) *book {
	for _, v := range h.values {
		if v.name == value {
			return v.book
		}
	}
	return nil
}

func (h *History) Setup() {
	h.ReadTextsFromWorkDir()
}

func (h *History) ReadTextsFromWorkDir() {
	files, err := filepath.Glob("texts/*.txt")
	if err != nil {
		panic(err)
	}
	for _, v := range files {
		t := newTxt()
		t.filename = v
		t.Setup()
		fmt.Println(":>", v, t.ShortString())
	}
}

func (h *History) Add(value *Txt) {
	h.count += 1
	h.values[h.count] = value
}

func (h *History) GetList() (strs []string) {
	for _, v := range h.values {
		strs = append(strs, v.ShortString())
	}
	return
}

func (h *History) String() string {
	s := ""
	for k, v := range h.values {
		s += fmt.Sprintf("%v: %v\n", k, v)
	}
	return s
}
