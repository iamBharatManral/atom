package token

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"
	NEWLINE = "NEWLINE"

	PLUS  = "PLUS"
	MINUS = "MINUS"
	STAR  = "STAR"
	SLASH = "SLASH"
	MOD   = "MOD"

	STRING  = "STRING"
	INTEGER = "INTEGER"
	FLOAT   = "FLOAT"

	LE     = "LE"
	EQ     = "EQ"
	GT     = "GT"
	GE     = "GE"
	LT     = "LT"
	NE     = "NE"
	ASSIGN = "ASSIGN"

	NOT = "NOT"

	ARROW  = "ARROW"
	BAR    = "BAR"
	COMMA  = "COMMA"
	LPAREN = "LPAREN"
	RPAREN = "RPAREN"

	IDENTIFIER = "IDENTIFIER"
)

var keywords = make(map[string]string)
var priorities = make(map[string][]any)

func RegisterPriorities() {
	priorities["NONE"] = []any{0, "left"}
	priorities["or"] = []any{1, "left"}
	priorities["and"] = []any{2, "left"}

	priorities["=="] = []any{3, "left"}
	priorities["!="] = []any{3, "left"}

	priorities["<="] = []any{4, "left"}
	priorities["<"] = []any{4, "left"}
	priorities[">"] = []any{4, "left"}
	priorities[">="] = []any{4, "left"}

	priorities["+"] = []any{5, "left"}
	priorities["-"] = []any{5, "left"}

	priorities["*"] = []any{6, "left"}
	priorities["/"] = []any{6, "left"}
	priorities["%"] = []any{6, "left"}

	priorities["!"] = []any{7, "left"}

}

func GetPriority(token string) []any {
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
	keywords["return"] = "return"
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
