package lexer

import (
	"fmt"
	"strconv"
	"unicode"

	"github.com/iamBharatManral/atom.git/cmd/internal/token"
)

type Lexer struct {
	input       []rune
	currentPos  uint
	currentChar rune
	line        uint
}

func New(input []rune) *Lexer {
	return &Lexer{
		input:       input,
		currentPos:  0,
		currentChar: 0,
		line:        1,
	}
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
	case '"':
		{
			tok, err := l.stringToken()
			if err != nil {
				fmt.Println(err.Error())
				return l.endOfFileToken()
			}
			return tok
		}
	default:
		{
			if unicode.IsDigit(l.currentChar) {
				tok, err := l.numberToken()
				if err != nil {
					fmt.Println(err.Error())
					return l.endOfFileToken()
				}
				return tok
			}
		}
	}
	return illegalToken()
}

func (l *Lexer) ignoreWhiteSpace() {
	for l.isWhiteSpace() {
		l.readChar()
	}
}

func illegalToken() token.Token {
	return token.New("", "", token.ILLEGAL, 0, 0)
}

func (l *Lexer) endOfFileToken() token.Token {
	return token.New(token.EOF, "", "", l.currentPos, l.currentPos)
}

func (l *Lexer) stringToken() (token.Token, error) {
	start := l.currentPos
	l.readChar()
	for l.currentChar != '"' {
		l.readChar()
	}
	if l.isAtEnd() {
		return token.Token{}, fmt.Errorf("error: not a valid string")
	}
	if l.currentChar == '"' {
		return token.New(token.STRING, string(l.input[start:l.currentPos+1]), string(l.input[start+1:l.currentPos]), start, l.currentPos), nil
	}
	return token.Token{}, fmt.Errorf("error: not a valid string")
}

func (l *Lexer) isAtEnd() bool {
	return int(l.currentPos) >= len(l.input)
}

func (l *Lexer) readChar() {
	if l.currentChar == 0 {
		l.currentChar = l.input[l.currentPos]
	} else {
		l.currentPos += 1
		if l.isAtEnd() {
			l.currentChar = 0
			return
		}
		l.currentChar = l.input[l.currentPos]
	}
}

func (l *Lexer) numberToken() (token.Token, error) {
	start := l.currentPos
	beforeDecPoint, err := l.readNumber()
	if err != nil {
		return token.Token{}, err
	}
	if l.currentChar == '.' {
		afterDecPoint, err := l.readNumber()
		if err != nil {
			return token.Token{}, err
		}
		decimal := append(beforeDecPoint, afterDecPoint...)
		number, _ := strconv.ParseFloat(string(decimal), 64)
		return token.New(token.FLOAT, "", number, start, l.currentPos-1), nil
	}
	number, _ := strconv.Atoi(string(beforeDecPoint))
	return token.New(token.INTEGER, "", number, start, l.currentPos-1), nil
}

func (l *Lexer) isWhiteSpace() bool {
	ch := l.currentChar
	if ch == '\n' {
		l.line += 1
	}
	return ch == '\n' || ch == '\t' || ch == '\r' || ch == ' '
}

func (l *Lexer) readNumber() ([]rune, error) {
	start := l.currentPos
	l.readChar()
	if l.isAtEnd() {
		return l.input[start:l.currentPos], nil
	}
	for unicode.IsDigit(l.currentChar) {
		l.readChar()
	}
	if l.isAtEnd() {
		return l.input[start:l.currentPos], nil
	}
	if l.currentChar == '.' || l.isWhiteSpace() {
		return l.input[start:l.currentPos], nil
	}
	return nil, fmt.Errorf("error: illegal number at line: %d, columns: %d", l.line, l.currentPos-1)
}

func (l *Lexer) Len() uint {
	return uint(len(l.input))
}
