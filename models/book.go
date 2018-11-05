package models

import (
	"database/sql"
	"fmt"
	"time"
)


type DBBook struct {
	Db       *sql.DB
	Book BookTB
}

type BookTB struct {
	Id   int
	Name sql.NullString
	Author string
	SiteName string
	SiteUrl string
	Image string
}

func InitBook()(*DBBook){
	dbw := DBBook{}
	dbw.Db=Init()
	return &dbw
}

func (dbw *DBBook)IsExistBook(name string)(bool) {
	dbw.QueryDataPre()
	err := dbw.Db.QueryRow("select name From book where name='"+name+"'").Scan(&dbw.Book.Name)
	if err!=nil{
		return false
	}
	if dbw.Book.Name.Valid {
		return true
	} else {
		return false
	}
}

func (dbw *DBBook) QueryDataPre() {
	dbw.Book = BookTB{}
}

func (dbw *DBBook) QueryData(name string)(bool) {
	dbw.QueryDataPre()
	rows, err := dbw.Db.Query("SELECT * From url where  name='{$name}'")
	defer rows.Close()
	if err != nil {
		fmt.Printf("query data error: %v\n", err)
		return false
	}
	for rows.Next() {
		rows.Scan(&dbw.Book.Id, &dbw.Book.Name, &dbw.Book.Author,&dbw.Book.Image)
		if err != nil {
			fmt.Printf(err.Error())
			continue
		}
	}
	return true
}

func (dbw *DBBook) InsertData(tb BookTB) (int64) {
	tb.SiteName="顶点小说"
	stmt, _ := dbw.Db.Prepare("INSERT INTO book (name, author,image,site_name,site_url,created_at) VALUES (?,?,?,?,?,?)")
	defer stmt.Close()
	ret, err := stmt.Exec(tb.Name.String, tb.Author, tb.Image, tb.SiteName, tb.SiteUrl, time.Now().Unix())
	if err != nil {
		fmt.Printf("insert data error: %v\n", err)
		return 0
	} else {
		if LastInsertId, err := ret.LastInsertId(); nil == err {
			return LastInsertId
		}
		return 0
	}
}

