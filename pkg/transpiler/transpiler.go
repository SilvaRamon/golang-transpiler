package transpiler

import (
	ty "golang-transpiler/pkg/types"
	ps "golang-transpiler/pkg/parser"
)

func eatExpr(expr *[]ps.CallExpr) ps.CallExpr {
	var first ps.CallExpr = (*expr)[0]
	*expr = (*expr)[1:]
	return first
}

func Transpile(source []ty.SourceCode) []string {
	var output []string = []string{}

	expressions := ps.Parse(source)

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