package stmt

import "github.com/vishrudh-raj-rs-14/lox/expr"


type PrintStmt struct{
	Expression expr.Expr
}

func (printStmt *PrintStmt) Accept(visitor StmtVisitor) interface{}{
	return	visitor.VisitPrintStmt(printStmt);
}