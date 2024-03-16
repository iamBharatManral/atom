package interpreter

import "github.com/iamBharatManral/atom.git/cmd/internal/ast"

func Eval(program ast.Program) any {
	if len(program.Body) == 0 {
		return nil
	}
	node := program.Body[0]
	switch node.(type) {
	case ast.IntegerLiteral:
		return node.(ast.IntegerLiteral).Value
	case ast.FloatLiteral:
		return node.(ast.FloatLiteral).Value
	default:
		return nil
	}
}
