package parser

import (
	"github.com/iamBharatManral/atom.git/cmd/internal/ast"
	"github.com/iamBharatManral/atom.git/cmd/internal/lexer"
	"github.com/iamBharatManral/atom.git/cmd/internal/token"
)

type Parser struct {
	lexer        *lexer.Lexer
	currentToken token.Token
	peekToken    token.Token
}

func New(lexer *lexer.Lexer) *Parser {
	return &Parser{
		lexer:        lexer,
		currentToken: lexer.NextToken(),
		peekToken:    lexer.NextToken(),
	}
}

func (p *Parser) Parse() ast.Program {
	program := ast.Program{
		Node: ast.Node{
			Type:  "Program",
			Start: 0,
			End:   int(p.lexer.Len()) - 1,
		},
		Body: []ast.Statement{},
	}
	for p.currentToken.TokenType() != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Body = append(program.Body, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currentToken.TokenType() {
	case token.INTEGER, token.FLOAT, token.STRING:
		return p.parseLiteral()
	}
	return nil
}

func (p *Parser) parseLiteral() ast.Statement {
	switch p.peekToken.TokenType() {
	case token.PLUS, token.MINUS, token.STAR, token.SLASH:
		return p.parseBinaryExpression()
	default:
		return ast.Literal{
			Node: ast.Node{
				Start: p.currentToken.Start(),
				End:   p.currentToken.End(),
				Type:  "Literal",
			},
			Value: p.currentToken.Value(),
		}
	}
}

func (p *Parser) parseBinaryExpression() ast.Statement {
	start := p.currentToken.Start()
	left := ast.Literal{
		Node: ast.Node{
			Start: p.currentToken.Start(),
			End:   p.currentToken.End(),
			Type:  "Literal",
		},
		Value: p.currentToken.Value(),
	}
	operator := p.peekToken.Lexeme()
	p.nextToken()
	end := p.peekToken.End()
	right := ast.Literal{
		Node: ast.Node{
			Start: p.peekToken.Start(),
			End:   end,
			Type:  "Literal",
		},
		Value: p.peekToken.Value(),
	}
	p.nextToken()
	return ast.BinaryExpression{
		Node: ast.Node{
			Start: start,
			End:   end,
			Type:  "BinaryExpression",
		},
		Left:     left,
		Right:    right,
		Operator: operator,
	}
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}
