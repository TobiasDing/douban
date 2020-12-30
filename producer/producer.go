package producer

import (
	"douban/spider"
	"douban/wait"
	"sync"
)
var tagChan chan string
func Producing() {
	var wg sync.WaitGroup
	tagChan := spider.GetTags("http://book.douban.com")
	for i := 0; i < 2; i = i + 1 {
		wg.Add(1)
		go spider.Spider("http://book.douban.com", wg, tagChan)
	}


	wg.Wait()


	wait.Wg.Done()
}
