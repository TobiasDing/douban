package main

import (
	"douban/producer"
	"douban/wait"
)


func main() {
	wait.Wg.Add(1)
	go producer.Producing()

	//time.Sleep(time.Second * 10)
	wait.Wg.Wait()

}
