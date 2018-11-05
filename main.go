
package main

import (
	"strconv"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"fmt"
	"os"
	"fiction/spider"
	"fiction/models"
	"database/sql"
	"fiction/common"
)

const (
	baseUrl  = "https://studygolang.com/topics?p="
	fictionDDXS  ="https://www.booktxt.net/"
)

func main() {
	//goPage()
	//GoDDXS()
}

func GoDDXS(){
	fmt.Println("spider is running")
	//get book list
	var count =0
	dom:=spider.UrlResponse(fictionDDXS)
	if dom!=nil{
		dbBook,err:=models.InitBook()
		if err!=nil{
			log.Fatalf(err.Error())
		}
		dbChapter,err:=models.InitChapter()
		if err!=nil{
			log.Fatalf(err.Error())
		}
		dom.Find("#hotcontent > div.l > div").Each(func(i int, content *goquery.Selection) {
			count++
			name := common.GbkToUtf8(content.Find("dl > dt > a").Text())
			fmt.Println("find book:"+name)
			isExist:=dbBook.IsExistBook(name)
			if isExist{
				fmt.Println("exist book:"+name)
			}else{
				author := common.GbkToUtf8(content.Find("dl > dt > span").Text())
				image,_ := content.Find("div.image > a > img").Attr("src")
				fmt.Println("author:"+author+"__image:"+image)
				tbBook:=models.BookTB{}
				var bookName=sql.NullString{name,true}
				tbBook.Name=bookName
				tbBook.Author=author
				tbBook.Image=image
				tbBook.SiteUrl=fictionDDXS
				insertLastID:=dbBook.InsertData(tbBook)
				if insertLastID>0{
					//get chapter list
					chapterListUrl,_:= content.Find("div.image > a").Attr("href")
					dom=spider.UrlResponse(chapterListUrl)
					var chapterNum=0
					dom.Find("#list > dl > dt:nth-child(10)").NextAll().Each(func(i int, content *goquery.Selection) {
						//get chapter
						chapterUrl,_:= content.Find("a").Attr("href")
						chapterUrl=chapterListUrl+chapterUrl
						dom=spider.UrlResponse(chapterUrl)
						if dom!=nil{
							chapterNum++
							title:=common.GbkToUtf8(dom.Find("#wrapper > div.content_read > div > div.bookname > h1").Text())
							isExist=dbChapter.IsExistChapter(title)
							if isExist{
								fmt.Println("exist chapter:"+title)
							}else{
								chapterContent:=common.TrimHtml(dom.Find("#content").Text())
								chapterContent=common.GbkToUtf8(chapterContent)
								pre,_:=dom.Find("#wrapper > div.content_read > div > div.bottem2 > a:nth-child(2)").Attr("href")
								pre=chapterListUrl+pre
								next,_:=dom.Find("#wrapper > div.content_read > div > div.bottem2 > a:nth-child(4)").Attr("href")
								next=chapterListUrl+next
								tbChapter:=models.ChapterTB{}
								tbChapter.Bid=insertLastID
								var chapterTitle=sql.NullString{title,true}
								tbChapter.Title=chapterTitle
								tbChapter.Content=chapterContent
								tbChapter.Sort=chapterNum
								tbChapter.Pre=pre
								tbChapter.Next=next
								isInsert:=dbChapter.InsertData(tbChapter)
								if isInsert==true{
									fmt.Println("add chapter success")
								}else{
									fmt.Println("add chapter fail")
								}
							}
						}
					})
				}
			}

		})
		dbBook.Db.Close()
		dbChapter.Db.Close()
	}

	fmt.Println("--------------------------------数据拉取完成共"+strconv.Itoa(count)+"条------------------------------------")
	os.Exit(1)
}


func goPage(){
	var page int = 1
	//var count int =getPageCount()
	var count int =2

	for  {
		str := baseUrl + strconv.Itoa(page)
		response := getResponse(str)
		if response.StatusCode == 403 {
			fmt.Println("IP 已被禁止访问")
			os.Exit(1)
		}
		if response.StatusCode == 200 {
			dom, err := goquery.NewDocumentFromReader(response.Body)
			if err != nil {
				log.Fatalf("失败原因", response.StatusCode)
			}
			dom.Find(".topics .topic").Each(func(i int, content *goquery.Selection) {
				title := content.Find(".title a").Text()
				fmt.Println(title)
			})
		}
		page++
		if page >= count{
			fmt.Println("--------------------------------数据拉取完成共"+strconv.Itoa(page)+"条------------------------------------")
			os.Exit(1)
		}
	}
}

/**
* 返回response
*/
func getResponse(url string) *http.Response {
	client := &http.Client{}
	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:50.0) Gecko/20100101 Firefox/50.0")
	response, _ := client.Do(request)
	return response
}


