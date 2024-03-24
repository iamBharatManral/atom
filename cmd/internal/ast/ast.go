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

type ErrorStatement struct {
	Reason string
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

type Consequent struct {
	Body []Statement
	Node
}

type Alternate struct {
	Body []Statement
	Node
}
type IfElseEBlock struct {
	Consequent
	Alternate
	Test any
	Node
}

type IfBlock struct {
	Node
	Consequent
	Test any
}

type Program struct {
	Body []Statement
	Node
}
