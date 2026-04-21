package main

import (
	"fmt"
	"os"
	"unicode"
)

func main() {
	js_code := read_file()
	fmt.Println("From JavaScript:\n" + js_code)
	tokens := lexer(js_code)
	fmt.Println(tokens)
}

func read_file() string {
	file, err := os.ReadFile("case1.js")
	if err != nil {
		panic(err)
	}
	return string(file)
}

const (
	StateStart        = "StateStart"
	StateInIdentifier = "StateInIdentifier"
	StateInInteger    = "StateInInteger"
	StateInString     = "StateInString"
	StateInOperator   = "StateInOperator"
)

type TokenType int

const (
	TokenIdentifier = iota
	TokenString
	TokenInteger
	TokenEqual
)

type Token struct {
	Type  TokenType
	Value string
}

func lexer(js_code string) []Token {
	tokens := []Token{}

	state := StateStart
	start_index := 0

	for current_index := 0; current_index < len(js_code); current_index++ {
		current_char := rune(js_code[current_index])

		switch state {
		case StateStart:
			{
				start_index = current_index

				if unicode.IsLetter(current_char) {
					state = StateInIdentifier
				} else if current_char == '"' {
					state = StateInString
				} else if unicode.IsDigit(current_char) {
					state = StateInInteger
				} else if current_char == ';' || current_char == ' ' || current_char == '\n' || current_char == '\r' || current_char == '\t' {
				} else if current_char == '=' {
					state = StateInOperator
				} else {
					panic(fmt.Sprintf("Lexer Error: Unexpected character '%s' for state '%s'", string(current_char), state))
				}
			}
		case StateInIdentifier:
			{
				if !(unicode.IsLetter(current_char) || unicode.IsDigit(current_char) || current_char == '_') {
					tokens = append(tokens, Token{Value: js_code[start_index:current_index]})
					state = StateStart
				}
			}
		case StateInString:
			{
				if current_char == '"' {
					tokens = append(tokens, Token{Type: TokenString, Value: js_code[start_index+1 : current_index]})
					state = StateStart
				}
			}
		case StateInInteger:
			{
				if !unicode.IsDigit(current_char) {
					tokens = append(tokens, Token{Type: TokenInteger, Value: js_code[start_index:current_index]})
					state = StateStart
				}
			}
		case StateInOperator:
			{
				if current_char != '=' {
					tokens = append(tokens, Token{Type: TokenEqual, Value: "Equal"})
					state = StateStart
				}
			}
		default:
			{
				panic(fmt.Sprintf("Lexer Error: Unexpected character '%s' for state '%s'", string(current_char), state))
			}
		}
	}

	return tokens
}
