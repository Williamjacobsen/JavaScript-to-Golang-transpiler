package main

import (
	"fmt"
	"os"
)

func main() {
	js_code := read_file()
	fmt.Println("From JavaScript:\n" + js_code)
	lexer(js_code)
}

func read_file() string {
	file, err := os.ReadFile("case1.js")
	if err != nil {
		panic(err)
	}
	return string(file)
}

const (
	StateInIdentifier = iota
	StateInInteger
	StateInString
)

func lexer(js_code string) {
	for i := 0; i < len(js_code); i++ {
		current_char := string(js_code[i])
	}
}
