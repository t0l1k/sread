package app

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/t0l1k/sread/ui"
	"golang.design/x/clipboard"
)

type Txt struct {
	dt               time.Time // create time or last access time
	filename, name   string    // filename in drive, name in list
	count, lastSpeed int       // read count, last read speed
	size             int64
	book             *book
}

func newTxt() *Txt {
	speed := ui.GetPreferences().Get("default words per minute speed").(int)
	return &Txt{
		count:     1,
		lastSpeed: speed,
	}
}

func (t *Txt) Setup() {
	info, err := os.Stat(t.filename)
	if err != nil {
		panic(err)
	}
	t.dt = info.ModTime()
	t.size = info.Size()
	t.book = loadBook(t.filename)
	t.name = t.book.data[0]
	GetHistory().Add(t)
}

func LoadBookAndSaveFromClipboard() *book {
	t := newTxt()
	t.filename = saveTextFromClipboard()
	t.Setup()
	return t.book
}

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

func loadBook(filename string) *book {
	rfile, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer rfile.Close()
	fscanner := bufio.NewScanner(rfile)
	fscanner.Split(bufio.ScanLines)
	book := newBook()
	for fscanner.Scan() {
		book.Add(fscanner.Text())
	}
	return book
}

func (t *Txt) GetBook() *book {
	return t.book
}

func (t *Txt) GetName() string {
	return t.name
}

func (t *Txt) GetFileName() string {
	return t.filename
}

func (t *Txt) ShortString() string {
	s := fmt.Sprintf("%v", t.name)
	return s
}

func (t *Txt) String() string {
	s := fmt.Sprintf("Name:%v, read %v times, last read %v times, at speed:%v, from:%v, size:%v", t.name, t.count, t.dt.Format("2006-01-02 15:04:05.000"), t.lastSpeed, t.filename, t.size)
	return s
}
