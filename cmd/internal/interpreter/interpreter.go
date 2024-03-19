package interpreter

import (
	"fmt"

	"github.com/iamBharatManral/atom.git/cmd/internal/ast"
	"github.com/iamBharatManral/atom.git/cmd/internal/error"
	"github.com/iamBharatManral/atom.git/cmd/internal/result"
)

func Eval(node ast.AstNode) result.Result {
	switch node := node.(type) {
	case ast.Program:
		return evalStatements(node.Body)
	case ast.Literal:
		return evalLiteral(node)
	case ast.BinaryExpression:
		return evalArithmetic(node)
	default:
		return error.UnsupportedTokensError()
	}
}

func evalStatements(stmts []ast.Statement) result.Result {
	var completeResult string
	for i := range stmts {
		result := Eval(stmts[i])
		if result.Type == "error" {
			return result
		}
		completeResult += fmt.Sprintf("%v\n", result.Value)
	}
	return result.Result{
		Type:  "result",
		Value: completeResult,
	}
}

func evalLiteral(stmt ast.Literal) result.Result {
	return result.Result{
		Type:  typeAsString(stmt.Value),
		Value: stmt.Value,
	}
}

func typeAsString(v any) string {
	return fmt.Sprintf("%v", v)
}

func evalArithmetic(stmt ast.BinaryExpression) result.Result {
	switch stmt.Operator {
	case "+":
		return evalAddition(stmt)
	case "-":
		return evalSubtraction(stmt)
	case "*":
		return evalMultiplication(stmt)
	case "/":
		return evalDivision(stmt)
	default:
		return error.UnsupportedOperatorError(stmt.Operator)
	}
}

func evalAddition(stmt ast.BinaryExpression) result.Result {
	left := Eval(stmt.Left)
	if left.Type == "error" {
		return left
	}
	right := Eval(stmt.Right)
	if right.Type == "error" {
		return right
	}
	switch left := left.Value.(type) {
	case int:
		if right, ok := right.Value.(int); ok {
			return createResult("int", left+right)
		}
		return error.TypeMismatchError(left, right.Value)
	case float64:
		if right, ok := right.Value.(float64); ok {
			return createResult("float", left+right)
		}
		return error.TypeMismatchError(left, right.Value)
	case string:
		if right, ok := right.Value.(string); ok {
			return createResult("string", left+right)
		}
		return error.TypeMismatchError(left, right.Value)
	}
	return error.UnsupportedTypeError(left, "+")
}

func evalSubtraction(stmt ast.BinaryExpression) result.Result {
	left := Eval(stmt.Left)
	if left.Type == "error" {
		return left
	}
	right := Eval(stmt.Right)
	if right.Type == "error" {
		return right
	}
	switch left := left.Value.(type) {
	case int:
		if right, ok := right.Value.(int); ok {
			return createResult("int", left-right)
		}
		return error.TypeMismatchError(left, right.Value)
	case float64:
		if right, ok := right.Value.(float64); ok {
			return createResult("float", left-right)
		}
		return error.TypeMismatchError(left, right.Value)
	}
	return error.UnsupportedTypeError(left, "-")
}

func createResult(t string, v any) result.Result {
	return result.Result{
		Type:  t,
		Value: v,
	}
}

func evalMultiplication(stmt ast.BinaryExpression) result.Result {
	left := Eval(stmt.Left)
	if left.Type == "error" {
		return left
	}
	right := Eval(stmt.Right)
	if right.Type == "error" {
		return right
	}
	switch left := left.Value.(type) {
	case int:
		if right, ok := right.Value.(int); ok {
			return createResult("int", left*right)
		}
		return error.TypeMismatchError(left, right.Value)
	case float64:
		if right, ok := right.Value.(float64); ok {
			return createResult("float", left*right)
		}
		return error.TypeMismatchError(left, right.Value)
	}
	return error.UnsupportedTypeError(left, "*")
}

func evalDivision(stmt ast.BinaryExpression) result.Result {
	left := Eval(stmt.Left)
	if left.Type == "error" {
		return left
	}
	right := Eval(stmt.Right)
	if right.Type == "error" {
		return right
	}
	if right.Value == 0 {
		return error.DivisonByZeroError()
	}
	switch left := left.Value.(type) {
	case int:
		if right, ok := right.Value.(int); ok {
			return createResult("int", left/right)
		}
		return error.TypeMismatchError(left, right.Value)
	case float64:
		if right, ok := right.Value.(float64); ok {
			return createResult("float", left/right)
		}
		return error.TypeMismatchError(left, right.Value)
	}
	return error.UnsupportedTypeError(left, "/")
}
