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
			},
		}, input: "1933;"},
		{name: "string literal \"word is my oyster\"", want: []ast.Statement{
			ast.Literal{
				Value: "word is my oyster",
				Node: ast.Node{
					Type:  "Literal",
					Start: 0,
					End:   18,
				},
			},
		}, input: "\"word is my oyster\";"},
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
		}, input: "12 + 13;"},
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
						Start: 7,
						End:   8,
					},
					Value: 96,
				},
				Right: ast.Literal{
					Node: ast.Node{
						Type:  "Literal",
						Start: 12,
						End:   12,
					},
					Value: 4,
				},
				Operator: "/",
				Node: ast.Node{
					Start: 7,
					End:   12,
					Type:  "BinaryExpression",
				},
			},
		}, input: `5 * 9;
96 / 4;
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
					Start: 7,
					End:   19,
					Type:  "Literal",
				},
				Value: "hello world",
			},
			ast.Literal{
				Node: ast.Node{
					Start: 22,
					End:   22,
					Type:  "Literal",
				},
				Value: 1,
			},
			ast.BinaryExpression{
				Left: ast.Literal{
					Node: ast.Node{
						Type:  "Literal",
						Start: 25,
						End:   25,
					},
					Value: 4,
				},
				Right: ast.Literal{
					Node: ast.Node{
						Type:  "Literal",
						Start: 28,
						End:   28,
					},
					Value: 5,
				},
				Operator: "*",
				Node: ast.Node{
					Start: 25,
					End:   28,
					Type:  "BinaryExpression",
				},
			},
		}, input: `5 * 9;
"hello world";
1;
4 *5;
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
		}, input: `let a = 10;`},
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
		}, input: `name = "hello";`},
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
		}, input: `51 + 23 * 4;`},
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
		}, input: "12 < 13;"},
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
		}, input: "\"hello\" >= \"bye\";"},
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
		}, input: `if 38 <= 121 do 1;`},
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
						Start: 24,
						End:   24,
						Type:  "Literal",
					},
					Value: 2,
				},
				Node: ast.Node{
					Start: 0,
					End:   24,
					Type:  "IfElseExpression",
				},
			},
		}, input: `if 38 <= 121 do 1; else 2;`},
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
						Start: 19,
						End:   19,
						Type:  "Literal",
					},
					Value: 2,
				},
				Node: ast.Node{
					Start: 0,
					End:   19,
					Type:  "IfElseExpression",
				},
			},
		}, input: `if true do 1; else 2;`},
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
