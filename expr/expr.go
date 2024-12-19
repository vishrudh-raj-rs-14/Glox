package expr


type Expr interface{
	Accept(obj ExprVisitor) interface{}
}

type ExprVisitor interface{
	VisitBinaryExpr(expr *Binary) interface{}
	VisitGroupingExpr(expr *Grouping) interface{}
	VisitLiteralExpr(expr *Literal) interface{}
	VisitUnaryExpr(expr *Unary) interface{}
	VisitVariablexpr(expr *Variable) interface{}
	VisitAssignxpr(expr *Assign) interface{}
	VisitLogicalxpr(expr *Logical) interface{}
	VisitCallxpr(expr *Call) interface{}
}

