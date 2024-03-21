package ast

type Node struct {
	Type  string
	Start int
	End   int
}

type (
	Statement interface{}
	AstNode   interface{}
)

type Literal struct {
	Value any
	Node
}

type Expression interface{}

type BinaryExpression struct {
	Left     Expression
	Right    Expression
	Operator string
	Node
}

type Identifier struct {
	Value string
	Node
}

type LetStatement struct {
	Left     Identifier
	Right    any
	Operator string
	Node
}

type AssignmentStatement struct {
	Left     Identifier
	Right    any
	Operator string
	Node
}

type Program struct {
	Body []Statement
	Node
}
