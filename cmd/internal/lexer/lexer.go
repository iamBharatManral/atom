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
		return token.New(token.EOF, "", "", 0, 0)
	}
	l.readChar()
	if l.isAtEnd() {
		return token.New(token.ILLEGAL, "", "", 0, 0)
	}
	switch l.currentChar {
	case ';':
		return token.New(token.SEMICOLON, ";", "", l.currentPos, l.currentPos)
	default:
		{
			if unicode.IsDigit(l.currentChar) {
				tok, err := l.numberToken()
				if err != nil {
					fmt.Println(err.Error())
					return token.New(token.EOF, "", "", l.currentPos-1, l.currentPos-1)
				}
				return tok
			}
		}
	}
	return token.New("", "", token.ILLEGAL, 0, 0)
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

func isWhiteSpace(ch rune) bool {
	return ch == '\n' || ch == '\t' || ch == '\r' || ch == ' '
}

func (l *Lexer) readNumber() ([]rune, error) {
	start := l.currentPos
	l.readChar()
	if l.currentChar == 0 {
		return nil, fmt.Errorf("error: missing semicolon at line: %d, column: %d", l.line, l.currentPos-1)
	}
	for unicode.IsDigit(l.currentChar) {
		l.readChar()
	}
	if l.isAtEnd() {
		return nil, fmt.Errorf("error: missing semicolon at line: %d, column: %d", l.line, l.currentPos-1)
	}
	if l.currentChar == '.' || l.currentChar == ';' {
		return l.input[start:l.currentPos], nil
	}
	return nil, fmt.Errorf("error: illegal number at line: %d, columns: %d", l.line, l.currentPos-1)
}

func (l *Lexer) Len() uint {
	return uint(len(l.input))
}
