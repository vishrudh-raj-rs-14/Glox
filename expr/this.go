package expr

import "github.com/vishrudh-raj-rs-14/lox/token"


type This struct{
	Keyword token.Token
}

func (thisExpr *This) Accept(obj ExprVisitor) interface{}{
	return obj.VisitThisxpr(thisExpr);
}