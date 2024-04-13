package parser

import (
	"reflect"
	"testing"

	"github.com/iamBharatManral/atom.git/cmd/internal/ast"
	"github.com/iamBharatManral/atom.git/cmd/internal/lexer"
)

func TestLiteralsAndExpressions(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []ast.Statement
	}{
		{name: "integer literal 1933", want: []ast.Statement{
			ast.Literal{
				Value: 1933,
				Node: ast.Node{
					Type:  "Literal",
					Start: 0,
					End:   3,
				},
				UnaryOp: "",
			},
		}, input: "1933"},
		{name: "negative number", want: []ast.Statement{
			ast.Literal{
				Value: 10,
				Node: ast.Node{
					Type:  "Literal",
					Start: 2,
					End:   3,
				},
				UnaryOp: "-",
			},
		}, input: "(-10)"},
		{name: "not '!' operator", want: []ast.Statement{
			ast.Identifier{
				Value: "true",
				Node: ast.Node{
					Type:  "Identifier",
					Start: 1,
					End:   4,
				},
				UnaryOp: "!",
			},
		}, input: "!true"},
		{name: "string literal", want: []ast.Statement{
			ast.Literal{
				Value: "hello world",
				Node: ast.Node{
					Type:  "Literal",
					Start: 0,
					End:   12,
				},
				UnaryOp: "",
			},
		}, input: "\"hello world\""},

		{name: "addtion of two numbers 12 + 13", want: []ast.Statement{
			ast.BinaryExpression{
				Left: ast.Literal{
					Node: ast.Node{
						Type:  "Literal",
						Start: 0,
						End:   1,
					},
					Value: 12,
				},
				Right: ast.Literal{
					Node: ast.Node{
						Type:  "Literal",
						Start: 5,
						End:   6,
					},
					Value: 13,
				},
				Operator: "+",
				Node: ast.Node{
					Start: 0,
					End:   6,
					Type:  "BinaryExpression",
				},
			},
		}, input: "12 + 13"},
		{name: "addtion of two numbers within bracket(12 + 13)", want: []ast.Statement{
			ast.BinaryExpression{
				Left: ast.Literal{
					Node: ast.Node{
						Type:  "Literal",
						Start: 1,
						End:   2,
					},
					Value: 12,
				},
				Right: ast.Literal{
					Node: ast.Node{
						Type:  "Literal",
						Start: 6,
						End:   7,
					},
					Value: 13,
				},
				Operator: "+",
				Node: ast.Node{
					Start: 1,
					End:   7,
					Type:  "BinaryExpression",
				},
			},
		}, input: "(12 + 13)"},
		{name: "multiplication of two numbers 5 * 9 and division of two numbers 96 / 4", want: []ast.Statement{
			ast.BinaryExpression{
				Left: ast.Literal{
					Node: ast.Node{
						Type:  "Literal",
						Start: 0,
						End:   0,
					},
					Value: 5,
				},
				Right: ast.Literal{
					Node: ast.Node{
						Type:  "Literal",
						Start: 4,
						End:   4,
					},
					Value: 9,
				},
				Operator: "*",
				Node: ast.Node{
					Start: 0,
					End:   4,
					Type:  "BinaryExpression",
				},
			},
			ast.BinaryExpression{
				Left: ast.Literal{
					Node: ast.Node{
						Type:  "Literal",
						Start: 6,
						End:   7,
					},
					Value: 96,
				},
				Right: ast.Literal{
					Node: ast.Node{
						Type:  "Literal",
						Start: 11,
						End:   11,
					},
					Value: 4,
				},
				Operator: "/",
				Node: ast.Node{
					Start: 6,
					End:   11,
					Type:  "BinaryExpression",
				},
			},
		}, input: `5 * 9
96 / 4
      `},
		{name: "multiplication of two numbers 5 * 9 and literal string hello with 1 and another binary expression", want: []ast.Statement{
			ast.BinaryExpression{
				Left: ast.Literal{
					Node: ast.Node{
						Type:  "Literal",
						Start: 0,
						End:   0,
					},
					Value: 5,
				},
				Right: ast.Literal{
					Node: ast.Node{
						Type:  "Literal",
						Start: 4,
						End:   4,
					},
					Value: 9,
				},
				Operator: "*",
				Node: ast.Node{
					Start: 0,
					End:   4,
					Type:  "BinaryExpression",
				},
			},
			ast.Literal{
				Node: ast.Node{
					Start: 6,
					End:   18,
					Type:  "Literal",
				},
				Value: "hello world",
			},
			ast.Literal{
				Node: ast.Node{
					Start: 20,
					End:   20,
					Type:  "Literal",
				},
				Value: 1,
			},
			ast.BinaryExpression{
				Left: ast.Literal{
					Node: ast.Node{
						Type:  "Literal",
						Start: 22,
						End:   22,
					},
					Value: 4,
				},
				Right: ast.Literal{
					Node: ast.Node{
						Type:  "Literal",
						Start: 25,
						End:   25,
					},
					Value: 5,
				},
				Operator: "*",
				Node: ast.Node{
					Start: 22,
					End:   25,
					Type:  "BinaryExpression",
				},
			},
		}, input: `5 * 9
"hello world"
1
4 *5
`},
		{name: "let declaration", want: []ast.Statement{
			ast.LetStatement{
				Left: ast.Identifier{
					Node: ast.Node{
						Type:  "Identifier",
						Start: 4,
						End:   4,
					},
					Value: "a",
				},
				Right: ast.Literal{
					Node: ast.Node{
						Type:  "Literal",
						Start: 8,
						End:   9,
					},
					Value: 10,
				},
				Operator: "=",
				Node: ast.Node{
					Start: 0,
					End:   9,
					Type:  "LetStatement",
				},
			},
		}, input: `let a = 10`},
		{name: "assigment operation name = \"hello\"", want: []ast.Statement{
			ast.AssignmentStatement{
				Left: ast.Identifier{
					Node: ast.Node{
						Type:  "Identifier",
						Start: 0,
						End:   3,
					},
					Value: "name",
				},
				Right: ast.Literal{
					Node: ast.Node{
						Type:  "Literal",
						Start: 7,
						End:   13,
					},
					Value: "hello",
				},
				Operator: "=",
				Node: ast.Node{
					Start: 0,
					End:   13,
					Type:  "Assignment",
				},
			},
		}, input: `name = "hello"`},
		{name: "multiple arithmetic expressions", want: []ast.Statement{
			ast.BinaryExpression{
				Left: ast.Literal{
					Node: ast.Node{
						Type:  "Literal",
						Start: 0,
						End:   1,
					},
					Value: 51,
				},
				Right: ast.BinaryExpression{
					Node: ast.Node{
						Type:  "BinaryExpression",
						Start: 5,
						End:   10,
					},
					Operator: "*",
					Left: ast.Literal{
						Node: ast.Node{
							Start: 5,
							End:   6,
							Type:  "Literal",
						},
						Value: 23,
					},
					Right: ast.Literal{
						Node: ast.Node{
							Start: 10,
							End:   10,
							Type:  "Literal",
						},
						Value: 4,
					},
				},
				Operator: "+",
				Node: ast.Node{
					Start: 0,
					End:   10,
					Type:  "BinaryExpression",
				},
			},
		}, input: `51 + 23 * 4`},
		{name: "less than comparison between 2 numbers", want: []ast.Statement{
			ast.BinaryExpression{
				Left: ast.Literal{
					Node: ast.Node{
						Type:  "Literal",
						Start: 0,
						End:   1,
					},
					Value: 12,
				},
				Right: ast.Literal{
					Node: ast.Node{
						Type:  "Literal",
						Start: 5,
						End:   6,
					},
					Value: 13,
				},
				Operator: "<",
				Node: ast.Node{
					Start: 0,
					End:   6,
					Type:  "BinaryExpression",
				},
			},
		}, input: "12 < 13"},
		{name: "greater than equal to between 2 strings", want: []ast.Statement{
			ast.BinaryExpression{
				Left: ast.Literal{
					Node: ast.Node{
						Type:  "Literal",
						Start: 0,
						End:   6,
					},
					Value: "hello",
				},
				Right: ast.Literal{
					Node: ast.Node{
						Type:  "Literal",
						Start: 11,
						End:   15,
					},
					Value: "bye",
				},
				Operator: ">=",
				Node: ast.Node{
					Start: 0,
					End:   15,
					Type:  "BinaryExpression",
				},
			},
		}, input: "\"hello\" >= \"bye\""},
		{name: "if block", want: []ast.Statement{
			ast.IfBlock{
				Test: ast.BinaryExpression{
					Left: ast.Literal{
						Node: ast.Node{
							Type:  "Literal",
							Start: 3,
							End:   4,
						},
						Value: 38,
					},
					Right: ast.Literal{
						Node: ast.Node{
							Type:  "Literal",
							Start: 9,
							End:   11,
						},
						Value: 121,
					},
					Operator: "<=",
					Node: ast.Node{
						Start: 3,
						End:   11,
						Type:  "BinaryExpression",
					},
				},
				Consequent: ast.Literal{
					Node: ast.Node{
						Start: 16,
						End:   16,
						Type:  "Literal",
					},
					Value: 1,
				},
				Node: ast.Node{
					Start: 0,
					End:   16,
					Type:  "IfExpression",
				},
			},
		}, input: `if 38 <= 121 do 1`},
		{name: "if else block", want: []ast.Statement{
			ast.IfElseBlock{
				Test: ast.BinaryExpression{
					Left: ast.Literal{
						Node: ast.Node{
							Type:  "Literal",
							Start: 3,
							End:   4,
						},
						Value: 38,
					},
					Right: ast.Literal{
						Node: ast.Node{
							Type:  "Literal",
							Start: 9,
							End:   11,
						},
						Value: 121,
					},
					Operator: "<=",
					Node: ast.Node{
						Start: 3,
						End:   11,
						Type:  "BinaryExpression",
					},
				},
				Consequent: ast.Literal{
					Node: ast.Node{
						Start: 16,
						End:   16,
						Type:  "Literal",
					},
					Value: 1,
				},
				Alternate: ast.Literal{
					Node: ast.Node{
						Start: 23,
						End:   23,
						Type:  "Literal",
					},
					Value: 2,
				},
				Node: ast.Node{
					Start: 0,
					End:   23,
					Type:  "IfElseExpression",
				},
			},
		}, input: `if 38 <= 121 do 1 else 2`},
		{name: "if else block with true keyword", want: []ast.Statement{
			ast.IfElseBlock{
				Test: ast.Identifier{
					Node: ast.Node{
						Start: 3,
						End:   6,
						Type:  "Identifier",
					},
					Value: "true",
				},
				Consequent: ast.Literal{
					Node: ast.Node{
						Start: 11,
						End:   11,
						Type:  "Literal",
					},
					Value: 1,
				},
				Alternate: ast.Literal{
					Node: ast.Node{
						Start: 18,
						End:   18,
						Type:  "Literal",
					},
					Value: 2,
				},
				Node: ast.Node{
					Start: 0,
					End:   18,
					Type:  "IfElseExpression",
				},
			},
		}, input: `if true do 1 else 2`},
		{name: "function declaration", want: []ast.Statement{
			ast.FunctionExpression{
				Body: []ast.Statement{
					ast.Identifier{
						Node: ast.Node{
							Start: 19,
							End:   19,
							Type:  "Identifier",
						},
						Value: "a",
					},
				},
				Node: ast.Node{
					Start: 0,
					End:   23,
					Type:  "FunctionExpression",
				},
				Name: ast.Identifier{
					Node: ast.Node{
						Start: 3,
						End:   7,
						Type:  "Identifier",
					},
					Value: "hello",
				},

				Parameters: []ast.Identifier{
					{
						Node: ast.Node{
							Start: 10,
							End:   10,
							Type:  "Identifier",
						},
						Value: "a",
					},
					{
						Node: ast.Node{
							Start: 13,
							End:   13,
							Type:  "Identifier",
						},
						Value: "b",
					},
				},
			},
		}, input: `fn hello |a, b| -> a end`},
		{name: "function evaluation", want: []ast.Statement{
			ast.FunctionEvaluation{
				Node: ast.Node{
					Start: 0,
					End:   10,
					Type:  "FunctionEvaluation",
				},
				Name: ast.Identifier{
					Node: ast.Node{
						Start: 0,
						End:   4,
						Type:  "Identifier",
					},
					Value: "hello",
				},
				Arguments: []ast.Statement{
					ast.Identifier{
						Node: ast.Node{
							Start: 6,
							End:   6,
							Type:  "Identifier",
						},
						Value: "a",
					},
					ast.Identifier{
						Node: ast.Node{
							Start: 9,
							End:   9,
							Type:  "Identifier",
						},
						Value: "b",
					},
				},
			}},
			input: `hello(a, b)`},
		{name: "function declaration with multiple statements", want: []ast.Statement{
			ast.FunctionExpression{
				Body: []ast.Statement{
					ast.Identifier{
						Node: ast.Node{
							Start: 19,
							End:   19,
							Type:  "Identifier",
						},
						Value: "a",
					},
					ast.Identifier{
						Node: ast.Node{
							Start: 21,
							End:   21,
							Type:  "Identifier",
						},
						Value: "b",
					},
				},
				Node: ast.Node{
					Start: 0,
					End:   25,
					Type:  "FunctionExpression",
				},
				Name: ast.Identifier{
					Node: ast.Node{
						Start: 3,
						End:   7,
						Type:  "Identifier",
					},
					Value: "hello",
				},

				Parameters: []ast.Identifier{
					{
						Node: ast.Node{
							Start: 10,
							End:   10,
							Type:  "Identifier",
						},
						Value: "a",
					},
					{
						Node: ast.Node{
							Start: 13,
							End:   13,
							Type:  "Identifier",
						},
						Value: "b",
					},
				},
			},
		}, input: `fn hello |a, b| ->
a
b
end`},
		{name: "let declaration", want: []ast.Statement{
			ast.LetStatement{
				Left: ast.Identifier{
					Value: "incr",
					Node: ast.Node{
						Start: 4,
						End:   7,
						Type:  "Identifier",
					},
				},
				Right: ast.FunctionExpression{
					Body: []ast.Statement{
						ast.BinaryExpression{
							Left: ast.Identifier{
								Node: ast.Node{
									Start: 21,
									End:   21,
									Type:  "Identifier",
								},
								Value: "a",
							},
							Right: ast.Literal{
								Value: 1,
								Node: ast.Node{
									Start: 25,
									End:   25,
									Type:  "Literal",
								},
							},
							Node: ast.Node{
								Start: 21,
								End:   25,
								Type:  "BinaryExpression",
							},
							Operator: "+",
						},
					},
					Parameters: []ast.Identifier{
						{
							Node: ast.Node{
								Start: 15,
								End:   15,
								Type:  "Identifier",
							},
							Value: "a",
						},
					},
					Node: ast.Node{
						Start: 11,
						End:   29,
						Type:  "FunctionExpression",
					},
					Name: ast.Identifier{},
				},
				Operator: "=",
				Node: ast.Node{
					Start: 0,
					End:   29,
					Type:  "LetStatement",
				},
			},
		}, input: `let incr = fn |a| -> a + 1 end`},
		{name: "let declaration in multiple lines", want: []ast.Statement{
			ast.LetStatement{
				Left: ast.Identifier{
					Value: "incr",
					Node: ast.Node{
						Start: 4,
						End:   7,
						Type:  "Identifier",
					},
				},
				Right: ast.FunctionExpression{
					Body: []ast.Statement{
						ast.BinaryExpression{
							Left: ast.Identifier{
								Node: ast.Node{
									Start: 21,
									End:   21,
									Type:  "Identifier",
								},
								Value: "a",
							},
							Right: ast.Literal{
								Value: 1,
								Node: ast.Node{
									Start: 25,
									End:   25,
									Type:  "Literal",
								},
							},
							Node: ast.Node{
								Start: 21,
								End:   25,
								Type:  "BinaryExpression",
							},
							Operator: "+",
						},
					},
					Parameters: []ast.Identifier{
						{
							Node: ast.Node{
								Start: 15,
								End:   15,
								Type:  "Identifier",
							},
							Value: "a",
						},
					},
					Node: ast.Node{
						Start: 11,
						End:   29,
						Type:  "FunctionExpression",
					},
					Name: ast.Identifier{},
				},
				Operator: "=",
				Node: ast.Node{
					Start: 0,
					End:   29,
					Type:  "LetStatement",
				},
			},
		}, input: `let incr = fn |a| ->
a + 1
end`},
		{name: "logical AND", want: []ast.Statement{
			ast.BinaryExpression{
				Left: ast.Identifier{
					Node: ast.Node{
						Type:  "Identifier",
						Start: 0,
						End:   3,
					},
					Value: "true",
				},
				Right: ast.Identifier{
					Node: ast.Node{
						Type:  "Identifier",
						Start: 9,
						End:   13,
					},
					Value: "false",
				},
				Operator: "and",
				Node: ast.Node{
					Start: 0,
					End:   13,
					Type:  "BinaryExpression",
				},
			},
		}, input: "true and false"},
		{name: "logical binary expression", want: []ast.Statement{
			ast.BinaryExpression{
				Left: ast.BinaryExpression{
					Left: ast.Literal{
						Node: ast.Node{
							Start: 0,
							End:   1,
							Type:  "Literal",
						},
						Value: 10,
					},
					Right: ast.Literal{
						Node: ast.Node{
							Start: 5,
							End:   6,
							Type:  "Literal",
						},
						Value: 12,
					},
					Operator: "<",
					Node: ast.Node{
						Start: 0,
						End:   6,
						Type:  "BinaryExpression",
					},
				},
				Right: ast.BinaryExpression{
					Left: ast.Literal{
						Node: ast.Node{
							Start: 11,
							End:   12,
							Type:  "Literal",
						},
						Value: 14,
					},
					Right: ast.Literal{
						Node: ast.Node{
							Start: 16,
							End:   17,
							Type:  "Literal",
						},
						Value: 34,
					},
					Operator: ">",
					Node: ast.Node{
						Start: 11,
						End:   17,
						Type:  "BinaryExpression",
					},
				},
				Operator: "or",
				Node: ast.Node{
					Start: 0,
					End:   17,
					Type:  "BinaryExpression",
				},
			},
		}, input: "10 < 12 or 14 > 34"},
		{name: "function evaluation (left side) in binary expression", want: []ast.Statement{
			ast.BinaryExpression{
				Left: ast.FunctionEvaluation{
					Node: ast.Node{
						Type:  "FunctionEvaluation",
						Start: 0,
						End:   6,
					},
					Arguments: []ast.Statement{
						ast.Literal{
							Node: ast.Node{
								Start: 5,
								End:   5,
								Type:  "Literal",
							},
							Value: 1,
						},
					},
					Name: ast.Identifier{
						Node: ast.Node{
							Start: 0,
							End:   3,
							Type:  "Identifier",
						},
						Value: "incr",
					},
				},
				Right: ast.Literal{
					Node: ast.Node{
						Type:  "Literal",
						Start: 10,
						End:   11,
					},
					Value: 13,
				},
				Operator: "+",
				Node: ast.Node{
					Start: 0,
					End:   11,
					Type:  "BinaryExpression",
				},
			},
		}, input: "incr(1) + 13"},
		{name: "function evaluation (right side) in binary expression", want: []ast.Statement{
			ast.BinaryExpression{
				Right: ast.FunctionEvaluation{
					Node: ast.Node{
						Type:  "FunctionEvaluation",
						Start: 5,
						End:   11,
					},
					Arguments: []ast.Statement{
						ast.Literal{
							Node: ast.Node{
								Start: 10,
								End:   10,
								Type:  "Literal",
							},
							Value: 1,
						},
					},
					Name: ast.Identifier{
						Node: ast.Node{
							Start: 5,
							End:   8,
							Type:  "Identifier",
						},
						Value: "incr",
					},
				},
				Left: ast.Literal{
					Node: ast.Node{
						Type:  "Literal",
						Start: 0,
						End:   1,
					},
					Value: 13,
				},
				Operator: "+",
				Node: ast.Node{
					Start: 0,
					End:   11,
					Type:  "BinaryExpression",
				},
			},
		}, input: "13 + incr(1)"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := lexer.New([]rune(tt.input))
			parser := Parser{
				lexer:        lexer,
				currentToken: lexer.NextToken(),
				peekToken:    lexer.NextToken(),
			}
			result := parser.Parse()
			for i := range result.Body {
				if !reflect.DeepEqual(result.Body[i], tt.want[i]) {
					t.Errorf("got %+v, want %+v", result.Body[i], tt.want[i])
				}
			}
		})
	}
}
