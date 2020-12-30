package parser

import (
	"douban/model"
	"regexp"
	"strconv"
)
var nameRe = regexp.MustCompile(`<span property="v:itemreviewed">([^<]+)</span>`)
var authorRe = regexp.MustCompile(`<span class="pl"> 作者</span>:[\d\D]*?<a.*?>([^<]+)</a>`)
var publisherRe = regexp.MustCompile(`<span class="pl">出版社:</span> ([^<]+)<br/>`)
var numOfPageRe = regexp.MustCompile(`<span class="pl">页数:</span> (\d+)<br/>`)
var priceRe = regexp.MustCompile(`<span class="pl">定价:</span> ([^<]+)元<br/>`)
var scoreRe = regexp.MustCompile(`<strong class="ll rating_num " property="v:average"> (\d\.\d) </strong>`)
var infoRe = regexp.MustCompile(`<div class="intro">[\d\D]*?<p>([^<]+)</p>`)

func ParseDetail(contents []byte) (detail model.BookDetail) {
	detail = model.BookDetail{}
	detail.BookName = ExtraString(contents, nameRe)
	detail.Author = ExtraString(contents, authorRe)
	num, err := strconv.Atoi(ExtraString(contents, numOfPageRe))
	if err == nil {
		detail.NumOfPages = num
	} else {
		detail.NumOfPages = 0
	}
	detail.Publisher = ExtraString(contents, publisherRe)
	score, err := strconv.ParseFloat(ExtraString(contents, scoreRe), 2)
	if err == nil {
		detail.Score = float32(score)
	} else {
		detail.Score = 0.0
	}
	detail.Info = ExtraString(contents, infoRe)

	price, err := strconv.ParseFloat(ExtraString(contents, priceRe), 2)

	if err == nil {
		detail.Price = float32(price)
	} else {
		detail.Price = 0.0
	}

	return detail
}

func ExtraString(contents []byte, re *regexp.Regexp) string {
	match := re.FindSubmatch(contents)
	if len(match) >= 2 {
		return string(match[1])
	} else {
		return ""
	}
}