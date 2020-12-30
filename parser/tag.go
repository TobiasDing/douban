package parser

import (
	"regexp"
)

const regexpStr = `<a href="([^"]+)" class="tag">([^<]+)</a>`
func ParseTag(contents []byte) []string {
	var result []string
	re := regexp.MustCompile(regexpStr)
	matches := re.FindAllSubmatch(contents, -1)


	for _, m := range matches {

		result = append(result, string(m[1]))

	}
	//fmt.Println(result)
	return result
}
