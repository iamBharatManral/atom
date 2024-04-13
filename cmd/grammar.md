Program := Statement*

Statement := 
    Expression 
    | LetDeclaration 
    | FunctionDeclaration
    | Assignment

Expression := 
    Literal
    | Identifier
    | UnaryExpression
    | BinaryExpression 
    | IfElseExpression 
    | ReturnExpression
    | FunctionEvaluation 
    | '(' Expression ')'


UnaryExpression := 
    < MINUS | BANG > Literal
    | < MINUS | BANG > Identifier

BinaryExpression :=
    Expression <BinaryOp> Expression

BinaryOp := 
    PLUS | MINUS | STAR | SLASH | AND | OR | LT | GT | LE | GE | EQ | NE


