package token

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	PLUS  = "PLUS"
	MINUS = "MINUS"
	STAR  = "STAR"
	SLASH = "SLASH"

	STRING  = "STRING"
	INTEGER = "INTEGER"
	FLOAT   = "FLOAT"

	IDENTIFIER = "IDENTIFIER"
)

type Token struct {
	literal   any
	lexeme    string
	tokenType string
	start     int
	end       int
}

func New(tokenType string, lexeme string, literal any, start int, end int) Token {
	return Token{
		literal,
		lexeme,
		tokenType,
		start,
		end,
	}
}

func (t Token) TokenType() string {
	return t.tokenType
}

func (t Token) Value() any {
	return t.literal
}

func (t Token) Lexeme() string {
	return t.lexeme
}

func (t Token) Start() int {
	return t.start
}

func (t Token) End() int {
	return t.end
}
