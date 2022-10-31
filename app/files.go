package app

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"golang.design/x/clipboard"
)

func LoadBookAndSaveFromClipboard() *Book {
	t := newBook()
	t.filename = saveTextFromClipboard()
	t.Setup()
	GetDb().InsertBook(t)
	return t
}

// Save Text From Clipboard in texts dir with filename generated by uuid
func saveTextFromClipboard() string {
	// init clipboard
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}
	// got text
	data := clipboard.Read(clipboard.FmtText)
	// check dir present
	if _, err := os.Stat("texts"); os.IsNotExist(err) {
		err := os.Mkdir("texts", os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	// get filename
	id := uuid.New()
	filename := fmt.Sprintf("texts/%v.txt", id.String())
	// save text to file
	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		panic(err)
	}
	return filename
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

func loadBook(filename string) *paragraphs {
	rfile, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer rfile.Close()
	fscanner := bufio.NewScanner(rfile)
	fscanner.Split(bufio.ScanLines)
	book := newParagraphs()
	for fscanner.Scan() {
		w := fscanner.Text()
		book.Add(w)
	}
	return book
}
