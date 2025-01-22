package main

import (
	"fmt"
	"regexp"
)

func main() {
	input := "바보<ol><li>1</li><li>2</li><li>3</li></ol>야"
	re := regexp.MustCompile(`(?s)<ol>.*?</ol>`)
	result := re.FindString(input)
	fmt.Println(result)
}
