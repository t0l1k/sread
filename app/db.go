package app

import (
	"fmt"

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
	d.createBooksTable()
}

func (d *Db) createBooksTable() {
	var err error
	d.conn, err = sql.Open("sqlite3", "texts/books.db")
	if err != nil {
		panic(err)
	}
	var createGameDB string = "CREATE TABLE IF NOT EXISTS books(id INTEGER PRIMARY KEY AUTOINCREMENT,filename TEXT, name TEXT, dt TEXT, count INTEGER, idxa INTEGER, idxb INTEGER, lastspeed INTEGER, status INTEGER)"
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
	insStr := "INSERT INTO books(filename, name, dt, count, idxa, idxb, lastspeed, status) VALUES(?,?,?,?,?,?,?,?)"
	cur, err := d.conn.Prepare(insStr)
	if err != nil {
		log.Println("Error in DB:", insStr, values)
		panic(err)
	}
	defer cur.Close()
	filename := values.filename
	name := values.name
	dt := values.dt
	count := values.count
	idxa := values.idxA
	idxb := values.idxB
	lastspeed := values.lastSpeed
	status := values.status
	cur.Exec(filename, name, dt, count, idxa, idxb, lastspeed, status)
	log.Println("DB:Inserted:", values)
}

func (d *Db) UpdateBook(values *Book) {
	fmt.Println(values)
	updateStr := `UPDATE "books" SET count = ? , idxa = ? , idxb = ? , lastspeed = ? , status = ? WHERE filename = ?`
	res, err := d.conn.Exec(updateStr, values.count, values.idxA, values.idxB, values.
		lastSpeed, values.status, values.filename)
	if err != nil {
		fmt.Println(updateStr, err)
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
		err = rows.Scan(&id, &txt.filename, &txt.name, &txt.dt, &txt.count, &txt.idxA, &txt.idxB, &txt.lastSpeed, &txt.status)
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
		err = rows.Scan(&id, &book.filename, &book.name, &book.dt, &book.count, &book.idxA, &book.idxB, &book.lastSpeed, &book.status)
		if err != nil && err != sql.ErrNoRows {
			panic(err)
		}
		names = append(names, book.name)
	}
	return names
}
