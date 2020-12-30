package spider

import (
	"douban/fetcher"
	"douban/kafka"
	"douban/parser"
	"encoding/json"
	"fmt"
	"sync"
)


func Spider(url string, wg sync.WaitGroup, tagChan chan string) {

	err := kafka.Init([]string{"127.0.0.1:9092"}, 10000)
	if err != nil {
		fmt.Println("kafka initial failed, err:", err)
		return
	}


	tag := <-tagChan

	contents, err := fetcher.Fetch(fmt.Sprintf("%s%s", url, tag))
	if err != nil {
		panic(err)
	}
	list := parser.ParseList(contents)
	for _, book := range list {
		contents, err := fetcher.Fetch(book)
		if err != nil {
			panic(err)
		}
		detail := parser.ParseDetail(contents)

		bookJson, err := ConvertToJson(detail)
		if err != nil {
			panic(err)
		}

		kafka.SendToChan("douban_book", string(bookJson))
		//time.Sleep(time.Millisecond*200)
	}

	wg.Done()

}
func GetTags(url string) chan string {


	bytes, err := fetcher.Fetch(url)

	if err != nil {
		panic(err)
	}
	tags := parser.ParseTag(bytes)
	//fmt.Println(tags)
	tagChan := make(chan string, 100)
	for _, tag := range tags {

		tagChan <- tag

	}
	return tagChan
}

func ConvertToJson(detail interface{}) (m []byte, err error) {
	m, err = json.Marshal(detail)
	if err != err {
		fmt.Println("Marshal failed, err:", err)
		return
	}
	return
}
