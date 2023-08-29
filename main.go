package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

type SourceCode struct {
	Line string
	LineNumber int
}

type Program struct {
	Type string
	Body interface{}
}

type CallExpr struct {
	Type string
	Value string
	Parameters []Token
}

func (c CallExpr) validate(value string) bool {
	var result = map[string]string {
		"Entity": "Entity",
		"Database": "Database",
		"Rel": "Rel",
	}[value]
	return result == value
}

type Token struct {
	Type  TokenType
	Value string
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
	var tokens []Token

	for len(input) > 0 {
		currentValue := current(input)
		if currentValue == "(" {
			tokens = append(tokens, Token{OpenParen, shift(&input)})
		} else if currentValue == ")" {
			tokens = append(tokens, Token{CloseParen, shift(&input)})

			for len(input) > 0 {
				if shift(&input) != "\n" || shift(&input) != "\t" || shift(&input) != " " {
					panic("Unexpected token after "+currentValue)
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
				tokens = append(tokens, Token{NumberLiteral, string(number)})
			} else if isAlpha(currentValue) {
				var identifier string
				for len(input) > 0 && isAlpha(current(input)) {
					identifier += shift(&input)
				}
				tokens = append(tokens, Token{Identifier, string(identifier)})
			} else if currentValue == "\"" {
				shift(&input)
				var stringLiteral string
				for len(input) > 0 && current(input) != "\"" {
					stringLiteral += shift(&input)
				}
				shift(&input)
				tokens = append(tokens, Token{StringLiteral, string(stringLiteral)})
			} else {
				panic("Unexpected token: "+currentValue)
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
			
			if ! expression.validate(current.Value) {
				panic("Unexpected '"+current.Value+"' call expression")
			}
			
			if currentToken(tokens).Type != OpenParen {
				panic("Expected \"(\" character after "+current.Value+" identifier")
			}

			eatToken(&tokens)

			if currentToken(tokens).Type != NumberLiteral && 
			   currentToken(tokens).Type != StringLiteral && 
			   currentToken(tokens).Type != Identifier {
				panic("Unexpected "+currentToken(tokens).Type.String()+" token after \"(\" character")
			}
			
			for len(tokens) > 0 && currentToken(tokens).Type != CloseParen {
				c := eatToken(&tokens)
				if c.Type == NumberLiteral || c.Type == StringLiteral || c.Type == Identifier {
					expression.Parameters = append(expression.Parameters, c)
				}
			}

			if currentToken(tokens).Type != CloseParen {
				panic("Expected \")\" character after "+current.Value+" identifier")
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

func readFile(filename string) []SourceCode {
	f, err := os.Open(filename)

	if err != nil {
		panic(err)
	}

	defer f.Close()
	
	scanner := bufio.NewScanner(f)
	
	var source []SourceCode

	var lineNumber int = 0

	for scanner.Scan() {
		source = append(source, SourceCode{ Line: scanner.Text(), LineNumber: lineNumber })
		lineNumber++
    }
	fmt.Println(source)
	return source
}

func identifiersToMap(expressions []CallExpr) map[string]string {
	var identifiers map[string]string = make(map[string]string)

	for _, expr := range expressions {
		if expr.Value == "Entity" {
			identifiers[expr.Parameters[0].Value] = "["+expr.Parameters[1].Value+"]"
			eatExpr(&expressions)
		} else if expr.Value == "Database" {
			identifiers[expr.Parameters[0].Value] = "[("+expr.Parameters[1].Value+")]"
			eatExpr(&expressions)
		}
	}

	return identifiers
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
		} else if expr.Value == "Rel" {
			output = append(output,expr.Parameters[0].Value+"-->|"+expr.Parameters[1].Value+"|"+expr.Parameters[2].Value)
		}
	}
	return output
}

func writeFile(source []string) {
	file, err := os.Create("output.txt")
    if err != nil {
        panic(err)
    }
    writer := bufio.NewWriter(file)
	for _, line := range source {
        bytesWritten, err := writer.WriteString(line + "\n")
        if err != nil {
            panic(err)
        }
        fmt.Printf("Bytes Written: %d\n", bytesWritten)
        fmt.Printf("Available: %d\n", writer.Available())
        fmt.Printf("Buffered : %d\n", writer.Buffered())
    }
    writer.Flush()
}

func main() {
	ast := parser(readFile("source.txt"))
	output := transpiler(ast)
	writeFile(output)
	fmt.Println(output)
}