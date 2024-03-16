package interpreter

import (
	"fmt"

	"github.com/iamBharatManral/atom.git/cmd/internal/ast"
)

func Eval(program ast.Program) {
	if len(program.Body) == 0 {
		return
	}
	for _, node := range program.Body {
		switch node.(type) {
		case ast.IntegerLiteral:
			fmt.Println(node.(ast.IntegerLiteral).Value)
		case ast.FloatLiteral:
			fmt.Println(node.(ast.FloatLiteral).Value)
		case ast.StringLiteral:
			fmt.Println(node.(ast.StringLiteral).Value)
		default:
			return
		}
	}
}
