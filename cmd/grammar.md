Program :=
StatementList

StatementList :=
Statement
StatementList

Statement :=
Expression
Statment

Expression :=
BinaryExpression
Literal

Literal :=
IntegerLiteral
FloatLiteral
StringLiteral

BinaryExpression :=
Expression Operator Expression

Operator :=
"+"
"-"
"\*"
"/"

IntegerLiteral := Integer
FloatLiteral := Float
StringLiteral := String
