package expr

import (
	"github.com/vishrudh-raj-rs-14/lox/token"
)

type Call struct{
	Callee Expr
	Arguments []Expr
	CloseParen token.Token
}

func (call *Call) Accept(visitor ExprVisitor) interface{}{
	return visitor.VisitCallxpr(call);
}