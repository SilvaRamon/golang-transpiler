package tokenizer

import (
	"fmt"
	"golang-transpiler/pkg/types"
	"regexp"
)

type Token struct {
	Type  TokenType
	Value string
	LineNumber int
}

type TokenType int

const (
	NumberLiteral TokenType = iota
	StringLiteral
	Identifier
	OpenParen
	CloseParen
	Quote
)

func (t TokenType) String() string {
	return [...]string{
		"NumberLiteral",
		"StringLiteral",
		"Identifier",
		"OpenParen",
		"CloseParen",
		"Quote",
	}[t]
}

func shift(input *string) string {
	var firstChar string = string((*input)[0])
	*input = (*input)[1:]
	return firstChar
}

func isNumber(input string) bool {
	return regexp.MustCompile(`[0-9]`).MatchString(input)
}

func isAlpha(input string) bool {
	return regexp.MustCompile(`[a-zA-Z]`).MatchString(input)
}

func isWhitespace(input string) bool {
	return regexp.MustCompile(`\s`).MatchString(input)
}

func current(input string) string {
	return string(input[0])
}

func Tokenizer(source types.SourceCode) []Token {
	var input string = source.Line
	var line int = source.LineNumber
	var tokens []Token

	for len(input) > 0 {
		currentValue := current(input)
		if currentValue == "(" {
			tokens = append(tokens, Token{OpenParen, shift(&input), line})
		} else if currentValue == ")" {
			tokens = append(tokens, Token{CloseParen, shift(&input), line})

			for len(input) > 0 {
				if shift(&input) != "\n" || shift(&input) != "\t" || shift(&input) != " " {
					panic(fmt.Sprintf("Line %v: Unexpected token after %v", line, currentValue))
				}
			}
		} else if currentValue == " " || currentValue == "\n" || currentValue == "\t" {
			shift(&input)
		} else {
			if len(input) == 0 {
				break
			}

			if isNumber(currentValue) {
				var number string
				for len(input) > 0 && isNumber(current(input)) {
					number += shift(&input)
				}
				tokens = append(tokens, Token{NumberLiteral, string(number), line})
			} else if isAlpha(currentValue) {
				var identifier string
				for len(input) > 0 && isAlpha(current(input)) {
					identifier += shift(&input)
				}
				tokens = append(tokens, Token{Identifier, string(identifier), line})
			} else if currentValue == "\"" {
				shift(&input)
				var stringLiteral string
				for len(input) > 0 && current(input) != "\"" {
					stringLiteral += shift(&input)
				}
				shift(&input)
				tokens = append(tokens, Token{StringLiteral, string(stringLiteral), line})
			} else {
				panic(fmt.Sprintf("Line %v: Unexpected token: %v", line, currentValue))
			}
		}
	}
	return tokens
}