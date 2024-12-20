package expr

import "github.com/vishrudh-raj-rs-14/lox/token"


type SetExpr struct{
	Name token.Token
	Object Expr
	Value Expr
}

func (setExpr *SetExpr) Accept(obj ExprVisitor) interface{}{
	return obj.VisitSetxpr(setExpr);
}