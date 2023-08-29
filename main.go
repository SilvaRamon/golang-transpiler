package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

type SourceCode struct {
	Line       string
	LineNumber int
}

type Program struct {
	Type string
	Body interface{}
}

type CallExpr struct {
	Type       string
	Value      string
	Parameters []Token
}

func (c CallExpr) validate(value string) bool {
	var result = map[string]string{
		"Entity":   "Entity",
		"Database": "Database",
		"Rel":      "Rel",
		"Queue":    "Queue",
		"Decision": "Decision",
	}[value]
	return result == value
}

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

func tokenizer(source SourceCode) []Token {
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

func eatToken(tokens *[]Token) Token {
	var token Token = (*tokens)[0]
	*tokens = (*tokens)[1:]
	return token
}

func currentToken(tokens []Token) Token {
	return tokens[0]
}

func parseCallExpr(tokens []Token) []CallExpr {
	var expressions []CallExpr

	for len(tokens) > 0 {
		current := eatToken(&tokens)
		if current.Type == Identifier {
			expression := CallExpr{Value: current.Value, Type: "CallExpr"}

			if !expression.validate(current.Value) {
				panic(fmt.Sprintf("Line %v: Unexpected %v call expression", current.LineNumber, current.Value))
			}

			if currentToken(tokens).Type != OpenParen {
				panic(fmt.Sprintf("Line %v: Expected \"(\" character after %v identifier.", currentToken(tokens).LineNumber, current.Value))
			}

			eatToken(&tokens)

			if currentToken(tokens).Type != NumberLiteral && currentToken(tokens).Type != StringLiteral && currentToken(tokens).Type != Identifier {
				panic(fmt.Sprintf("Line %v: Unexpected %v token after \"(\" character.", currentToken(tokens).LineNumber, currentToken(tokens).Type.String()))
			}

			for len(tokens) > 0 && currentToken(tokens).Type != CloseParen {
				c := eatToken(&tokens)
				if c.Type == NumberLiteral || c.Type == StringLiteral || c.Type == Identifier {
					expression.Parameters = append(expression.Parameters, c)
				}
			}

			if currentToken(tokens).Type != CloseParen {
				panic(fmt.Sprintf("Line %v: Expected \")\" character after %v identifier.", currentToken(tokens).LineNumber, current.Value))
			}

			eatToken(&tokens)

			expressions = append(expressions, expression)
		}
	}
	return expressions
}

func parser(source []SourceCode) Program {
	var ast Program = Program{Type: "Program"}
	var tokens []Token

	for len(source) > 0 {
		first := source[0]
		tokens = append(tokens, tokenizer(first)...)
		source = source[1:]
	}

	ast.Body = parseCallExpr(tokens)

	return ast
}

func eatExpr(expr *[]CallExpr) CallExpr {
	var first CallExpr = (*expr)[0]
	*expr = (*expr)[1:]
	return first
}

func transpiler(ast Program) []string {
	var output []string = []string{}
	expressions := ast.Body.([]CallExpr)

	for len(expressions) > 0 {
		expr := eatExpr(&expressions)
		if expr.Value == "Entity" {
			output = append(output, expr.Parameters[0].Value+"["+expr.Parameters[1].Value+"]")
		} else if expr.Value == "Database" {
			output = append(output, expr.Parameters[0].Value+"[("+expr.Parameters[1].Value+")]")
		} else if expr.Value == "Queue" {
			output = append(output, expr.Parameters[0].Value+"[["+expr.Parameters[1].Value+"]]")
		} else if expr.Value == "Decision" {
			output = append(output, expr.Parameters[0].Value+"{"+expr.Parameters[1].Value+"}")
		} else if expr.Value == "Rel" {
			output = append(output, expr.Parameters[0].Value+"-->|"+expr.Parameters[1].Value+"|"+expr.Parameters[2].Value)
		}
	}
	return output
}

type FileScanner struct {
	Scanner *bufio.Scanner
}

func (f *FileScanner) newReader(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	f.Scanner = bufio.NewScanner(file)
}

func (f *FileScanner) ReadFile(file string) []SourceCode {
	f.newReader(file)
	var source []SourceCode
	var lineNumber int = 1
	
	for f.Scanner.Scan() {
		source = append(source, SourceCode{Line: f.Scanner.Text(), LineNumber: lineNumber})
		lineNumber++
	}

	return source
}

type FileWriter struct {
	Writer *bufio.Writer
}

func (f *FileWriter) newWriter(filename string) {
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	f.Writer = bufio.NewWriter(file)
}

func (f *FileWriter) writeLine(line string) {
	_, err := f.Writer.WriteString(line + "\n")
	if err != nil {
		panic(err)
	}
}

func (f *FileWriter) WriteFile(source []string) {
	f.newWriter("output.txt")
	f.writeLine("flowchart LR")
	for _, line := range source {
		f.writeLine(line)
	}
	f.Writer.Flush()
}

func main() {
	var writer FileWriter = FileWriter{}
	var scanner FileScanner = FileScanner{} 
	
	ast := parser(scanner.ReadFile("source.txt"))
	output := transpiler(ast)
	writer.WriteFile(output)
	fmt.Println(output)
}
