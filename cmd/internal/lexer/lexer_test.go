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
		{name: "integer literal 1933", want: []token.Token{
			token.New(token.INTEGER, "", 1933, 0, 3),
			token.New(token.SEMICOLON, ";", "", 4, 4),
			token.New(token.EOF, "", "", 5, 5),
		}, input: "1933;"},
		{name: "float literal 145.6", want: []token.Token{
			token.New(token.FLOAT, "", 145.6, 0, 4),
			token.New(token.SEMICOLON, ";", "", 5, 5),
			token.New(token.EOF, "", "", 6, 6),
		}, input: "145.6;"},
		{name: "string literal hello", want: []token.Token{
			token.New(token.STRING, "\"hello\"", "hello", 0, 6),
			token.New(token.SEMICOLON, ";", "", 7, 7),
			token.New(token.EOF, "", "", 8, 8),
		}, input: "\"hello\";"},
		{name: "all the literals", want: []token.Token{
			token.New(token.STRING, "\"hello\"", "hello", 0, 6),
			token.New(token.INTEGER, "", 123, 8, 10),
			token.New(token.FLOAT, "", 123.4, 12, 16),
			token.New(token.INTEGER, "", 45678, 18, 22),
			token.New(token.SEMICOLON, ";", "", 23, 23),
			token.New(token.EOF, "", "", 24, 24),
		}, input: "\"hello\" 123 123.4 45678;"},
		{name: "binary expression 2 + 3", want: []token.Token{
			token.New(token.INTEGER, "", 2, 0, 0),
			token.New(token.PLUS, "+", "", 2, 2),
			token.New(token.INTEGER, "", 3, 4, 4),
			token.New(token.SEMICOLON, ";", "", 5, 5),
			token.New(token.EOF, "", "", 6, 6),
		}, input: "2 + 3;"},
		{name: "binary expressions without space 2+3", want: []token.Token{
			token.New(token.INTEGER, "", 2, 0, 0),
			token.New(token.PLUS, "+", "", 1, 1),
			token.New(token.INTEGER, "", 3, 2, 2),
			token.New(token.SEMICOLON, ";", "", 3, 3),
			token.New(token.EOF, "", "", 4, 4),
		}, input: "2+3;"},
		{name: "identifier let", want: []token.Token{
			token.New(token.IDENTIFIER, "let", "", 0, 2),
			token.New(token.SEMICOLON, ";", "", 4, 4),
			token.New(token.EOF, "", "", 5, 5),
		}, input: "let ;"},
		{name: "let statement: let id = 100", want: []token.Token{
			token.New(token.IDENTIFIER, "let", "", 0, 2),
			token.New(token.IDENTIFIER, "id", "", 4, 5),
			token.New(token.ASSIGN, "=", "", 7, 7),
			token.New(token.INTEGER, "", 100, 9, 11),
			token.New(token.SEMICOLON, ";", "", 12, 12),
			token.New(token.EOF, "", "", 13, 13),
		}, input: "let id = 100;"},
		{name: "assignment operation: name = \"hello\"", want: []token.Token{
			token.New(token.IDENTIFIER, "name", "", 0, 3),
			token.New(token.ASSIGN, "=", "", 5, 5),
			token.New(token.STRING, "\"hello\"", "hello", 7, 13),
			token.New(token.SEMICOLON, ";", "", 14, 14),
			token.New(token.EOF, "", "", 15, 15),
		},
			input: "name = \"hello\";"},
		{name: "multiple expressios", want: []token.Token{
			token.New(token.INTEGER, "", 2, 0, 0),
			token.New(token.PLUS, "+", "", 2, 2),
			token.New(token.INTEGER, "", 3, 4, 4),
			token.New(token.PLUS, "+", "", 6, 6),
			token.New(token.INTEGER, "", 5, 8, 8),
			token.New(token.STAR, "*", "", 10, 10),
			token.New(token.INTEGER, "", 6, 11, 11),
			token.New(token.SEMICOLON, ";", "", 12, 12),
			token.New(token.EOF, "", "", 13, 13),
		}, input: "2 + 3 + 5 *6;"},
		{name: "multiple statements", want: []token.Token{
			token.New(token.INTEGER, "", 3, 0, 0),
			token.New(token.STAR, "*", "", 2, 2),
			token.New(token.INTEGER, "", 5, 4, 4),
			token.New(token.SEMICOLON, ";", "", 5, 5),
			token.New(token.INTEGER, "", 5, 7, 7),
			token.New(token.PLUS, "+", "", 9, 9),
			token.New(token.INTEGER, "", 9, 11, 11),
			token.New(token.SEMICOLON, ";", "", 12, 12),
			token.New(token.EOF, "", "", 13, 13),
		}, input: `3 * 5;
5 + 9;`},
		{name: "conditional operator", want: []token.Token{
			token.New(token.INTEGER, "", 3, 0, 0),
			token.New(token.LE, "<=", "", 2, 3),
			token.New(token.INTEGER, "", 5, 5, 5),
			token.New(token.SEMICOLON, ";", "", 6, 6),
			token.New(token.EOF, "", "", 7, 7),
		}, input: `3 <= 5;`},
		{name: "if condition", want: []token.Token{
			token.New(token.IDENTIFIER, "if", "", 0, 1),
			token.New(token.INTEGER, "", 10, 3, 4),
			token.New(token.LT, "<", "", 6, 6),
			token.New(token.INTEGER, "", 12, 8, 9),
			token.New(token.IDENTIFIER, "do", "", 11, 12),
			token.New(token.STRING, "\"true\"", "true", 14, 19),
			token.New(token.SEMICOLON, ";", "", 20, 20),
			token.New(token.EOF, "", "", 21, 21),
		}, input: `if 10 < 12 do "true";`},
		{name: "if else condition", want: []token.Token{
			token.New(token.IDENTIFIER, "if", "", 0, 1),
			token.New(token.INTEGER, "", 10, 3, 4),
			token.New(token.LT, "<", "", 6, 6),
			token.New(token.INTEGER, "", 12, 8, 9),
			token.New(token.IDENTIFIER, "do", "", 11, 12),
			token.New(token.STRING, "\"true\"", "true", 14, 19),
			token.New(token.SEMICOLON, ";", "", 20, 20),
			token.New(token.IDENTIFIER, "else", "", 22, 25),
			token.New(token.STRING, "\"false\"", "false", 27, 33),
			token.New(token.SEMICOLON, ";", "", 34, 34),
			token.New(token.EOF, "", "", 35, 35),
		}, input: `if 10 < 12 do "true"; else "false";`},
		{name: "if else conditon with true keyword", want: []token.Token{
			token.New(token.IDENTIFIER, "if", "", 0, 1),
			token.New(token.IDENTIFIER, "true", "", 3, 6),
			token.New(token.IDENTIFIER, "do", "", 8, 9),
			token.New(token.STRING, "\"true\"", "true", 11, 16),
			token.New(token.SEMICOLON, ";", "", 17, 17),
			token.New(token.IDENTIFIER, "else", "", 19, 22),
			token.New(token.STRING, "\"false\"", "false", 24, 30),
			token.New(token.SEMICOLON, ";", "", 31, 31),
			token.New(token.EOF, "", "", 32, 32),
		}, input: `if true do "true"; else "false";`},
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
			token.New(token.SEMICOLON, ";", "", 20, 20),
			token.New(token.IDENTIFIER, "end", "", 22, 24),
			token.New(token.SEMICOLON, ";", "", 25, 25),
			token.New(token.EOF, "", "", 26, 26),
		}, input: `fn hello |a, b| -> a; end;`},
		{name: "function evaluation", want: []token.Token{
			token.New(token.IDENTIFIER, "hello", "", 0, 4),
			token.New(token.LPAREN, "(", "", 5, 5),
			token.New(token.IDENTIFIER, "a", "", 6, 6),
			token.New(token.COMMA, ",", "", 7, 7),
			token.New(token.IDENTIFIER, "b", "", 9, 9),
			token.New(token.RPAREN, ")", "", 10, 10),
			token.New(token.SEMICOLON, ";", "", 11, 11),
			token.New(token.EOF, "", "", 12, 12),
		}, input: `hello(a, b);`},
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
