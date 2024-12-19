package expr

import "github.com/vishrudh-raj-rs-14/lox/token"

//implements Expr
type Unary struct{
	Right Expr
	Operator token.Token
}

func (unary *Unary) Accept(visitor ExprVisitor) interface{}{
	return visitor.VisitUnaryExpr(unary);
}