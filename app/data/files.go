package data

import (
	"bufio"
	"bytes"
	"io"
	"time"
	"unicode"

	"golang.design/x/clipboard"
)

// Сохранить в БД из буфера обмена
func LoadBookFromClipboardAndSave() *Book {
	t := newBook()
	t.dt = time.Now().Format("2006-01-02 15:04:05.000")
	t.count = 1
	t.content = saveTextFromClipboard()
	t.Setup()
	GetDb().InsertBook(t)
	return t
}

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

// Прочитать из БД по имени
func LoadBookByFilename(name string) *Book {
	books := GetDb().GetFromDbHistory()
	for _, v := range books.books {
		if v.name == name {
			return v
		}
	}
	return nil
}

// Подготовка к чтению и генерация имени файла
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
