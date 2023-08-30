package parser

import (
	"fmt"
	ty "golang-transpiler/pkg/types"
	tk "golang-transpiler/pkg/tokenizer"
)

type CallExpr struct {
	Type       string
	Value      string
	Parameters []tk.Token
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

func eatToken(tokens *[]tk.Token) tk.Token {
	var token tk.Token = (*tokens)[0]
	*tokens = (*tokens)[1:]
	return token
}

func currentToken(tokens []tk.Token) tk.Token {
	return tokens[0]
}

func parseCallExpr(tokens []tk.Token) []CallExpr {
	var expressions []CallExpr

	for len(tokens) > 0 {
		current := eatToken(&tokens)
		if current.Type == tk.Identifier {
			expression := CallExpr{Value: current.Value, Type: "CallExpr"}

			if !expression.validate(current.Value) {
				panic(fmt.Sprintf("Line %v: Unexpected %v call expression", current.LineNumber, current.Value))
			}

			if currentToken(tokens).Type != tk.OpenParen {
				panic(fmt.Sprintf("Line %v: Expected \"(\" character after %v identifier.", currentToken(tokens).LineNumber, current.Value))
			}

			eatToken(&tokens)

			if currentToken(tokens).Type != tk.NumberLiteral && 
				currentToken(tokens).Type != tk.StringLiteral && 
				currentToken(tokens).Type != tk.Identifier {
				panic(fmt.Sprintf("Line %v: Unexpected %v token after \"(\" character.", currentToken(tokens).LineNumber, currentToken(tokens).Type.String()))
			}

			for len(tokens) > 0 && currentToken(tokens).Type != tk.CloseParen {
				c := eatToken(&tokens)
				if c.Type == tk.NumberLiteral || c.Type == tk.StringLiteral || c.Type == tk.Identifier {
					expression.Parameters = append(expression.Parameters, c)
				}
			}

			if currentToken(tokens).Type != tk.CloseParen {
				panic(fmt.Sprintf("Line %v: Expected \")\" character after %v identifier.", currentToken(tokens).LineNumber, current.Value))
			}

			eatToken(&tokens)

			expressions = append(expressions, expression)
		}
	}
	return expressions
}

func Parse(source []ty.SourceCode) []CallExpr {
	var tokens []tk.Token

	for len(source) > 0 {
		first := source[0]
		tokens = append(tokens, tk.Tokenizer(first)...)
		source = source[1:]
	}

	return parseCallExpr(tokens)
}