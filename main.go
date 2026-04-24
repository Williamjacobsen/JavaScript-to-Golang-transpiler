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
	root_node := build_ast(tokens)
	if root_node == nil {
		panic("Main: Root node is nil.")
	}
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

// The Lexer should not handle build-in methods like console.log,
// they are just TokenIdentifier's, the parser will handle it.
// The Lexer only handles keywords, operators, etc...

type TokenType string

const (
	TokenIdentifier = "TokenIdentifier"

	TokenString  = "TokenString"
	TokenInteger = "TokenInteger"

	TokenEqual      = "TokenEqual"
	TokenDot        = "TokenDot"
	TokenLeftParen  = "TokenLeftParen"
	TokenRightParen = "TokenRightParen"

	TokenVar   = "TokenVar"
	TokenLet   = "TokenLet"
	TokenConst = "TokenConst"
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
					panic(fmt.Sprintf("Lexer Error: Unexpected character '%s' for state '%s'.", string(current_char), state))
				}
			}
		case StateInIdentifier:
			{
				if !(unicode.IsLetter(current_char) || unicode.IsDigit(current_char) || current_char == '_') {
					token := js_code[start_index:current_index]

					switch token {
					case "let":
						{
							tokens = append(tokens, Token{Type: TokenLet})
						}
					case "const":
						{
							tokens = append(tokens, Token{Type: TokenConst})
						}
					default:
						{
							tokens = append(tokens, Token{Type: TokenIdentifier, Value: token})
						}
					}

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
							tokens = append(tokens, Token{Type: TokenEqual})
						}
					case ".":
						{
							tokens = append(tokens, Token{Type: TokenDot})
						}
					case "(":
						{
							tokens = append(tokens, Token{Type: TokenLeftParen})
						}
					case ")":
						{
							tokens = append(tokens, Token{Type: TokenRightParen})
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
				panic(fmt.Sprintf("Lexer Error: Unexpected character '%s' for state '%s'.", string(current_char), state))
			}
		}
	}

	return tokens
}

type NodeType int

const (
	NodeProgram NodeType = iota
	NodeVariable
)

type Node any

type ProgramNode struct {
	Body []Node
}

type VariableNode struct {
	Name     string
	Operator string
	Value    Node
}

type ConsoleMethod int

const (
	ConsoleLog ConsoleMethod = iota
	ConsoleWarn
	ConsoleError
	ConsoleInfo
	ConsoleDebug
)

type ConsoleCallNode struct {
	Method ConsoleMethod
	Args   []Node
}

type Parser struct {
	tokens []Token
	index  int
}

func NewParser(tokens []Token) *Parser {
	return &Parser{tokens: tokens, index: 0}
}

func (p *Parser) current() Token {
	return p.tokens[p.index]
}

func (p *Parser) advance() {
	p.index++
}

func (p *Parser) consume_expect(expected_token TokenType) Token {
	if p.current().Type != expected_token {
		panic(fmt.Sprintf("Parser - consume_expect: Expected token '%s', got token '%s'", expected_token, p.current().Type))
	}
	p.advance()
	return p.tokens[p.index-1]
}

func (p *Parser) peek() Token {
	return p.tokens[p.index+1]
}

func (p *Parser) parse_variable() Node {
	p.advance()
	node_variable := VariableNode{}

	node_variable.Name = p.consume_expect(TokenIdentifier).Value
	node_variable.Operator = string(p.consume_expect(TokenEqual).Type)
	node_variable.Value = p.consume_expect(TokenInteger).Value

	fmt.Printf("Variable:\n\tName: %s\n\tOperator: %s\n\tValue: %s\n", node_variable.Name, node_variable.Operator, node_variable.Value)

	return node_variable
}

func (p *Parser) parse_program() Node {
	program_node := ProgramNode{}

	for p.index < len(p.tokens) {
		switch p.current().Type {
		case TokenVar, TokenLet:
			{
				program_node.Body = append(program_node.Body, p.parse_variable())
			}
		default:
			{
				panic(fmt.Sprintf("Parser: Case not handled for %s", p.current().Type))
			}
		}
	}

	return program_node
}

func build_ast(tokens []Token) Node {
	parser := NewParser(tokens)
	program_node := parser.parse_program()
	return program_node
}
