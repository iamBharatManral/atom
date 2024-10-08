package lexer

import (
	"fmt"
	"strconv"
	"unicode"

	"github.com/iamBharatManral/atom.git/cmd/internal/token"
)

type Lexer struct {
	input       []rune
	currentPos  int
	currentChar rune
	line        uint
}

func New(input []rune) *Lexer {
	return &Lexer{
		input:       input,
		currentPos:  -1,
		currentChar: 0,
		line:        1,
	}
}
func (l *Lexer) Line() uint {
	return l.line
}
func (l *Lexer) NextToken() token.Token {
	if l.isAtEnd() {
		return l.endOfFileToken()
	}
	l.readChar()
	l.ignoreWhiteSpace()
	if l.isAtEnd() {
		return l.endOfFileToken()
	}
	switch l.currentChar {
	case '+':
		return token.New(token.PLUS, "+", "", l.currentPos, l.currentPos)
	case '-':
		if l.peek() == '>' {
			l.readChar()
			return token.New(token.ARROW, "->", "", l.currentPos-1, l.currentPos)
		}
		return token.New(token.MINUS, "-", "", l.currentPos, l.currentPos)
	case '*':
		return token.New(token.STAR, "*", "", l.currentPos, l.currentPos)
	case '/':
		return token.New(token.SLASH, "/", "", l.currentPos, l.currentPos)
	case '%':
		return token.New(token.MOD, "%", "", l.currentPos, l.currentPos)

	case '=':
		if l.peek() == '=' {
			l.readChar()
			return token.New(token.EQ, "==", "", l.currentPos-1, l.currentPos)
		}
		return token.New(token.ASSIGN, "=", "", l.currentPos, l.currentPos)
	case '"':
		{
			tok, err := l.stringToken()
			if err != nil {
				fmt.Println(err.Error())
				return l.endOfFileToken()
			}
			return tok
		}
	case '<':
		if l.peek() == '=' {
			l.readChar()
			return token.New(token.LE, "<=", "", l.currentPos-1, l.currentPos)
		}
		return token.New(token.LT, "<", "", l.currentPos, l.currentPos)
	case '>':
		if l.peek() == '=' {
			l.readChar()
			return token.New(token.GE, ">=", "", l.currentPos-1, l.currentPos)
		}
		return token.New(token.GT, ">", "", l.currentPos, l.currentPos)
	case '|':
		return token.New(token.BAR, "|", "", l.currentPos, l.currentPos)
	case '!':
		if l.peek() == '=' {
			l.readChar()
			return token.New(token.NE, "!=", "", l.currentPos-1, l.currentPos)
		}
		return token.New(token.NOT, "!", "", l.currentPos, l.currentPos)
	case ',':
		return token.New(token.COMMA, ",", "", l.currentPos, l.currentPos)
	case '(':
		return token.New(token.LPAREN, "(", "", l.currentPos, l.currentPos)
	case ')':
		return token.New(token.RPAREN, ")", "", l.currentPos, l.currentPos)

	case '\n':
		return token.New(token.NEWLINE, "", "", l.currentPos, l.currentPos)
	default:
		{
			if unicode.IsDigit(l.currentChar) {
				tok, err := l.numberToken()
				if err != nil {
					fmt.Println(err.Error())
					return l.endOfFileToken()
				}
				return tok
			} else if unicode.IsLetter(l.currentChar) {
				tok, err := l.identifier()
				if err != nil {
					fmt.Println(err.Error())
					return l.endOfFileToken()
				}
				return tok
			}
		}
	}
	return illegalToken(l.currentPos)
}

func (l *Lexer) PeekToken(peek int) token.Token {
	oldPos := l.currentPos
	oldChar := l.currentChar
	oldLine := l.line
	var tt token.Token
	for i := 1; i <= peek; i += 1 {
		tt = l.NextToken()
		if tt.TokenType() == token.EOF {
			break
		}
	}
	l.currentPos = oldPos
	l.currentChar = oldChar
	l.line = oldLine
	return tt
}

func (l *Lexer) identifier() (token.Token, error) {
	// let
	start := l.currentPos
	for unicode.IsLetter(l.peek()) {
		l.readChar()
	}
	if l.currentPos+1 == len(l.input) {
		id := string(l.input[start : l.currentPos+1])
		return token.New(token.IDENTIFIER, id, "", start, l.currentPos), nil
	}
	if unicode.IsDigit(l.peek()) {
		return token.Token{}, fmt.Errorf("error: invalid identifer at line: %d, column %d", l.line, l.currentPos+1)
	}
	id := string(l.input[start : l.currentPos+1])
	return token.New(token.IDENTIFIER, id, "", start, l.currentPos), nil
}

func (l *Lexer) ignoreWhiteSpace() {
	for l.isWhiteSpace() {
		l.readChar()
	}
}

func illegalToken(column int) token.Token {
	return token.New("", "", token.ILLEGAL, column, column)
}

func (l *Lexer) endOfFileToken() token.Token {
	return token.New(token.EOF, "", "", l.currentPos, l.currentPos)
}

func (l *Lexer) stringToken() (token.Token, error) {
	start := l.currentPos
	l.readChar()
	if l.isAtEnd() {
		return token.Token{}, fmt.Errorf("error: unclosed string at line: %d, column: %d", l.line, l.currentPos-1)
	}
	for l.peek() != '"' && l.peek() != 0 {
		l.readChar()
	}
	if l.peek() == '"' {
		l.readChar()
		return token.New(token.STRING, string(l.input[start:l.currentPos+1]), string(l.input[start+1:l.currentPos]), start, l.currentPos), nil
	}
	return token.Token{}, fmt.Errorf("error: unclosed string at line: %d, column: %d", l.line, l.currentPos)
}

func (l *Lexer) isAtEnd() bool {
	return int(l.currentPos) >= len(l.input)
}

func (l *Lexer) readChar() {
	l.currentPos += 1
	if l.isAtEnd() {
		l.currentChar = 0
		return
	}
	l.currentChar = l.input[l.currentPos]
}

func (l *Lexer) peek() rune {
	l.currentPos += 1
	if l.isAtEnd() {
		l.currentPos -= 1
		return 0
	}
	l.currentPos -= 1
	return l.input[l.currentPos+1]
}

func (l *Lexer) numberToken() (token.Token, error) {
	start := l.currentPos
	beforeDecPoint, err := l.readNumber()
	if err != nil {
		return token.Token{}, err
	}
	if l.peek() == '.' {
		l.readChar()
		afterDecPoint, err := l.readNumber()
		if err != nil {
			return token.Token{}, err
		}
		decimal := append(beforeDecPoint, afterDecPoint...)
		number, _ := strconv.ParseFloat(string(decimal), 64)
		return token.New(token.FLOAT, "", number, start, l.currentPos), nil
	}
	number, _ := strconv.Atoi(string(beforeDecPoint))
	return token.New(token.INTEGER, "", number, start, l.currentPos), nil
}

func (l *Lexer) isWhiteSpace() bool {
	ch := l.currentChar
	return ch == '\t' || ch == '\r' || ch == ' '
}

func (l *Lexer) readNumber() ([]rune, error) {
	start := l.currentPos
	for unicode.IsDigit(l.peek()) {
		l.readChar()
	}
	nextChar := l.peek()
	if nextChar == '.' || unicode.IsSpace(nextChar) || !unicode.IsLetter(nextChar) {
		return l.input[start : l.currentPos+1], nil
	}
	return nil, fmt.Errorf("error: illegal number at line: %d, columns: %d", l.line, l.currentPos+1)
}

func (l *Lexer) Len() uint {
	return uint(len(l.input))
}
