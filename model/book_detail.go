package model

import "fmt"

type BookDetail struct {
	Id         int
	BookName   string
	Author     string
	Publisher  string
	NumOfPages int
	Price      float32
	Score      float32
	Info       string
}

func (b BookDetail) String() string {
	str := fmt.Sprintf("BookName: %s\nAuthor: %s\nPublisher: %s\nNum of Pages: %d\nPrice: %.2f\nScore: %.1f\nInfo: %s",
		b.BookName, b.Author, b.Publisher, b.NumOfPages, b.Price, b.Score, b.Info)
	return str
}
