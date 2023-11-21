package data

import (
	"os"

	_ "github.com/mattn/go-sqlite3"

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
	// check dir present
	if _, err := os.Stat("texts"); os.IsNotExist(err) {
		err := os.Mkdir("texts", os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	d.createBooksTable()
}

func (d *Db) createBooksTable() {
	var err error
	d.conn, err = sql.Open("sqlite3", "texts/books.db")
	if err != nil {
		panic(err)
	}
	var createGameDB string = "CREATE TABLE IF NOT EXISTS books(id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT, dt TEXT, count INTEGER, idx INTEGER, lastspeed INTEGER, status INTEGER, content TEXT)"
	cur, err := d.conn.Prepare(createGameDB)
	if err != nil {
		panic(err)
	}
	cur.Exec()
	cur.Close()
	log.Println("Created table for books.")
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
	log.Println(values)
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
		panic(err)
	}
	defer rows.Close()
	history := newHistory()
	for rows.Next() {
		txt := newBook()
		id := 0
		err = rows.Scan(&id, &txt.name, &txt.dt, &txt.count, &txt.idx, &txt.lastSpeed, &txt.status, &txt.content)
		if err != nil && err != sql.ErrNoRows {
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
		panic(err)
	}
	defer rows.Close()
	var names []string
	for rows.Next() {
		book := newBook()
		id := 0
		err = rows.Scan(&id, &book.name, &book.dt, &book.count, &book.idx, &book.lastSpeed, &book.status, &book.content)
		if err != nil && err != sql.ErrNoRows {
			panic(err)
		}
		names = append(names, book.name)
	}
	return names
}
