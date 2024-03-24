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
		{name: "comparison between floats", want: true, input: `1.2 <= 3.4;`},
		{name: "comparison between strings", want: false, input: `"greater" > "less";`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := lexer.New([]rune(tt.input))
			parser := parser.New(lexer)
			program := parser.Parse()
			env := env.New()
			for i := range program.Body {
				output := Eval(program.Body[i], env)
				if output.Value != tt.want {
					t.Errorf("got %+v, want %+v", output, tt.want)
				}
			}
		})
	}

}
