package expr

import "github.com/vishrudh-raj-rs-14/lox/token"

type Assign struct{
	Name token.Token
	Value Expr
}

func (assign *Assign) Accept(visitor ExprVisitor) interface{}{
	return visitor.VisitAssignxpr(assign);
}