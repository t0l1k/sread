package data

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"unicode"

	"golang.design/x/clipboard"
)

func LoadBookFromClipboardAndSave() *Book {
	t := newBook()
	t.content = saveTextFromClipboard()
	t.Setup()
	GetDb().InsertBook(t)
	return t
}

// Save Text From Clipboard in texts dir with filename generated by uuid
func saveTextFromClipboard() []byte {
	// init clipboard
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}
	// got text
	data := clipboard.Read(clipboard.FmtText)
	return data
}

func LoadBookByFilename(name string) *Book {
	books := GetDb().GetFromDbHistory()
	for _, v := range books.books {
		if v.name == name {
			log.Printf("Found:%v", v)
			return v
		}
	}
	return nil
}

func loadBook(content []byte) (*paragraph, string) {
	rfile := bytes.NewReader(content)

	book := newParagraph()
	ln := 30 //max word lenght
	var (
		w string
		l int
	)
	r := bufio.NewReader(rfile)
	for {
		if c, _, err := r.ReadRune(); err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		} else {
			if checkSeparator(c) {
				if len(w) > 0 {
					book.Add(w)
					w = ""
					l = 0
				}
			} else if l > ln {
				w += string(c)
				book.Add(w)
				w = ""
				l = 0
			} else {
				if checkSeparator(c) {
					continue
				}
				w += string(c)
				l++
			}
		}
	}
	rfile.Seek(0, io.SeekStart)
	fscanner := bufio.NewScanner(rfile)
	fscanner.Split(bufio.ScanWords)
	var tmp string
	for fscanner.Scan() {
		tmp += fscanner.Text() + " "
		if len(tmp) > 120 {
			break
		}
	}
	tmp1 := []rune(tmp)
	if len(tmp1) >= 120 {
		tmp1 = tmp1[0:120]
	}

	return book, string(tmp1)
}

func checkSeparator(r rune) bool {
	if unicode.IsSpace(r) {
		return true
	} else if r == '-' {
		return true
	}
	return false
}
