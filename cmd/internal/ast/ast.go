package ast

type Node struct {
	Type  string
	Start int
	End   int
}

type (
	Statement interface{}
)

type Literal struct {
	Value any
	Node
	UnaryOp string
}

type Expression interface{}

type UnaryExpression struct {
	Value    Expression
	Operator string
	Node
}

type BinaryExpression struct {
	Left     Expression
	Right    Expression
	Operator string
	Node
}

type Identifier struct {
	Value string
	Node
	UnaryOp string
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

type IfElseBlock struct {
	Consequent Statement
	Alternate  Statement
	Test       any
	Node
}

type IfBlock struct {
	Node
	Consequent Statement
	Test       any
}

type FunctionExpression struct {
	Body []Statement
	Node
	Name       Identifier
	Parameters []Identifier
}

type FunctionEvaluation struct {
	Node
	Arguments []Statement
	Name      Identifier
}

type ReturnStatement struct {
	Value Statement
	Node
}

type Program struct {
	Body []Statement
	Node
}
