package ast

type Node struct {
	Type  string
	Start uint
	End   uint
}

type Statement interface{}

type Literal interface{}

type IntegerLiteral struct {
	Node
	Value int
}

type FloatLiteral struct {
	Node
	Value float64
}

type StringLiteral struct {
	Value string
	Node
}

type Program struct {
	Body []Statement
	Node
}
