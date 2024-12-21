package expr

import "github.com/vishrudh-raj-rs-14/lox/token"


type Super struct{
	Keyword token.Token
	Method token.Token
}

func (superExpr *Super) Accept(obj ExprVisitor) interface{}{
	return obj.VisitSuperxpr(superExpr);
}