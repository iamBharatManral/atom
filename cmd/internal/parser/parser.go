package parser

import (
	"fmt"

	"github.com/iamBharatManral/atom.git/cmd/internal/ast"
	"github.com/iamBharatManral/atom.git/cmd/internal/lexer"
	"github.com/iamBharatManral/atom.git/cmd/internal/token"
)

type Parser struct {
	lexer        *lexer.Lexer
	currentToken token.Token
	peekToken    token.Token
	Errors       []string
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
	token.RegisterPriorities()
	program := ast.Program{
		Node: ast.Node{
			Type:  "Program",
			Start: 0,
			End:   int(p.lexer.Len()) - 1,
		},
		Body: []ast.Statement{},
	}
	for p.peekToken.TokenType() != token.EOF {
		stmt := p.parseSingleStatement()
		if stmt != nil {
			program.Body = append(program.Body, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseSingleStatement() ast.Statement {
	switch p.peekToken.TokenType() {
	case token.PLUS, token.MINUS, token.STAR, token.SLASH, token.SEMICOLON:
		return p.parseExpression()
	default:
		return p.parseStatement()
	}
}

func (p *Parser) parseExpression() ast.Statement {
	switch p.peekToken.TokenType() {
	case token.PLUS, token.MINUS, token.STAR, token.SLASH:
		return p.parseBinaryExpression()
	case token.SEMICOLON:
		return p.parseUnaryExpression()
	default:
		p.addError(fmt.Sprintf("error: wrong token in right side of the expression at line: %d, column: %d", p.lexer.Line(), p.peekToken.Start()+1))
		return nil
	}
}

func (p *Parser) parseUnaryExpression() ast.Statement {
	switch p.currentToken.TokenType() {
	case token.IDENTIFIER:
		return p.parseIdentifier()
	default:
		l := ast.Literal{
			Node: ast.Node{
				Start: p.currentToken.Start(),
				End:   p.currentToken.End(),
				Type:  "Literal",
			},
			Value: p.currentToken.Value(),
		}
		p.nextToken()
		return l
	}
}

func (p *Parser) parseStatement() ast.Statement {
	if token.GetKeyword(p.currentToken.Lexeme()) == "let" {
		return p.parseLetDeclaration()
	}
	if p.peekToken.TokenType() == token.ASSIGN {
		return p.parseAssignment()
	}
	return p.parseIdentifier()
}

func (p *Parser) parseIdentifier() ast.Statement {
	i := ast.Identifier{
		Node: ast.Node{
			Start: p.currentToken.Start(),
			End:   p.currentToken.End(),
			Type:  "Identifier",
		},
		Value: p.currentToken.Lexeme(),
	}
	p.nextToken()
	return i
}

func (p *Parser) parseAssignment() ast.Statement {
	// a = 10
	start := p.currentToken.Start()
	if p.currentToken.TokenType() != token.IDENTIFIER {
		p.addError("error: wrong type in left side of assignment")
		return nil
	}
	left := ast.Identifier{
		Value: p.currentToken.Lexeme(),
		Node: ast.Node{
			Start: p.currentToken.Start(),
			End:   p.currentToken.End(),
			Type:  "Identifier",
		},
	}
	p.nextToken()
	return p.parseRHS("assign", left, start)
}

func (p *Parser) parseRHS(kind string, left ast.Identifier, start int) ast.Statement {
	var tp string
	if kind == "let" {
		tp = "LetStatement"
	} else {
		tp = "Assignment"
	}
	if p.currentToken.TokenType() != token.ASSIGN {
		p.addError("error: invalid operator in let statement, should be = operator")
		return nil
	}
	operator := p.currentToken.Lexeme()
	var rightSide any
	nextTokenOperator := p.lexer.PeekToken(1)
	nextTokenOpeatorType := nextTokenOperator.TokenType()
	if nextTokenOpeatorType == token.PLUS || nextTokenOpeatorType == token.MINUS || nextTokenOpeatorType == token.STAR || nextTokenOpeatorType == token.SLASH {
		p.nextToken()
		rightSide = p.parseBinaryExpression()
		return ast.LetStatement{
			Left:     left,
			Right:    rightSide,
			Operator: operator,
			Node: ast.Node{
				Start: start,
				End:   rightSide.(ast.BinaryExpression).Node.End,
				Type:  tp,
			},
		}
	} else {
		p.nextToken()
		var end int
		switch p.currentToken.TokenType() {
		case token.IDENTIFIER:
			rightSide = p.parseIdentifier()
			end = rightSide.(ast.Identifier).End
		case token.STRING, token.INTEGER, token.FLOAT:
			rightSide = ast.Literal{
				Value: p.currentToken.Value(),
				Node: ast.Node{
					Start: p.currentToken.Start(),
					End:   p.currentToken.End(),
					Type:  "Literal",
				},
			}
			end = p.currentToken.End()
			p.nextToken()
		default:
			p.addError("error: wrong type in right side of assignment")
			return nil
		}
		if kind == "let" {
			return ast.LetStatement{
				Left:     left,
				Right:    rightSide,
				Operator: operator,
				Node: ast.Node{
					Start: start,
					End:   end,
					Type:  tp,
				},
			}

		} else {
			return ast.AssignmentStatement{
				Left:     left,
				Right:    rightSide,
				Operator: operator,
				Node: ast.Node{
					Start: start,
					End:   end,
					Type:  tp,
				},
			}
		}
	}

}

func (p *Parser) parseLetDeclaration() ast.Statement {
	// let a = 10
	start := p.currentToken.Start()
	if p.peekToken.TokenType() != token.IDENTIFIER {
		p.addError("error: wrong type in left side of assignment")
		return nil
	}
	p.nextToken()
	left := ast.Identifier{
		Value: p.currentToken.Lexeme(),
		Node: ast.Node{
			Start: p.currentToken.Start(),
			End:   p.currentToken.End(),
			Type:  "Identifier",
		},
	}
	p.nextToken()
	return p.parseRHS("let", left, start)
}

func (p *Parser) parseBinaryExpression() ast.Statement {
	// 2 + 3
	start := p.currentToken.Start()
	var leftSide any
	if p.currentToken.TokenType() == token.IDENTIFIER {
		leftSide = p.parseIdentifier()
	} else {
		leftSide = ast.Literal{
			Node: ast.Node{
				Start: p.currentToken.Start(),
				End:   p.currentToken.End(),
				Type:  "Literal",
			},
			Value: p.currentToken.Value(),
		}
		p.nextToken()

	}
	firstOperatorToken := p.currentToken
	secondOperatorToken := p.lexer.PeekToken(1)
	if secondOperatorToken.TokenType() != token.SEMICOLON {
		if token.GetPriority(secondOperatorToken.TokenType()) > token.GetPriority(firstOperatorToken.TokenType()) {
			p.nextToken()
			rightSide := p.parseBinaryExpression()
			return ast.BinaryExpression{
				Node: ast.Node{
					Start: start,
					End:   rightSide.(ast.BinaryExpression).Node.End,
					Type:  "BinaryExpression",
				},
				Left:     leftSide,
				Right:    rightSide,
				Operator: firstOperatorToken.Lexeme(),
			}

		} else {
			p.nextToken()
			var rightSide any
			if p.currentToken.TokenType() == token.IDENTIFIER {
				rightSide = p.parseIdentifier()
				return ast.BinaryExpression{
					Node: ast.Node{
						Start: start,
						End:   rightSide.(ast.Identifier).Node.End,
						Type:  "BinaryExpression",
					},
					Right: p.parseExpression(),
					Left: ast.BinaryExpression{
						Node: ast.Node{
							Start: start,
							End:   rightSide.(ast.Identifier).Node.End,
							Type:  "BinaryExpression",
						},
						Left:     leftSide,
						Right:    rightSide,
						Operator: firstOperatorToken.Lexeme(),
					},
					Operator: secondOperatorToken.Lexeme(),
				}
			} else {
				rightSide = ast.Literal{
					Node: ast.Node{
						Start: p.currentToken.Start(),
						End:   p.currentToken.End(),
						Type:  "Literal",
					},
					Value: p.currentToken.Value(),
				}
				p.nextToken()
				p.nextToken()
				return ast.BinaryExpression{
					Node: ast.Node{
						Start: start,
						End:   rightSide.(ast.Literal).End,
						Type:  "BinaryExpression",
					},
					Right: p.parseExpression(),
					Left: ast.BinaryExpression{
						Node: ast.Node{
							Start: start,
							End:   rightSide.(ast.Literal).End,
							Type:  "BinaryExpression",
						},
						Left:     leftSide,
						Right:    rightSide,
						Operator: firstOperatorToken.Lexeme(),
					},
					Operator: secondOperatorToken.Lexeme(),
				}
			}
		}

	} else {
		operator := p.currentToken.Lexeme()
		p.nextToken()
		var rightSide any
		if p.currentToken.TokenType() == token.IDENTIFIER {
			rightSide = p.parseIdentifier()
			return ast.BinaryExpression{
				Node: ast.Node{
					Start: start,
					End:   rightSide.(ast.Identifier).End,
					Type:  "BinaryExpression",
				},
				Left:     leftSide,
				Right:    rightSide,
				Operator: operator,
			}
		} else {
			rightSide := ast.Literal{
				Node: ast.Node{
					Start: p.currentToken.Start(),
					End:   p.currentToken.End(),
					Type:  "Literal",
				},
				Value: p.currentToken.Value(),
			}
			p.nextToken()
			return ast.BinaryExpression{
				Node: ast.Node{
					Start: start,
					End:   rightSide.Node.End,
					Type:  "BinaryExpression",
				},
				Left:     leftSide,
				Right:    rightSide,
				Operator: firstOperatorToken.Lexeme(),
			}

		}
	}
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) addError(error string) {
	p.Errors = append(p.Errors, error)
}
