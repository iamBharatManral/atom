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
	case token.PLUS, token.MINUS, token.STAR, token.SLASH, token.SEMICOLON, token.EQ, token.NE, token.GT, token.LT, token.GE, token.LE:
		return p.parseExpression()
	case token.LPAREN:
		return p.parseFunctionEvaluation()
	default:
		return p.parseStatement()
	}
}
func (p *Parser) parseFunctionEvaluation() ast.Statement {
	// hello(a, b);
	name := p.parseIdentifier().(ast.Identifier)
	p.nextToken()
	var params []ast.Statement
	for p.currentToken.TokenType() != token.RPAREN {
		if p.currentToken.TokenType() != token.COMMA {
			switch p.currentToken.TokenType() {
			case token.IDENTIFIER:
				params = append(params, p.parseIdentifier())
				continue
			case token.INTEGER, token.FLOAT, token.STRING:
				params = append(params, ast.Literal{
					Node: ast.Node{
						Start: p.currentToken.Start(),
						End:   p.currentToken.End(),
						Type:  "Literal",
					},
					Value: p.currentToken.Value(),
				})
				p.nextToken()
				continue
			}
		}
		p.nextToken()
	}
	p.nextToken()
	return ast.FunctionEvaluation{
		Node: ast.Node{
			Start: name.Start,
			End:   p.currentToken.End() - 1,
			Type:  "FunctionEvaluation",
		},
		Parameters: params,
		Name:       name,
	}

}

func (p *Parser) parseExpression() ast.Statement {
	switch p.peekToken.TokenType() {
	case token.PLUS, token.MINUS, token.STAR, token.SLASH, token.EQ, token.NE, token.GT, token.LT, token.GE, token.LE:
		return p.parseBinaryExpression()
	case token.SEMICOLON, token.NOT:
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
	if token.GetKeyword(p.currentToken.Lexeme()) == "if" {
		return p.parseIfExpression()
	}
	if token.GetKeyword(p.currentToken.Lexeme()) == "fn" {
		return p.parseFunctionExpression()
	}
	if p.peekToken.TokenType() == token.ASSIGN {
		return p.parseAssignment()
	}
	return p.parseIdentifier()
}

func (p *Parser) parseIfExpression() ast.Statement {
	// if 10 < 20 do a; else b;
	start := p.currentToken.Start()
	var test any
	p.nextToken()
	nextTokenType := p.peekToken.TokenType()
	switch nextTokenType {
	case token.LE, token.LT, token.GT, token.GE, token.NE, token.EQ:
		test = p.parseBinaryExpression()
	case token.IDENTIFIER:
		value := p.currentToken.Lexeme()
		if value == "true" || value == "false" {
			test = ast.Identifier{
				Node: ast.Node{
					Start: p.currentToken.Start(),
					End:   p.currentToken.End(),
					Type:  "Identifier",
				},
				Value: p.currentToken.Lexeme(),
			}
			p.nextToken()

		} else {
			p.addError("error: wrong conditional type, allowed types are 'true' or 'false' keyword or any conditional expression")
			return nil
		}

	}
	if keyword := token.GetKeyword(p.currentToken.Lexeme()); keyword != "do" {
		p.addError("error: missing 'do' symbol")
		return nil
	}
	p.nextToken()
	consequent := p.parseSingleStatement()
	end := p.currentToken.Start() - 1
	p.nextToken()
	if keyword := token.GetKeyword(p.currentToken.Lexeme()); keyword != "else" {
		return ast.IfBlock{
			Node: ast.Node{
				Start: start,
				End:   end,
				Type:  "IfExpression",
			},
			Consequent: consequent,
			Test:       test,
		}
	}
	p.nextToken()
	alternate := p.parseSingleStatement()
	return ast.IfElseBlock{
		Consequent: consequent,
		Alternate:  alternate,
		Node: ast.Node{
			Start: start,
			End:   p.currentToken.Start() - 1,
			Type:  "IfElseExpression",
		},
		Test: test,
	}

}

func (p *Parser) parseFunctionExpression() ast.Statement {
	// fn hello |a, b| -> a; end;
	start := p.currentToken.Start()
	p.nextToken()
	var name ast.Identifier
	if p.currentToken.TokenType() != token.IDENTIFIER {
		name = ast.Identifier{}
	} else {
		name = p.parseIdentifier().(ast.Identifier)
	}
	var parameters []ast.Identifier
	p.nextToken()
	for p.currentToken.TokenType() != token.BAR {
		if p.currentToken.TokenType() == token.IDENTIFIER {
			parameters = append(parameters, p.parseIdentifier().(ast.Identifier))
		} else {

			p.nextToken()
		}
	}
	p.nextToken()
	p.nextToken()
	var body []ast.Statement
	for p.currentToken.Lexeme() != "end" {
		body = append(body, p.parseSingleStatement())
		p.nextToken()
	}
	end := p.currentToken.End()
	p.nextToken()
	return ast.FunctionExpression{
		Node: ast.Node{
			Start: start,
			End:   end,
			Type:  "FunctionExpression",
		},
		Body:       body,
		Parameters: parameters,
		Name:       name,
	}
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
	if p.peekToken.Lexeme() == "fn" {
		p.nextToken()
		rightSide := p.parseFunctionExpression()
		return ast.LetStatement{
			Left:  left,
			Right: rightSide,
			Node: ast.Node{
				Start: start,
				End:   rightSide.(ast.FunctionExpression).End,
				Type:  "LetStatement",
			},
			Operator: "=",
		}
	}
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
	nextTokenType := nextTokenOperator.TokenType()
	if nextTokenType == token.PLUS || nextTokenType == token.MINUS || nextTokenType == token.STAR || nextTokenType == token.SLASH || nextTokenType == token.EQ || nextTokenType == token.NE || nextTokenType == token.GE || nextTokenType == token.GT || nextTokenType == token.LE || nextTokenType == token.LT {
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
	if secondOperatorToken.TokenType() != token.SEMICOLON && secondOperatorToken.TokenType() != token.IDENTIFIER {
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
