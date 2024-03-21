package interpreter

import (
	"fmt"

	"github.com/iamBharatManral/atom.git/cmd/internal/ast"
	"github.com/iamBharatManral/atom.git/cmd/internal/env"
	"github.com/iamBharatManral/atom.git/cmd/internal/error"
	"github.com/iamBharatManral/atom.git/cmd/internal/result"
)

func Eval(node ast.AstNode, env *env.Environment) result.Result {
	switch node := node.(type) {
	case ast.Program:
		return evalStatements(node.Body, env)
	case ast.Literal:
		return evalLiteral(node)
	case ast.BinaryExpression:
		return evalArithmetic(node, env)
	case ast.LetStatement:
		return evalLetStatement(node, env)
	case ast.Identifier:
		return evalIdentifier(node, env)
	case ast.AssignmentStatement:
		return evalAssignment(node, env)
	default:
		return error.UnsupportedTokensError()
	}
}

func evalAssignment(stmt ast.AssignmentStatement, env *env.Environment) result.Result {
	id := stmt.Left.Value
	if _, ok := env.Get(string(id)); !ok {
		return error.UndefinedError(id)
	}
	var rightValue any
	switch right := stmt.Right.(type) {
	case ast.Identifier:
		if _, ok := env.Get(string(right.Value)); ok {
			rightValue = Eval(right, env).Value
		} else {
			return error.UndefinedError(right.Value)
		}
	case ast.Literal:
		rightValue = right.Value
	default:
		return error.SyntaxError("error: wrong type in right hand side of assignment")
	}
	env.Set(id, createResult("identifier", rightValue))
	return result.Result{}
}

func evalStatements(stmts []ast.Statement, env *env.Environment) result.Result {
	var completeResult string
	for i := range stmts {
		result := Eval(stmts[i], env)
		if result.Type == "error" {
			return result
		} else if result.Type == "" {
			continue
		}
		completeResult += fmt.Sprintf("%v\n", result.Value)
	}
	return result.Result{
		Type:  "result",
		Value: completeResult,
	}
}

func evalLetStatement(stmt ast.LetStatement, env *env.Environment) result.Result {
	id := stmt.Left.Value
	switch right := stmt.Right.(type) {
	case ast.Literal:
		env.Set(id, result.Result{
			Type:  "literal",
			Value: right.Value,
		})
	case ast.Identifier:
		return evalIdentifier(stmt.Right.(ast.Identifier), env)
	}
	return result.Result{}
}

func evalIdentifier(stmt ast.Identifier, env *env.Environment) result.Result {
	id := stmt.Value
	if result, ok := env.Get(id); ok {
		return createResult("identifier", result.Value)
	}
	return error.UndefinedError(id)
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

func evalArithmetic(stmt ast.BinaryExpression, env *env.Environment) result.Result {
	switch stmt.Operator {
	case "+":
		return evalAddition(stmt, env)
	case "-":
		return evalSubtraction(stmt, env)
	case "*":
		return evalMultiplication(stmt, env)
	case "/":
		return evalDivision(stmt, env)
	default:
		return error.UnsupportedOperatorError(stmt.Operator)
	}
}

func evalAddition(stmt ast.BinaryExpression, env *env.Environment) result.Result {
	left := Eval(stmt.Left, env)
	if left.Type == "error" {
		return left
	}
	right := Eval(stmt.Right, env)
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

func evalSubtraction(stmt ast.BinaryExpression, env *env.Environment) result.Result {
	left := Eval(stmt.Left, env)
	if left.Type == "error" {
		return left
	}
	right := Eval(stmt.Right, env)
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

func evalMultiplication(stmt ast.BinaryExpression, env *env.Environment) result.Result {
	left := Eval(stmt.Left, env)
	if left.Type == "error" {
		return left
	}
	right := Eval(stmt.Right, env)
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

func evalDivision(stmt ast.BinaryExpression, env *env.Environment) result.Result {
	left := Eval(stmt.Left, env)
	if left.Type == "error" {
		return left
	}
	right := Eval(stmt.Right, env)
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
