package parser

import (
	"github.com/iamBharatManral/atom.git/cmd/internal/ast"
	"github.com/iamBharatManral/atom.git/cmd/internal/lexer"
	"github.com/iamBharatManral/atom.git/cmd/internal/token"
)

type Parser struct {
	lexer        *lexer.Lexer
	currentToken token.Token
	Ast          ast.Program
}

func New(lexer *lexer.Lexer) *Parser {
	return &Parser{
		lexer: lexer,
	}
}

func (p *Parser) Parse() {
	p.readToken()
	p.Program()
}

func (p *Parser) Program() {
	p.Ast = ast.Program{
		Node: ast.Node{
			Type:  "Program",
			Start: 0,
			End:   p.lexer.Len() - 1,
		},
		Body: []ast.Statement{},
	}
	for p.currentToken.TokenType() != token.EOF {
		switch p.currentToken.TokenType() {
		case token.INTEGER:
			{
				p.Ast.Body = append(p.Ast.Body, ast.IntegerLiteral{
					Node: ast.Node{
						Start: p.currentToken.Start(),
						End:   p.currentToken.End(),
						Type:  "IntegerLiteral",
					},
					Value: p.currentToken.Value().(int),
				})
			}
		case token.FLOAT:
			p.Ast.Body = append(p.Ast.Body, ast.FloatLiteral{
				Node: ast.Node{
					Start: p.currentToken.Start(),
					End:   p.currentToken.End(),
					Type:  "FloatLiteral",
				},
				Value: p.currentToken.Value().(float64),
			})
		case token.STRING:
			p.Ast.Body = append(p.Ast.Body, ast.StringLiteral{
				Node: ast.Node{
					Start: p.currentToken.Start(),
					End:   p.currentToken.End(),
					Type:  "StringLiteral",
				},
				Value: p.currentToken.Value().(string),
			})
		}
		p.readToken()
	}
}

func (p *Parser) readToken() {
	p.currentToken = p.lexer.NextToken()
}
