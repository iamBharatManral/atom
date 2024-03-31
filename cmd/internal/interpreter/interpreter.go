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
		return evalBinaryExpression(node, env)
	case ast.LetStatement:
		return evalLetStatement(node, env)
	case ast.Identifier:
		return evalIdentifier(node, env)
	case ast.AssignmentStatement:
		return evalAssignment(node, env)
	case ast.IfBlock:
		return evalIfExpression(node, env)
	case ast.IfElseBlock:
		return evalIfElseExpression(node, env)
	case ast.FunctionExpression:
		return evalFunctionExpression(node, env, "")
	case ast.FunctionEvaluation:
		return evalFunction(node, env)
	case ast.ReturnStatement:
		return evalReturnStatement(node, env)
	default:
		return error.UnsupportedTokensError()
	}
}

func evalReturnStatement(node ast.ReturnStatement, ev *env.Environment) result.Result {
	switch v := node.Value.(type) {
	case ast.Identifier:
		return evalIdentifier(v, ev)
	case ast.Literal:
		return createResult("literal", v.Value)
	case ast.BinaryExpression:
		return evalBinaryExpression(v, ev)

	}
	return result.Result{}
}

func evalFunction(node ast.FunctionEvaluation, ev *env.Environment) result.Result {
	localEnv := env.New()
	fnName := node.Name.Value
	var funcDecl ast.FunctionExpression
	if fn, ok := ev.Get(fnName); !ok {
		return error.UndefinedError(fnName)
	} else {
		funcDecl = fn.Value.(ast.FunctionExpression)
	}
	if fnName != node.Name.Value {
		return error.UndefinedError(node.Name.Value)
	}
	if len(node.Parameters) != len(funcDecl.Parameters) {
		return error.NotEnoughArguments(fmt.Sprintf("not enough arguments. require: %d, got: %d", len(funcDecl.Parameters), len(node.Parameters)))
	}
	for i, stmt := range node.Parameters {
		switch stmt := stmt.(type) {
		case ast.Identifier:
			result := evalIdentifier(stmt, ev)
			if result.Type == "error" {
				return error.UndefinedError(stmt.Value)
			}
			localEnv.Set(stmt.Value, createResult("identifier", result.Value))
		case ast.Literal:
			localEnv.Set(funcDecl.Parameters[i].Value, createResult("identifier", stmt.Value))
		}
	}

	var results []result.Result
	for _, stmt := range funcDecl.Body {
		if _, ok := stmt.(ast.ReturnStatement); ok {
			return Eval(stmt, localEnv)
		} else {

			results = append(results, Eval(stmt, localEnv))
		}
	}
	if len(results) > 0 {
		return results[len(results)-1]
	}
	return result.Result{}
}

func evalFunctionExpression(stmt ast.FunctionExpression, env *env.Environment, fnName string) result.Result {
	name := stmt.Name.Value
	if name == "" {
		name = fnName
	}
	if _, ok := env.Get(name); ok {
		return error.UnsupportedOperation(fmt.Sprintf("symbol '%s' is already defined", name))
	}
	env.Set(name, createResult("fn", stmt))
	return createResult("function declaration", "()")
}

func evalAssignment(stmt ast.AssignmentStatement, env *env.Environment) result.Result {
	id := stmt.Left.Value
	if _, ok := env.Get(string(id)); !ok {
		return error.UndefinedError(id)
	}
	return error.UnsupportedOperation("re-assignment is not supported")
}

func evalRHS(stmt ast.LetStatement, env *env.Environment) result.Result {
	id := stmt.Left.Value
	if _, ok := env.Get(id); ok {
		env.Delete(id)
	}
	switch right := stmt.Right.(type) {
	case ast.Literal:
		env.Set(id, createResult("literal", right.Value))
	case ast.Identifier:
		r := evalIdentifier(stmt.Right.(ast.Identifier), env)
		if r.Type == "error" {
			return r
		}
		env.Set(id, createResult("identifier", r.Value))
	case ast.BinaryExpression:
		r := evalBinaryExpression(stmt.Right.(ast.BinaryExpression), env)
		if r.Type == "error" {
			return r
		}
		env.Set(id, createResult("BinaryExpression", r.Value))
	case ast.FunctionExpression:
		return evalFunctionExpression(stmt.Right.(ast.FunctionExpression), env, id)
	case ast.IfBlock:
		r := evalIfExpression(right, env)
		if r.Type == "error" {
			return r
		}
		env.Set(id, createResult("IfExpression", r.Value))
	case ast.IfElseBlock:
		r := evalIfElseExpression(right, env)
		if r.Type == "error" {
			return r
		}
		env.Set(id, createResult("IfElseExpression", r.Value))
	default:
		return error.UnsupportedTokensError()
	}
	return result.Result{}

}
func evalStatements(stmts []ast.Statement, env *env.Environment) result.Result {
	var completeResult string
	for i := range stmts {
		result := Eval(stmts[i], env)
		if result.Type == "error" {
			return result
		} else if result.Type == "" || result.Value == "()" {
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
	return evalRHS(stmt, env)
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

func evalIfExpression(stmt ast.IfBlock, env *env.Environment) result.Result {
	test := stmt.Test
	var testResult result.Result
	switch test := test.(type) {
	case ast.BinaryExpression:
		testResult = evalBinaryExpression(test, env)
	case ast.Literal:
		testResult = createResult("boolean", test.Value)
	}
	if testResult.Value == true {
		finalResult := Eval(stmt.Consequent, env)
		return createResult("conditional", finalResult.Value)
	}
	return result.Result{}

}

func evalIfElseExpression(stmt ast.IfElseBlock, env *env.Environment) result.Result {
	test := stmt.Test
	var testResult result.Result
	switch test := test.(type) {
	case ast.BinaryExpression:
		testResult = evalBinaryExpression(test, env)
	case ast.Literal:
		testResult = createResult("boolean", test.Value)
	}
	if testResult.Value == true {
		finalResult := Eval(stmt.Consequent, env)
		return createResult("conditional", finalResult.Value)
	} else {
		finalResult := Eval(stmt.Alternate, env)
		return createResult("conditional", finalResult.Value)
	}

}

func typeAsString(v any) string {
	return fmt.Sprintf("%v", v)
}

func evalBinaryExpression(stmt ast.BinaryExpression, env *env.Environment) result.Result {
	var tempLeft, tempRight any
	if _, ok := stmt.Left.(ast.BinaryExpression); ok {
		tempLeft = evalBinaryExpression(stmt.Left.(ast.BinaryExpression), env)
	} else {
		tempLeft = Eval(stmt.Left, env)
	}
	if _, ok := stmt.Right.(ast.BinaryExpression); ok {
		tempRight = evalBinaryExpression(stmt.Right.(ast.BinaryExpression), env)
	} else {
		tempRight = Eval(stmt.Right, env)
	}
	left := createResult(tempLeft.(result.Result).Type, tempLeft.(result.Result).Value)
	right := createResult(tempRight.(result.Result).Type, tempRight.(result.Result).Value)
	switch stmt.Operator {
	case "+":
		return evalAddition(left, right)
	case "-":
		return evalSubtraction(left, right)
	case "*":
		return evalMultiplication(left, right)
	case "/":
		return evalDivision(left, right)
	case "<":
		return evalLessThan(left, right)
	case "<=":
		return evalLessThanEqual(left, right)
	case ">":
		return evalGreaterThan(left, right)
	case ">=":
		return evalGreaterThanEqual(left, right)
	case "!=":
		return evalNotEqual(left, right)
	case "==":
		return evalEqualEqual(left, right)
	case "and":
		return evalLogicalAnd(left, right)
	case "or":
		return evalLogicalOr(left, right)
	default:
		return error.UnsupportedOperatorError(stmt.Operator)
	}
}

func evalLogicalAnd(left, right result.Result) result.Result {
	if left.Type == "error" {
		return left
	}
	if right.Type == "error" {
		return right
	}
	switch left := left.Value.(type) {
	case bool:
		if right, ok := right.Value.(bool); ok {
			return createResult("bool", left && right)
		}
		return error.TypeMismatchError(left, right.Value)
	}
	return error.UnsupportedTypeError(left, "and")

}

func evalLogicalOr(left, right result.Result) result.Result {
	if left.Type == "error" {
		return left
	}
	if right.Type == "error" {
		return right
	}
	switch left := left.Value.(type) {
	case bool:
		if right, ok := right.Value.(bool); ok {
			return createResult("bool", left || right)
		}
		return error.TypeMismatchError(left, right.Value)
	}
	return error.UnsupportedTypeError(left, "or")

}

func evalEqualEqual(left, right result.Result) result.Result {
	if left.Type == "error" {
		return left
	}
	if right.Type == "error" {
		return right
	}
	switch left := left.Value.(type) {
	case int:
		if right, ok := right.Value.(int); ok {
			return createResult("int", left == right)
		}
		return error.TypeMismatchError(left, right.Value)
	case float64:
		if right, ok := right.Value.(float64); ok {
			return createResult("float", left == right)
		}
		return error.TypeMismatchError(left, right.Value)
	case string:
		if right, ok := right.Value.(string); ok {
			return createResult("string", left == right)
		}
		return error.TypeMismatchError(left, right.Value)
	case bool:
		if right, ok := right.Value.(bool); ok {
			return createResult("bool", left == right)
		}
		return error.TypeMismatchError(left, right.Value)
	}
	return error.UnsupportedTypeError(left, "==")

}

func evalNotEqual(left, right result.Result) result.Result {
	if left.Type == "error" {
		return left
	}
	if right.Type == "error" {
		return right
	}
	switch left := left.Value.(type) {
	case int:
		if right, ok := right.Value.(int); ok {
			return createResult("int", left != right)
		}
		return error.TypeMismatchError(left, right.Value)
	case float64:
		if right, ok := right.Value.(float64); ok {
			return createResult("float", left != right)
		}
		return error.TypeMismatchError(left, right.Value)
	case string:
		if right, ok := right.Value.(string); ok {
			return createResult("string", left != right)
		}
		return error.TypeMismatchError(left, right.Value)
	case bool:
		if right, ok := right.Value.(bool); ok {
			return createResult("bool", left == right)
		}
		return error.TypeMismatchError(left, right.Value)
	}

	return error.UnsupportedTypeError(left, "!=")
}

func evalGreaterThanEqual(left, right result.Result) result.Result {
	if left.Type == "error" {
		return left
	}
	if right.Type == "error" {
		return right
	}
	switch left := left.Value.(type) {
	case int:
		if right, ok := right.Value.(int); ok {
			return createResult("int", left >= right)
		}
		return error.TypeMismatchError(left, right.Value)
	case float64:
		if right, ok := right.Value.(float64); ok {
			return createResult("float", left >= right)
		}
		return error.TypeMismatchError(left, right.Value)
	case string:
		if right, ok := right.Value.(string); ok {
			return createResult("string", left >= right)
		}
		return error.TypeMismatchError(left, right.Value)
	}
	return error.UnsupportedTypeError(left, ">=")
}

func evalGreaterThan(left, right result.Result) result.Result {
	if left.Type == "error" {
		return left
	}
	if right.Type == "error" {
		return right
	}
	switch left := left.Value.(type) {
	case int:
		if right, ok := right.Value.(int); ok {
			return createResult("int", left > right)
		}
		return error.TypeMismatchError(left, right.Value)
	case float64:
		if right, ok := right.Value.(float64); ok {
			return createResult("float", left > right)
		}
		return error.TypeMismatchError(left, right.Value)
	case string:
		if right, ok := right.Value.(string); ok {
			return createResult("string", left > right)
		}
		return error.TypeMismatchError(left, right.Value)
	}
	return error.UnsupportedTypeError(left, ">")
}

func evalLessThanEqual(left, right result.Result) result.Result {
	if left.Type == "error" {
		return left
	}
	if right.Type == "error" {
		return right
	}
	switch left := left.Value.(type) {
	case int:
		if right, ok := right.Value.(int); ok {
			return createResult("int", left <= right)
		}
		return error.TypeMismatchError(left, right.Value)
	case float64:
		if right, ok := right.Value.(float64); ok {
			return createResult("float", left <= right)
		}
		return error.TypeMismatchError(left, right.Value)
	case string:
		if right, ok := right.Value.(string); ok {
			return createResult("string", left <= right)
		}
		return error.TypeMismatchError(left, right.Value)
	}
	return error.UnsupportedTypeError(left, "<=")
}

func evalLessThan(left, right result.Result) result.Result {
	if left.Type == "error" {
		return left
	}
	if right.Type == "error" {
		return right
	}
	switch left := left.Value.(type) {
	case int:
		if right, ok := right.Value.(int); ok {
			return createResult("int", left < right)
		}
		return error.TypeMismatchError(left, right.Value)
	case float64:
		if right, ok := right.Value.(float64); ok {
			return createResult("float", left < right)
		}
		return error.TypeMismatchError(left, right.Value)
	case string:
		if right, ok := right.Value.(string); ok {
			return createResult("string", left < right)
		}
		return error.TypeMismatchError(left, right.Value)
	}
	return error.UnsupportedTypeError(left, "<")
}

func evalAddition(left, right result.Result) result.Result {
	if left.Type == "error" {
		return left
	}
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

func evalSubtraction(left, right result.Result) result.Result {
	if left.Type == "error" {
		return left
	}
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

func evalMultiplication(left, right result.Result) result.Result {
	if left.Type == "error" {
		return left
	}
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

func evalDivision(left, right result.Result) result.Result {
	if left.Type == "error" {
		return left
	}
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
