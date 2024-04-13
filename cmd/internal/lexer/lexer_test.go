package lexer

import (
	"reflect"
	"testing"

	"github.com/iamBharatManral/atom.git/cmd/internal/token"
)

func TestTokens(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []token.Token
	}{
		{name: "integer", want: []token.Token{
			token.New(token.INTEGER, "", 1933, 0, 3),
			token.New(token.EOF, "", "", 4, 4),
		}, input: "1933"},
		{name: "float", want: []token.Token{
			token.New(token.FLOAT, "", 145.6, 0, 4),
			token.New(token.EOF, "", "", 5, 5),
		}, input: "145.6"},
		{name: "string", want: []token.Token{
			token.New(token.STRING, "\"hello\"", "hello", 0, 6),
			token.New(token.EOF, "", "", 7, 7),
		}, input: "\"hello\""},
		{name: "boolean true", want: []token.Token{
			token.New(token.IDENTIFIER, "true", "", 0, 3),
			token.New(token.EOF, "", "", 4, 4),
		}, input: "true"},
		{name: "boolean false", want: []token.Token{
			token.New(token.IDENTIFIER, "false", "", 0, 4),
			token.New(token.EOF, "", "", 5, 5),
		}, input: "false"},
		{name: "unary operator minus", want: []token.Token{
			token.New(token.LPAREN, "(", "", 0, 0),
			token.New(token.MINUS, "-", "", 1, 1),
			token.New(token.INTEGER, "", 1, 2, 2),
			token.New(token.RPAREN, ")", "", 3, 3),
			token.New(token.EOF, "", "", 4, 4),
		}, input: "(-1)"},

		{name: "not operator", want: []token.Token{
			token.New(token.NOT, "!", "", 0, 0),
			token.New(token.IDENTIFIER, "false", "", 2, 6),
			token.New(token.EOF, "", "", 7, 7),
		}, input: "! false"},
		{name: "binary expression 2 + 3", want: []token.Token{
			token.New(token.INTEGER, "", 2, 0, 0),
			token.New(token.PLUS, "+", "", 2, 2),
			token.New(token.INTEGER, "", 3, 4, 4),
			token.New(token.EOF, "", "", 5, 5),
		}, input: "2 + 3"},
		{name: "binary expressions without space 2+3", want: []token.Token{
			token.New(token.INTEGER, "", 2, 0, 0),
			token.New(token.PLUS, "+", "", 1, 1),
			token.New(token.INTEGER, "", 3, 2, 2),
			token.New(token.EOF, "", "", 3, 3),
		}, input: "2+3"},
		{name: "binary expressions with brackets (2+3)", want: []token.Token{
			token.New(token.LPAREN, "(", "", 0, 0),
			token.New(token.INTEGER, "", 2, 1, 1),
			token.New(token.PLUS, "+", "", 2, 2),
			token.New(token.INTEGER, "", 3, 3, 3),
			token.New(token.RPAREN, ")", "", 4, 4),
			token.New(token.EOF, "", "", 5, 5),
		}, input: "(2+3)"},

		{name: "multiple binary expressions with brackets 15 * (2+3)", want: []token.Token{
			token.New(token.INTEGER, "", 15, 0, 1),
			token.New(token.STAR, "*", "", 3, 3),
			token.New(token.LPAREN, "(", "", 5, 5),
			token.New(token.INTEGER, "", 2, 6, 6),
			token.New(token.PLUS, "+", "", 7, 7),
			token.New(token.INTEGER, "", 3, 8, 8),
			token.New(token.RPAREN, ")", "", 9, 9),
			token.New(token.EOF, "", "", 10, 10),
		}, input: "15 * (2+3)"},

		{name: "identifier let", want: []token.Token{
			token.New(token.IDENTIFIER, "let", "", 0, 2),
			token.New(token.EOF, "", "", 3, 3),
		}, input: "let"},
		{name: "let statement: let id = 100", want: []token.Token{
			token.New(token.IDENTIFIER, "let", "", 0, 2),
			token.New(token.IDENTIFIER, "id", "", 4, 5),
			token.New(token.ASSIGN, "=", "", 7, 7),
			token.New(token.INTEGER, "", 100, 9, 11),
			token.New(token.EOF, "", "", 12, 12),
		}, input: "let id = 100"},
		{name: "assignment operation: name = \"hello\"", want: []token.Token{
			token.New(token.IDENTIFIER, "name", "", 0, 3),
			token.New(token.ASSIGN, "=", "", 5, 5),
			token.New(token.STRING, "\"hello\"", "hello", 7, 13),
			token.New(token.EOF, "", "", 14, 14),
		},
			input: "name = \"hello\""},
		{name: "binary expressions", want: []token.Token{
			token.New(token.INTEGER, "", 2, 0, 0),
			token.New(token.PLUS, "+", "", 2, 2),
			token.New(token.INTEGER, "", 3, 4, 4),
			token.New(token.PLUS, "+", "", 6, 6),
			token.New(token.INTEGER, "", 5, 8, 8),
			token.New(token.STAR, "*", "", 10, 10),
			token.New(token.INTEGER, "", 6, 11, 11),
			token.New(token.EOF, "", "", 12, 12),
		}, input: "2 + 3 + 5 *6"},
		{name: "multiple statements", want: []token.Token{
			token.New(token.INTEGER, "", 3, 0, 0),
			token.New(token.STAR, "*", "", 2, 2),
			token.New(token.INTEGER, "", 5, 4, 4),
			token.New(token.NEWLINE, "", "", 5, 5),
			token.New(token.INTEGER, "", 5, 6, 6),
			token.New(token.PLUS, "+", "", 8, 8),
			token.New(token.INTEGER, "", 9, 10, 10),
			token.New(token.EOF, "", "", 11, 11),
		}, input: `3 * 5
5 + 9`},
		{name: "conditional operator", want: []token.Token{
			token.New(token.INTEGER, "", 3, 0, 0),
			token.New(token.LE, "<=", "", 2, 3),
			token.New(token.INTEGER, "", 5, 5, 5),
			token.New(token.EOF, "", "", 6, 6),
		}, input: `3 <= 5`},
		{name: "if condition", want: []token.Token{
			token.New(token.IDENTIFIER, "if", "", 0, 1),
			token.New(token.INTEGER, "", 10, 3, 4),
			token.New(token.LT, "<", "", 6, 6),
			token.New(token.INTEGER, "", 12, 8, 9),
			token.New(token.IDENTIFIER, "do", "", 11, 12),
			token.New(token.STRING, "\"true\"", "true", 14, 19),
			token.New(token.EOF, "", "", 20, 20),
		}, input: `if 10 < 12 do "true"`},
		{name: "if else condition", want: []token.Token{
			token.New(token.IDENTIFIER, "if", "", 0, 1),
			token.New(token.INTEGER, "", 10, 3, 4),
			token.New(token.LT, "<", "", 6, 6),
			token.New(token.INTEGER, "", 12, 8, 9),
			token.New(token.IDENTIFIER, "do", "", 11, 12),
			token.New(token.STRING, "\"true\"", "true", 14, 19),
			token.New(token.IDENTIFIER, "else", "", 21, 24),
			token.New(token.STRING, "\"false\"", "false", 26, 32),
			token.New(token.EOF, "", "", 33, 33),
		}, input: `if 10 < 12 do "true" else "false"`},
		{name: "if else conditon with true keyword", want: []token.Token{
			token.New(token.IDENTIFIER, "if", "", 0, 1),
			token.New(token.IDENTIFIER, "true", "", 3, 6),
			token.New(token.IDENTIFIER, "do", "", 8, 9),
			token.New(token.STRING, "\"true\"", "true", 11, 16),
			token.New(token.IDENTIFIER, "else", "", 18, 21),
			token.New(token.STRING, "\"false\"", "false", 23, 29),
			token.New(token.EOF, "", "", 30, 30),
		}, input: `if true do "true" else "false"`},
		{name: "function declaration", want: []token.Token{
			token.New(token.IDENTIFIER, "fn", "", 0, 1),
			token.New(token.IDENTIFIER, "hello", "", 3, 7),
			token.New(token.BAR, "|", "", 9, 9),
			token.New(token.IDENTIFIER, "a", "", 10, 10),
			token.New(token.COMMA, ",", "", 11, 11),
			token.New(token.IDENTIFIER, "b", "", 13, 13),
			token.New(token.BAR, "|", "", 14, 14),
			token.New(token.ARROW, "->", "", 16, 17),
			token.New(token.IDENTIFIER, "a", "", 19, 19),
			token.New(token.IDENTIFIER, "end", "", 21, 23),
			token.New(token.EOF, "", "", 24, 24),
		}, input: `fn hello |a, b| -> a end`},
		{name: "function evaluation", want: []token.Token{
			token.New(token.IDENTIFIER, "hello", "", 0, 4),
			token.New(token.LPAREN, "(", "", 5, 5),
			token.New(token.IDENTIFIER, "a", "", 6, 6),
			token.New(token.COMMA, ",", "", 7, 7),
			token.New(token.IDENTIFIER, "b", "", 9, 9),
			token.New(token.RPAREN, ")", "", 10, 10),
			token.New(token.EOF, "", "", 11, 11),
		}, input: `hello(a, b)`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ans []token.Token
			lexer := New([]rune(tt.input))
			tok := lexer.NextToken()
			for tok.TokenType() != token.EOF {
				ans = append(ans, tok)
				tok = lexer.NextToken()
			}
			ans = append(ans, tok)
			if !reflect.DeepEqual(ans, tt.want) {
				t.Errorf("got %+v, want %+v", ans, tt.want)
			}
		})
	}
}
