package parser

import (
	"regexp"
)

const BookListRe = `<a href="([^"]+)" title="([^"]+)"`

func ParseList(contents []byte) (list []string) {

	re := regexp.MustCompile(BookListRe)

	matches := re.FindAllSubmatch(contents, -1)

	for _, m := range matches {
		list = append(list, string(m[1]))

	}

	return
}
