package spider

import (
	"fmt"
	"os"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
)


func UrlResponse(url string) (*goquery.Document){
	response := GetResponse(url)
	if response.StatusCode == 403 {
		fmt.Println("403:IP 已被禁止访问")
		os.Exit(1)
	}
	if response.StatusCode == 200 {
		fmt.Println("url response is success ")
		dom, err := goquery.NewDocumentFromReader(response.Body)
		if err != nil {
			log.Fatalf("失败原因: %d", response.StatusCode)
		}
		return dom
	}else{
		return nil
	}
}


/**
* 返回response
*/
func GetResponse(url string) *http.Response {
	client := &http.Client{}
	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:50.0) Gecko/20100101 Firefox/50.0")
	response, _ := client.Do(request)
	return response
}
