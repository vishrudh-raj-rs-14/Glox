package expr

//implements Expr
type Literal struct{
	Value interface{}
}

func (lit *Literal) Accept(visitor ExprVisitor) interface{}{
	return visitor.VisitLiteralExpr(lit);
}