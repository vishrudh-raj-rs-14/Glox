package expr

import "github.com/vishrudh-raj-rs-14/lox/token"


type GetExpr struct{
	Name token.Token
	Object Expr
}

func (getExpr *GetExpr) Accept(obj ExprVisitor) interface{}{
	return obj.VisitGetxpr(getExpr);
}