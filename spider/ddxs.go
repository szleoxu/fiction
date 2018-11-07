package spider

import (
	"fmt"
	"fiction/common"
	"fiction/models"
	"github.com/PuerkitoBio/goquery"
	"database/sql"
	"fiction/config"
	"time"
	"os"
)

const(
	bookListSelector="#hotcontent > div.l > div"
	chapterListSelector="#list > dl > dt:nth-child(10)"
	chapterTitleSelector="#wrapper > div.content_read > div > div.bookname > h1"
	chapterContentSelector="#content"
	chapterPreSelector="#wrapper > div.content_read > div > div.bottem2 > a:nth-child(2)"
	chapterNextSelector="#wrapper > div.content_read > div > div.bottem2 > a:nth-child(4)"
)

func DDXS(){
	start := time.Now()
	fmt.Println("spider is running")
	//get book list
	urlDDXS:=config.GetValue("spiderUrl","ddxs")
	dom:=common.UrlResponse(urlDDXS)
	if dom!=nil{
		dbBook:=models.InitBook()
		dbChapter:=models.InitChapter()
		dom.Find(bookListSelector).Each(func(i int, content *goquery.Selection) {
			name := common.GbkToUtf8(content.Find("dl > dt > a").Text())
			fmt.Println("Find book:"+name)
			isExist:=dbBook.IsExistBook(name)
			if isExist{
				fmt.Println("Exist book:"+name)
			}else{
				author := common.GbkToUtf8(content.Find("dl > dt > span").Text())
				image,_ := content.Find("div.image > a > img").Attr("src")
				chapterListUrl,_:= content.Find("div.image > a").Attr("href")
				tbBook:=models.BookTB{}
				var bookName=sql.NullString{name,true}
				tbBook.Name=bookName
				tbBook.Author=author
				tbBook.Image=image
				tbBook.SiteUrl=chapterListUrl
				insertLastID:=dbBook.Insert(tbBook)
				if insertLastID>0{
					//get chapter list
					dom=common.UrlResponse(chapterListUrl)
					var chapterNum=0
					dom.Find(chapterListSelector).NextAll().EachWithBreak(func(i int, content *goquery.Selection) bool{
						//get chapter
						chapterUrl,_:= content.Find("a").Attr("href")
						chapterUrl=chapterListUrl+chapterUrl
						dom=common.UrlResponse(chapterUrl)
						if dom!=nil{
							chapterNum++
							title:=common.GbkToUtf8(dom.Find(chapterTitleSelector).Text())
							isExist=dbChapter.IsExistChapter(insertLastID,title)
							if isExist{
								fmt.Println("Exist chapter:"+title)
							}else{
								chapterContent:=dom.Find(chapterContentSelector).Text()
								chapterContent=common.GbkToUtf8(chapterContent)
								pre,_:=dom.Find(chapterPreSelector).Attr("href")
								pre=chapterListUrl+pre
								next,_:=dom.Find(chapterNextSelector).Attr("href")
								next=chapterListUrl+next
								tbChapter:=models.ChapterTB{}
								tbChapter.Bid=insertLastID
								var chapterTitle=sql.NullString{title,true}
								tbChapter.Title=chapterTitle
								tbChapter.Content=chapterContent
								tbChapter.Sort=chapterNum
								tbChapter.Pre=pre
								tbChapter.Next=next
								isInsert:=dbChapter.Insert(tbChapter)
								if isInsert==true{
									fmt.Println("Insert chapter success")
								}else{
									fmt.Println("Insert chapter fail")
								}
							}
						}
						return true
					})
				}
			}
		})
		dbBook.DB.Close()
		dbChapter.DB.Close()
		fmt.Println("complete")
	}
	cost := time.Since(start)
	fmt.Printf("time=[%s]",cost)
	os.Exit(0)
}

func AfreshBook(bookID int64){
	start := time.Now()
	fmt.Println("spider is running")
	dbBook:=models.InitBook()
	dbBook.GetBook(bookID)
	fmt.Println(dbBook.Book.Name.String)
	dbChapter:=models.InitChapter()
	ret:=dbChapter.DeleteByBid(bookID)
	if ret{
		chapterListUrl:=dbBook.Book.SiteUrl
		//get chapter list
		dom:=common.UrlResponse(dbBook.Book.SiteUrl)
		var chapterNum=0
		dom.Find(chapterListSelector).NextAll().EachWithBreak(func(i int, content *goquery.Selection) bool{
			//get chapter
			chapterUrl,_:= content.Find("a").Attr("href")
			chapterUrl=chapterListUrl+chapterUrl
			dom=common.UrlResponse(chapterUrl)
			if dom!=nil{
				title:=common.GbkToUtf8(dom.Find(chapterTitleSelector).Text())
				isExist:=dbChapter.IsExistChapter(bookID,title)
				if isExist{
					fmt.Println("Exist chapter:"+title)
				}else{
					chapterContent:=dom.Find(chapterContentSelector).Text()
					chapterContent=common.GbkToUtf8(chapterContent)
					pre,_:=dom.Find(chapterPreSelector).Attr("href")
					pre=chapterListUrl+pre
					next,_:=dom.Find(chapterNextSelector).Attr("href")
					next=chapterListUrl+next
					tbChapter:=models.ChapterTB{}
					tbChapter.Bid=bookID
					var chapterTitle=sql.NullString{title,true}
					tbChapter.Title=chapterTitle
					tbChapter.Content=chapterContent
					chapterNum++
					tbChapter.Sort=chapterNum
					tbChapter.Pre=pre
					tbChapter.Next=next
					isInsert:=dbChapter.Insert(tbChapter)
					if isInsert==true{
						fmt.Println("Insert chapter success")
					}else{
						fmt.Println("Insert chapter fail")
					}
				}
			}
			return true
		})
	}
	dbBook.DB.Close()
	dbChapter.DB.Close()
	fmt.Println("complete")
	cost := time.Since(start)
	fmt.Printf("time=[%s]",cost)
	os.Exit(0)
}

func GetLastChapter(bookID int64){
	dbBook:=models.InitBook()
	dbChapter:=models.InitChapter()
	dbBook.GetBook(bookID)
	chapterListUrl:=dbBook.Book.SiteUrl
	dbChapter.GetLastSort(bookID)
	chapterNum:=dbChapter.Chapter.Sort
	dom:=common.UrlResponse(chapterListUrl)
	dom.Find(chapterListSelector).PrevAll().EachWithBreak(func(i int, content *goquery.Selection) bool{
		//get chapter
		chapterUrl,_:= content.Find("a").Attr("href")
		if chapterUrl!=""{
			chapterUrl=chapterListUrl+chapterUrl
			dom=common.UrlResponse(chapterUrl)
			if dom!=nil{
				title:=common.GbkToUtf8(dom.Find(chapterTitleSelector).Text())
				isExist:=dbChapter.IsExistChapter(bookID,title)
				if isExist{
					fmt.Println("Exist chapter:"+title)
				}else{
					chapterContent:=dom.Find(chapterContentSelector).Text()
					chapterContent=common.GbkToUtf8(chapterContent)
					pre,_:=dom.Find(chapterPreSelector).Attr("href")
					pre=chapterListUrl+pre
					next,_:=dom.Find(chapterNextSelector).Attr("href")
					next=chapterListUrl+next
					tbChapter:=models.ChapterTB{}
					tbChapter.Bid=bookID
					var chapterTitle=sql.NullString{title,true}
					tbChapter.Title=chapterTitle
					tbChapter.Content=chapterContent
					chapterNum++
					tbChapter.Sort=chapterNum
					tbChapter.Pre=pre
					tbChapter.Next=next
					isInsert:=dbChapter.Insert(tbChapter)
					if isInsert==true{
						fmt.Println("Insert chapter success:"+title)
					}else{
						fmt.Println("Insert chapter fail")
					}
				}
			}
		}
		return true
	})
	dbBook.DB.Close()
	dbChapter.DB.Close()
	fmt.Println("complete")
}