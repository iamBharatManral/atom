package token

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	STRING  = "STRING"
	INTEGER = "INTEGER"
	FLOAT   = "FLOAT"
)

type Token struct {
	literal   any
	lexeme    string
	tokenType string
	start     uint
	end       uint
}

func New(tokenType string, lexeme string, literal any, start uint, end uint) Token {
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

func (t Token) Start() uint {
	return t.start
}

func (t Token) End() uint {
	return t.end
}
