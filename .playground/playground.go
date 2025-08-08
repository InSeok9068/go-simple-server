package main

import (
	"fmt"
	"strings"
)

func main() {
	font := []string{"Roboto", "Roboto+Condensed", "Roboto+Slab"}
	fmt.Println(strings.Join(font, "&family="))
}
