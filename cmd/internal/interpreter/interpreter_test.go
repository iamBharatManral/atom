package interpreter

import (
	"testing"

	"github.com/iamBharatManral/atom.git/cmd/internal/env"
	"github.com/iamBharatManral/atom.git/cmd/internal/lexer"
	"github.com/iamBharatManral/atom.git/cmd/internal/parser"
)

func TestEvaluation(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  any
	}{
		{name: "single digit expression", want: 1, input: `1;`},
		{name: "binary expression", want: 35, input: `12+23;`},
		{name: "variable declaration", want: nil, input: `let a = 10;`},
		{name: "comparison between integers", want: true, input: `13 > 9;`},
		{name: "equality check", want: true, input: `"hello" == "hello";`},
		{name: "comparison between floats", want: true, input: `1.2 <= 3.4;`},
		{name: "comparison between strings", want: false, input: `"greater" > "less";`},
		{name: "if block", want: "greater", input: `if 12 > 10 do "greater";`},
		{name: "if else block with truthy condition", want: "greater than", input: `if 12 > 10 do "greater than"; else "less than";`},
		{name: "if else block with falsy condition", want: "false", input: `if 10 != 10 do "true"; else "false";`},
		{name: "if else block with true keyword", want: "true", input: `if true do "true";`},
		{name: "if else block with false keyword", want: "false", input: `if false do "true"; else "false";`},
		{name: "function declaration", want: "()", input: `fn hello|a,b| -> a; end;`},
		{name: "binary expression with logical and", want: false, input: `10 != 10 and 12 > 10;`},
		{name: "binary expression with logical or", want: true, input: `10 > 10 or 10 != 10 or 12 > 7;`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := lexer.New([]rune(tt.input))
			parser := parser.New(lexer)
			program := parser.Parse()
			env := env.New(nil)
			for i := range program.Body {
				output := Eval(program.Body[i], env)
				if output.Value != tt.want {
					t.Errorf("got %+v, want %+v", output, tt.want)
				}
			}
		})
	}

}
