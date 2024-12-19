package expr

import "github.com/vishrudh-raj-rs-14/lox/token"

type Logical struct{
	Left Expr
	Operator token.Token
	Right Expr
}

func (logical *Logical) Accept(visitor ExprVisitor) interface{}{
	return visitor.VisitLogicalxpr(logical);
}