package models

import (
	"database/sql"
	"fmt"
	"time"
	"strconv"
)

type DBChapter struct {
	DB      *sql.DB
	Chapter ChapterTB
}

type ChapterTB struct {
	Id   int
	Bid  int64
	Title sql.NullString
	Content string
	Sort int
	Pre string
	Next string
}

func InitChapter()(*DBChapter){
	dbw := DBChapter{}
	dbw.DB=Init()
	return &dbw
}

func (dbw *DBChapter) QueryDataPre() {
	dbw.Chapter = ChapterTB{}
}

func (dbw *DBChapter)IsExistChapter(bookID int64,title string)(bool) {
	dbw.QueryDataPre()
	err := dbw.DB.QueryRow("select title from chapter where bid="+strconv.FormatInt(bookID,10)+" and title='"+title+"'").Scan(&dbw.Chapter.Title)
	if err!=nil{
		return false
	}
	if dbw.Chapter.Title.Valid {
		return true
	} else {
		return false
	}
}


func (dbw *DBChapter) Insert(tb ChapterTB) (bool) {
	stmt, _ := dbw.DB.Prepare("INSERT INTO chapter (bid, title,content,sort,pre,next,created_at) VALUES (?,?,?,?,?,?,?)")
	defer stmt.Close()
	_, err := stmt.Exec(tb.Bid, tb.Title.String, tb.Content, tb.Sort,tb.Pre, tb.Next, time.Now().Unix())
	if err != nil {
		fmt.Printf("insert data error: %v\n", err)
		return false
	} else {
		return true
	}
}

func (dbw *DBChapter) DeleteByBid(bid int64) bool{
	stmt, _ := dbw.DB.Prepare("DELETE FROM chapter WHERE bid = ?")
	defer stmt.Close()
	_, err := stmt.Exec(bid)
	if err != nil {
		fmt.Printf("insert data error: %v\n", err)
		return false
	} else {
		return true
	}
}

