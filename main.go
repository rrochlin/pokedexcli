package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println("Hello, World!")
}

func cleanInput(text string) []string {
	processed := strings.ToLower(text)
	substrings := strings.Split(processed, " ")
	res := make([]string, 0)
	for _, str := range substrings {
		if str == "" {
			continue
		}
		res = append(res, str)
	}
	return res
}
