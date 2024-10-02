package models

type Expression interface {
	Accept(visitor ExpressionVisitor) interface{}
}

// ExpressionVisitor is the interface for the visitor pattern
type ExpressionVisitor interface {
	VisitLiteralExpr(expr *Literal) interface{}
	VisitUnaryExpr(expr *Unary) interface{}
	VisitBinaryExpr(expr *Binary) interface{}
	VisitGroupingExpr(expr *Grouping) interface{}
}

// Literal represents a literal value
type Literal struct {
	Value interface{}
}

func (l *Literal) Accept(visitor ExpressionVisitor) interface{} {
	return visitor.VisitLiteralExpr(l)
}

// Unary represents a unary operation
type Unary struct {
	Operator string
	Right    Expression
}

func (u *Unary) Accept(visitor ExpressionVisitor) interface{} {
	return visitor.VisitUnaryExpr(u)
}

// Binary represents a binary operation
type Binary struct {
	Left     Expression
	Operator string
	Right    Expression
}

func (b *Binary) Accept(visitor ExpressionVisitor) interface{} {
	return visitor.VisitBinaryExpr(b)
}

// Grouping represents a grouped expression
type Grouping struct {
	Expression Expression
}

func (g *Grouping) Accept(visitor ExpressionVisitor) interface{} {
	return visitor.VisitGroupingExpr(g)
}

// Token represents a token in the language
type Token struct {
	Type    string
	Lexeme  string
	Literal interface{}
	Line    int
}
