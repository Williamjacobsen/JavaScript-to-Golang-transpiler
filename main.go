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
	TokenDot
	TokenLeftParen
	TokenRightParen
)

type Token struct {
	Type  TokenType
	Value string
}

func is_supported_symbol(character rune) bool {
	return character == '=' || character == '.' || character == '(' || character == ')'
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
				} else if is_supported_symbol(current_char) {
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
					current_index-- // 'console.log', without this the '.' would be forgotten
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
					current_index-- // '10+5' <-- notice that the '+' is right after, it should also be consumed.
				}
			}
		case StateInOperator:
			{
				if !is_supported_symbol(current_char) {
					switch string(js_code[start_index:current_index]) {
					case "=":
						{
							tokens = append(tokens, Token{Type: TokenEqual, Value: "Equal"})
						}
					case ".":
						{
							tokens = append(tokens, Token{Type: TokenDot, Value: "Dot"})
						}
					case "(":
						{
							tokens = append(tokens, Token{Type: TokenLeftParen, Value: "LeftParen"})
						}
					case ")":
						{
							tokens = append(tokens, Token{Type: TokenRightParen, Value: "RightParen"})
						}
					default:
						{
							panic(fmt.Sprintf("Lexer Error: Operator '%s' is not supported.", js_code[start_index:current_index]))
						}
					}
					state = StateStart
					current_index--
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
