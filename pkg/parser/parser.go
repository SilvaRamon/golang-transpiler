package parser

import (
	"fmt"
	ty "golang-transpiler/pkg/types"
	tk "golang-transpiler/pkg/tokenizer"
)

type CallExpr struct {
	Type       string
	Value      string
	LineNumber int
	Parameters []tk.Token
}

func (c CallExpr) validateCallExprIdentifier(value string) bool {
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

func validateEntityParameters(expression CallExpr) {
	if expression.Value == "Rel" {
		return
	}
	if len(expression.Parameters) == 0 || len(expression.Parameters) != 2 {
		panic(fmt.Sprintf("Line %v: Entity call expression must have 2 parameters.", expression.LineNumber))
	}
	if expression.Parameters[0].Type != tk.Identifier {
		panic(fmt.Sprintf("Line %v: Entity call expression first parameter must be an identifier.", expression.Parameters[0].LineNumber))
	}
	if expression.Parameters[1].Type != tk.StringLiteral {
		panic(fmt.Sprintf("Line %v: Entity call expression second parameter must be a string literal.", expression.Parameters[1].LineNumber))
	}
}

func validateRelParameters(expression CallExpr) {
	if expression.Value != "Rel" {
		return
	}
	if len(expression.Parameters) == 0 || len(expression.Parameters) != 3 {
		panic(fmt.Sprintf("Line %v: Rel call expression must have 3 parameters.", expression.LineNumber))
	}
	if expression.Parameters[0].Type != tk.Identifier {
		panic(fmt.Sprintf("Line %v: Rel call expression first parameter must be an identifier.", expression.Parameters[1].LineNumber))
	}
	if expression.Parameters[1].Type != tk.StringLiteral {
		panic(fmt.Sprintf("Line %v: Rel call expression second parameter must be a string literal.", expression.Parameters[2].LineNumber))
	}
	if expression.Parameters[2].Type != tk.Identifier {
		panic(fmt.Sprintf("Line %v: Rel call expression third parameter must be an identifier.", expression.Parameters[2].LineNumber))
	}
}

func parseTokensToCallExprs(tokens []tk.Token) []CallExpr {
	var expressions []CallExpr

	for len(tokens) > 0 {
		current := eatToken(&tokens)
		if current.Type == tk.Identifier {
			expression := CallExpr{Value: current.Value, Type: "CallExpr", LineNumber: current.LineNumber}

			if !expression.validateCallExprIdentifier(current.Value) {
				panic(fmt.Sprintf("Line %v: Unexpected %v call expression", current.LineNumber, current.Value))
			}

			if currentToken(tokens).Type != tk.OpenParen {
				panic(fmt.Sprintf("Line %v: Expected \"(\" character after %v identifier.", currentToken(tokens).LineNumber, current.Value))
			}

			eatToken(&tokens)

			for len(tokens) > 0 && currentToken(tokens).Type != tk.CloseParen {
				c := eatToken(&tokens)
				if c.Type == tk.NumberLiteral || c.Type == tk.StringLiteral || c.Type == tk.Identifier {
					expression.Parameters = append(expression.Parameters, c)
				} else {
					panic(fmt.Sprintf("Line %v: Unexpected %v token after \"(\" character.", currentToken(tokens).LineNumber, currentToken(tokens).Type.String()))
				}
			}

			if currentToken(tokens).Type != tk.CloseParen {
				panic(fmt.Sprintf("Line %v: Expected \")\" character after %v identifier.", currentToken(tokens).LineNumber, current.Value))
			}

			eatToken(&tokens)

			validateEntityParameters(expression)
			validateRelParameters(expression)

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

	return parseTokensToCallExprs(tokens)
}