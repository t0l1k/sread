package data

import (
	"fmt"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/quasilyte/gdata"

	"database/sql"
	"log"
)

type Db struct {
	conn *sql.DB
}

var dbInstance *Db = nil

func init() {
	dbInstance = GetDb()
}

func GetDb() (db *Db) {
	if dbInstance == nil {
		db = &Db{}
	} else {
		db = dbInstance
	}
	return db
}

func (d *Db) Setup() {
	d.createBooksTable()
}

func (d *Db) createBooksTable() {
	var err error
	filename := d.getDbPath()
	d.conn, err = sql.Open("sqlite3", filename)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	var createGameDB string = "CREATE TABLE IF NOT EXISTS books(id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT, dt TEXT, count INTEGER, idx INTEGER, lastspeed INTEGER, status INTEGER, content TEXT)"
	cur, err := d.conn.Prepare(createGameDB)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	cur.Exec()
	cur.Close()
	log.Println("Created table for books.")
}

func (*Db) getDbPath() string {
	m, err := gdata.Open(gdata.Config{AppName: "sread"})
	if err != nil {
		log.Println(err)
		panic(err)
	}
	s := []byte("spritz reader")
	if err := m.SaveItem("sread.txt", s); err != nil {
		log.Println(err)
		panic(err)
	}
	filename := fmt.Sprintf("%v/books.db", strings.TrimSuffix(m.ItemPath("sread.txt"), "/sread.txt"))
	log.Println("to db path", strings.TrimSuffix(m.ItemPath("sread.txt"), "/sread.txt"), filename)
	return filename
}

func (d *Db) InsertBook(values *Book) {
	if d.conn == nil {
		d.Setup()
	}
	insStr := "INSERT INTO books(name, dt, count, idx, lastspeed, status, content) VALUES(?,?,?,?,?,?,?)"
	cur, err := d.conn.Prepare(insStr)
	if err != nil {
		log.Println("Error in DB:", insStr, values)
		panic(err)
	}
	defer cur.Close()
	name := values.name
	dt := values.dt
	count := values.count
	idx := values.idx
	lastspeed := values.lastSpeed
	status := values.status
	content := values.content
	cur.Exec(name, dt, count, idx, lastspeed, status, content)
	log.Println("DB:Inserted:", values)
}

func (d *Db) UpdateBook(values *Book) {
	log.Println("Update:", values)
	updateStr := `UPDATE "books" SET count = ? , idx = ? , lastspeed = ? , status = ? WHERE name = ?`
	res, err := d.conn.Exec(updateStr, values.count, values.idx, values.
		lastSpeed, values.status, values.name)
	if err != nil {
		log.Println(updateStr, err)
		panic(err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}
	log.Println("DB UPDATE AFFECTED:", count)
}

func (d *Db) GetFromDbHistory() *History {
	if d.conn == nil {
		d.Setup()
	}
	rows, err := d.conn.Query("SELECT * FROM books")
	if err != nil {
		log.Println(err)
		panic(err)
	}
	defer rows.Close()
	history := newHistory()
	for rows.Next() {
		txt := newBook()
		id := 0
		err = rows.Scan(&id, &txt.name, &txt.dt, &txt.count, &txt.idx, &txt.lastSpeed, &txt.status, &txt.content)
		if err != nil && err != sql.ErrNoRows {
			log.Println(err)
			panic(err)
		}
		history.New(id, txt)
	}
	log.Println("Done Read History from DB", history, len(history.books), "items")
	return history
}

func (d *Db) GetNames() []string {
	if d.conn == nil {
		d.Setup()
	}
	rows, err := d.conn.Query("SELECT * FROM books")
	if err != nil {
		log.Println(err)
		panic(err)
	}
	defer rows.Close()
	var names []string
	for rows.Next() {
		book := newBook()
		id := 0
		err = rows.Scan(&id, &book.name, &book.dt, &book.count, &book.idx, &book.lastSpeed, &book.status, &book.content)
		if err != nil && err != sql.ErrNoRows {
			log.Println(err)
			panic(err)
		}
		names = append(names, book.name)
	}
	return names
}

func (d *Db) DeleteRecords(names []string) {
	if d.conn == nil {
		d.Setup()
	}
	var (
		res sql.Result
		err error
	)
	deleteStr := `DELETE FROM "books" WHERE name = ?`
	for _, name := range names {
		res, err = d.conn.Exec(deleteStr, name)
		if err != nil {
			log.Println(deleteStr, err)
			panic(err)
		}
		log.Println("01-deleted from DB record:", name)
	}
	count, err := res.RowsAffected()
	if err != nil {
		log.Println(err)
		panic(err)
	}
	log.Println("DB DELETE AFFECTED:", count)
}
