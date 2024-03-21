package parser

import (
	"errors"
	"github.com/iamBharatManral/atom.git/cmd/internal/ast"
	"github.com/iamBharatManral/atom.git/cmd/internal/lexer"
	"github.com/iamBharatManral/atom.git/cmd/internal/token"
)

type Parser struct {
	lexer        *lexer.Lexer
	currentToken token.Token
	peekToken    token.Token
	errors       []error
}

func New(lexer *lexer.Lexer) *Parser {
	return &Parser{
		lexer:        lexer,
		currentToken: lexer.NextToken(),
		peekToken:    lexer.NextToken(),
	}
}

func (p *Parser) Parse() ast.Program {
	token.RegisterKeyWords()
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
	case token.IDENTIFIER:
		return p.parseDeclaration()
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

func (p *Parser) parseDeclaration() ast.Statement {
	switch token.GetKeyword(p.currentToken.Lexeme()) {
	case "let":
		return p.parseLetDeclaration()
	}
	return p.parseIdentifier()
}

func (p *Parser) parseIdentifier() ast.Statement {
	return ast.Identifier{
		Node: ast.Node{
			Start: p.currentToken.Start(),
			End:   p.currentToken.End(),
			Type:  "Identifier",
		},
		Value: p.currentToken.Lexeme(),
	}
}
func (p *Parser) addError(e error) {
	p.errors = append(p.errors, e)
}

func (p *Parser) parseLetDeclaration() ast.Statement {
	// let a = 10
	start := p.currentToken.Start()
	if p.peekToken.TokenType() != token.IDENTIFIER {
		p.addError(errors.New("error: wrong type in left side of assignment"))
		return nil
	}
	left := ast.Identifier{
		Value: p.peekToken.Lexeme(),
		Node: ast.Node{
			Start: p.peekToken.Start(),
			End:   p.peekToken.End(),
			Type:  "Identifier",
		},
	}
	p.nextToken()
	if p.peekToken.TokenType() != token.ASSIGN {
		p.addError(errors.New("error: invalid operator in let statement"))
		return nil
	}
	p.nextToken()
	operator := p.currentToken.Lexeme()
	var right any
	switch p.peekToken.TokenType() {
	case token.IDENTIFIER:
		right = ast.Identifier{
			Value: p.peekToken.Lexeme(),
			Node: ast.Node{
				Start: p.peekToken.Start(),
				End:   p.peekToken.End(),
				Type:  "Identifier",
			},
		}
	case token.STRING, token.INTEGER, token.FLOAT:
		right = ast.Literal{
			Value: p.peekToken.Value(),
			Node: ast.Node{
				Start: p.peekToken.Start(),
				End:   p.peekToken.End(),
				Type:  "Literal",
			},
		}
	default:
		p.addError(errors.New("error: wrong type in right side of assignment"))
		return nil
	}
	end := p.peekToken.End()
	p.nextToken()
	return ast.LetStatement{
		Left:     left,
		Right:    right,
		Operator: operator,
		Node: ast.Node{
			Start: start,
			End:   end,
			Type:  "Declaration",
		},
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
