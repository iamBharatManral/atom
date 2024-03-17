package lexer

import (
	"reflect"
	"testing"

	"github.com/iamBharatManral/atom.git/cmd/internal/token"
)

func TestLiterals(t *testing.T) {
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
			token.New(token.EOF, "", "", 33, 23),
		}, input: "\"hello\" 123 123.4 45678"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ans []token.Token
			lexer := Lexer{input: []rune(tt.input)}
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
