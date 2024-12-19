package stmt

import "github.com/vishrudh-raj-rs-14/lox/expr"

type ExpressionStmt struct{
	Expression expr.Expr
}

func (expression *ExpressionStmt) Accept(visitor StmtVisitor) interface{}{
	return visitor.VisitExpressionStmt(expression);
}