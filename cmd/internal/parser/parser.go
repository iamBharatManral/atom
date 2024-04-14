package parser

import (
	"fmt"
	"reflect"
	"slices"

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
	for p.currentToken.TokenType() != token.EOF {
		if p.currentToken.TokenType() == token.NEWLINE {
			p.nextToken()
			continue
		}
		stmt := p.parseStatement()
		if stmt != nil {
			program.Body = append(program.Body, stmt)
		}
	}
	return program
}

func (p *Parser) parseStatement() ast.Statement {
	if p.currentToken.Lexeme() == "let" {
		return p.parseLetDeclaration()
	} else if p.currentToken.Lexeme() == "fn" {
		return p.parseFunctionDeclaration()
	} else if p.peekToken.TokenType() == token.ASSIGN {
		return p.parseAssignment()
	}
	return p.parseExpression()
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
	p.nextToken()
	return p.parseRHS("let", left, start)
}

func (p *Parser) parseFunctionDeclaration() ast.Statement {
	// fn hello |a, b| -> a end
	start := p.currentToken.Start()
	p.nextToken()
	var name ast.Identifier
	if p.currentToken.TokenType() != token.IDENTIFIER {
		name = ast.Identifier{}
		p.nextToken()
	} else {
		name = p.parseIdentifier("").(ast.Identifier)
		p.nextToken()
		p.nextToken()
	}

	var parameters []ast.Identifier
	for p.currentToken.TokenType() != token.BAR {
		if p.currentToken.TokenType() == token.IDENTIFIER {
			parameters = append(parameters, p.parseIdentifier("").(ast.Identifier))
			p.nextToken()
		} else {
			p.nextToken()
		}
	}
	p.nextToken()
	p.nextToken()
	for p.currentToken.TokenType() == token.NEWLINE {
		p.nextToken()
	}
	var body []ast.Statement
	for p.currentToken.Lexeme() != "end" {
		for p.currentToken.TokenType() == token.NEWLINE {
			p.nextToken()
		}
		body = append(body, p.parseStatement())
		if p.currentToken.Lexeme() != "end" {
			p.nextToken()
		}
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
			Start: start,
			End:   p.currentToken.End(),
			Type:  "Identifier",
		},
	}
	p.nextToken()
	return p.parseRHS("assign", left, start)
}

func (p *Parser) parseExpression() ast.Statement {
	switch p.currentToken.Lexeme() {
	case "if":
		return p.parseIfExpression()
	case "return":
		return p.parseReturnExpression()
	default:
		postfix := p.convertToPostFixNotation()
		if len(postfix) == 1 {
			return postfix[0]
		}
		return p.createASTFromPostfixExpression(postfix)
	}
}

func (p *Parser) parseFunctionEvaluation() ast.Statement {
	// hello(a, b)
	name := p.parseIdentifier("").(ast.Identifier)
	p.nextToken()
	var args []ast.Statement
	p.nextToken()

	for p.currentToken.TokenType() != token.RPAREN {
		if p.peekToken.TokenType() == token.COMMA {
			args = append(args, p.callAppropriateFunction(p.currentToken.TokenType(), ""))
			p.nextToken()
			p.nextToken()
			continue
		} else if p.peekToken.TokenType() == token.RPAREN {
			args = append(args, p.callAppropriateFunction(p.currentToken.TokenType(), ""))
			p.nextToken()
			continue
		}
		args = append(args, p.parseExpression())
		p.nextToken()
	}

	fnEval := ast.FunctionEvaluation{
		Node: ast.Node{
			Start: name.Start,
			End:   p.currentToken.Start(),
			Type:  "FunctionEvaluation",
		},
		Arguments: args,
		Name:      name,
	}
	return fnEval

}

func (p *Parser) isBinaryOperator(token token.Token) bool {
	switch token.TokenType() {
	case "PLUS", "MINUS", "STAR", "SLASH", "EQ", "NE", "GT", "LT", "GE", "LE", "MOD":
		return true
	default:
		if token.Lexeme() == "and" || token.Lexeme() == "or" {
			return true
		}
		return false
	}
}

func (p *Parser) parseReturnExpression() ast.Statement {
	start := p.currentToken.Start()
	p.nextToken()
	if p.peekToken.TokenType() != token.NEWLINE {
		result := p.parseStatement()
		return ast.ReturnStatement{
			Value: result,
			Node: ast.Node{
				Start: start,
				End:   p.getEndOfStatement(result),
				Type:  "returnstatement",
			},
		}
	}
	switch p.currentToken.TokenType() {
	case token.IDENTIFIER:
		v := p.parseIdentifier("")
		return ast.ReturnStatement{
			Value: v,
			Node: ast.Node{
				Start: start,
				End:   v.(ast.Identifier).End,
				Type:  "returnstatement",
			},
		}
	case token.STRING, token.INTEGER, token.FLOAT:
		stmt := ast.ReturnStatement{
			Value: ast.Literal{
				Node: ast.Node{
					Start: p.currentToken.Start(),
					End:   p.currentToken.End(),
					Type:  "Literal",
				},
				Value: p.currentToken.Value(),
			},
			Node: ast.Node{
				Start: start,
				End:   p.currentToken.End(),
				Type:  "ReturnStatement",
			},
		}
		p.nextToken()
		return stmt

	}
	p.addError("error: unsupported return type, must be literal or identifier")
	return nil
}

func (p *Parser) parseIfExpression() ast.Statement {
	// if 10 < 20 do a else b
	start := p.currentToken.Start()
	var test any
	p.nextToken()
	if p.isBinaryOperator(p.peekToken) {
		test = p.parseExpression()
	} else if p.peekToken.Lexeme() == "do" {
		test = p.parseIdentifier("")
		p.nextToken()
	}
	if keyword := token.GetKeyword(p.currentToken.Lexeme()); keyword != "do" {
		p.addError("error: missing 'do' symbol")
		return nil
	}
	p.nextToken()
	consequent := p.parseExpression()
	end := p.currentToken.Start() - 1
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
	alternate := p.parseExpression()
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

func (p *Parser) parseRHS(kind string, left ast.Identifier, start int) ast.Statement {
	if p.currentToken.Lexeme() == "fn" {
		rightSide := p.parseFunctionDeclaration()
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
	} else if p.currentToken.Lexeme() == "if" {
		rightSide := p.parseIfExpression()
		return ast.LetStatement{
			Left:  left,
			Right: rightSide,
			Node: ast.Node{
				Start: start,
				End:   p.getEndOfStatement(rightSide),
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
		p.nextToken()
	}
	rightSide := p.parseExpression()
	if kind == "let" {
		return ast.LetStatement{
			Left:     left,
			Right:    rightSide,
			Operator: "=",
			Node: ast.Node{
				Start: start,
				End:   p.getEndOfStatement(rightSide),
				Type:  tp,
			},
		}

	} else {
		return ast.AssignmentStatement{
			Left:     left,
			Right:    rightSide,
			Operator: "=",
			Node: ast.Node{
				Start: start,
				End:   p.getEndOfStatement(rightSide),
				Type:  tp,
			},
		}
	}

}
func (p *Parser) createASTFromPostfixExpression(tokens []any) ast.Statement {
	stack := []any{}
	for _, val := range tokens {
		if reflect.TypeOf(val).Kind() == reflect.String {
			right := stack[len(stack)-1]
			left := stack[len(stack)-2]
			bExp := ast.BinaryExpression{
				Left:     left,
				Right:    right,
				Operator: val.(string),
				Node: ast.Node{
					Start: p.getStartOfStatement(left),
					End:   p.getEndOfStatement(right),
					Type:  "BinaryExpression",
				},
			}
			stack = stack[0 : len(stack)-2]
			stack = append(stack, bExp)

		} else {
			stack = append(stack, val)
		}
	}
	if len(stack) == 1 {
		return stack[0]
	}
	p.addError("error: syntax error!")
	return nil
}

func (p *Parser) convertToPostFixNotation() []any {
	queue := []any{}
	stack := []any{}
	var prevToken any
	tokens := []string{"NEWLINE", "EOF", "COMMA"}
	lexemes := []string{"else", "end", "do"}
	for (!slices.Contains(tokens, p.currentToken.TokenType())) && (!slices.Contains(lexemes, p.currentToken.Lexeme())) {
		current := p.currentToken
		if (current.TokenType() == token.MINUS || current.TokenType() == token.NOT) && (prevToken == nil || prevToken == "(") {
			p.nextToken()
			queue = append(queue, p.callAppropriateFunction(p.currentToken.TokenType(), current.Lexeme()))
			p.nextToken()
			continue
		}
		if (current.TokenType() == token.IDENTIFIER || current.TokenType() == token.INTEGER || current.TokenType() == token.FLOAT || current.TokenType() == token.STRING) && (current.Lexeme() != "and" && current.Lexeme() != "or") {
			queue = append(queue, p.callAppropriateFunction(current.TokenType(), ""))
		} else if p.isBinaryOperator(current) {
			for {
				if len(stack) == 0 {
					break
				}
				topElm := stack[len(stack)-1]
				if topElm == "(" {
					break
				}
				p1 := token.GetPriority(string(current.Lexeme()))
				p2 := token.GetPriority(topElm.(string))
				if p2[0].(int) > p1[0].(int) || (p1[0].(int) == p2[0].(int) && p1[1] == "left") {
					queue = append(queue, topElm)
					stack = stack[0 : len(stack)-1]
				} else {
					break
				}
			}

			stack = append(stack, current.Lexeme())
		} else if current.TokenType() == token.LPAREN {
			stack = append(stack, p.currentToken.Lexeme())
		} else {
			for {
				if len(stack) == 0 {
					p.addError(fmt.Sprintf("error: unbalanced parenthesis"))
					return nil
				}
				topElm := stack[len(stack)-1]
				if topElm == "(" {
					break
				} else {
					queue = append(queue, topElm)
					stack = stack[0 : len(stack)-1]
				}
			}
			stack = stack[0 : len(stack)-1]
		}
		prevToken = p.currentToken.Lexeme()
		p.nextToken()
	}
	for len(stack) != 0 {
		topElm := stack[len(stack)-1]
		if topElm == "(" {
			p.addError("error: unbalanced parenthesis")
			return []any{}
		}
		queue = append(queue, topElm)
		stack = stack[0 : len(stack)-1]
	}
	return queue
}

func (p *Parser) callAppropriateFunction(tokenType string, op string) ast.Statement {
	switch tokenType {
	case token.IDENTIFIER:
		if p.peekToken.TokenType() == token.LPAREN {
			return p.parseFunctionEvaluation()
		}
		return p.parseIdentifier(op)
	case token.FLOAT, token.INTEGER, token.STRING:
		return p.parseLiteral(op)
	default:
		return p.parseExpression()
	}
}

func (p *Parser) parseLiteral(op string) ast.Statement {
	return ast.Literal{
		Node: ast.Node{
			Start: p.currentToken.Start(),
			End:   p.currentToken.End(),
			Type:  "Literal",
		},
		Value:   p.currentToken.Value(),
		UnaryOp: op,
	}

}

func (p *Parser) parseIdentifier(op string) ast.Statement {
	return ast.Identifier{
		Node: ast.Node{
			Start: p.currentToken.Start(),
			End:   p.currentToken.End(),
			Type:  "Identifier",
		},
		Value:   p.currentToken.Lexeme(),
		UnaryOp: op,
	}
}

func (p *Parser) getStartOfStatement(stmt ast.Statement) int {
	switch stmt.(type) {
	case ast.Literal:
		return stmt.(ast.Literal).Start
	case ast.Identifier:
		return stmt.(ast.Identifier).Start
	case ast.BinaryExpression:
		return stmt.(ast.BinaryExpression).Start
	case ast.IfBlock:
		return stmt.(ast.IfBlock).Start
	case ast.IfElseBlock:
		return stmt.(ast.IfElseBlock).Start
	case ast.ReturnStatement:
		return stmt.(ast.ReturnStatement).Start
	case ast.LetStatement:
		return stmt.(ast.LetStatement).Start
	case ast.AssignmentStatement:
		return stmt.(ast.AssignmentStatement).Start
	case ast.FunctionEvaluation:
		return stmt.(ast.FunctionEvaluation).Start
	case ast.FunctionExpression:
		return stmt.(ast.FunctionExpression).Start
	}
	return 0
}
func (p *Parser) getEndOfStatement(stmt ast.Statement) int {
	switch stmt.(type) {
	case ast.Literal:
		return stmt.(ast.Literal).End
	case ast.Identifier:
		return stmt.(ast.Identifier).End
	case ast.BinaryExpression:
		return stmt.(ast.BinaryExpression).End
	case ast.IfBlock:
		return stmt.(ast.IfBlock).End
	case ast.IfElseBlock:
		return stmt.(ast.IfElseBlock).End
	case ast.ReturnStatement:
		return stmt.(ast.ReturnStatement).End
	case ast.LetStatement:
		return stmt.(ast.LetStatement).End
	case ast.AssignmentStatement:
		return stmt.(ast.AssignmentStatement).End
	case ast.FunctionEvaluation:
		return stmt.(ast.FunctionEvaluation).End
	case ast.FunctionExpression:
		return stmt.(ast.FunctionExpression).End
	}
	return 0
}
func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) addError(error string) {
	p.Errors = append(p.Errors, error)
}
