package token

const (
	ILLEGAL   = "ILLEGAL"
	EOF       = "EOF"
	SEMICOLON = "SEMICOLON"
	PLUS      = "PLUS"
	MINUS     = "MINUS"
	STAR      = "STAR"
	SLASH     = "SLASH"

	STRING  = "STRING"
	INTEGER = "INTEGER"
	FLOAT   = "FLOAT"

	LE = "LE"
	EQ = "EQ"
	GT = "GT"
	GE = "GE"
	LT = "LT"
	NE = "NE"

	NOT   = "NOT"
	AND   = "AND"
	OR    = "OR"
	TRUE  = "TRUE"
	FALSE = "FALSE"

	ARROW  = "ARROW"
	BAR    = "BAR"
	COMMA  = "COMMA"
	LPAREN = "LPAREN"
	RPAREN = "RPAREN"

	IDENTIFIER = "IDENTIFIER"
	ASSIGN     = "ASSIGN"
)

var keywords = make(map[string]string)
var priorities = make(map[string]int)

func RegisterPriorities() {
	priorities["NONE"] = 0
	priorities["PLUS"] = 1
	priorities["MINUS"] = 1
	priorities["STAR"] = 2
	priorities["SLASH"] = 2
	priorities["AND"] = 4
	priorities["OR"] = 5
	priorities["NOT"] = 6

	priorities["LE"] = 7
	priorities["NE"] = 7
	priorities["LT"] = 7
	priorities["GT"] = 7
	priorities["GE"] = 7
	priorities["EQ"] = 7
}

func GetPriority(token string) int {
	if priority, ok := priorities[token]; ok {
		return priority
	}
	return priorities["NONE"]

}
func RegisterKeyWords() {
	keywords["let"] = "let"
	keywords["if"] = "if"
	keywords["do"] = "do"
	keywords["else"] = "else"
	keywords["false"] = "false"
	keywords["true"] = "true"
	keywords["fn"] = "fn"
}

func IsKeyword(key string) bool {
	return keywords[key] == key
}

func GetKeyword(key string) string {
	return keywords[key]
}

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
