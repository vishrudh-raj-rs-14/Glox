package expr

//implements Expr
type Grouping struct{
	Expression     Expr
}

func (grp *Grouping) Accept(obj ExprVisitor) interface{}{
	return obj.VisitGroupingExpr(grp);
}