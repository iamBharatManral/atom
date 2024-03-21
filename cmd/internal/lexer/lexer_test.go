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
			token.New(token.EOF, "", "", 4, 4),
		}, input: "1933"},
		{name: "float literal 145.6", want: []token.Token{
			token.New(token.FLOAT, "", 145.6, 0, 4),
			token.New(token.EOF, "", "", 5, 5),
		}, input: "145.6"},
		{name: "string literal hello", want: []token.Token{
			token.New(token.STRING, "\"hello\"", "hello", 0, 6),
			token.New(token.EOF, "", "", 7, 7),
		}, input: "\"hello\""},
		{name: "all the literals", want: []token.Token{
			token.New(token.STRING, "\"hello\"", "hello", 0, 6),
			token.New(token.INTEGER, "", 123, 8, 10),
			token.New(token.FLOAT, "", 123.4, 12, 16),
			token.New(token.INTEGER, "", 45678, 18, 22),
			token.New(token.EOF, "", "", 23, 23),
		}, input: "\"hello\" 123 123.4 45678"},
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
		{name: "identifier let", want: []token.Token{
			token.New(token.IDENTIFIER, "let", "", 0, 2),
			token.New(token.EOF, "", "", 4, 4),
		}, input: "let "},
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
		}, input: "name = \"hello\""},
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
