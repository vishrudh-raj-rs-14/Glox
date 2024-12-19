package stmt

import "github.com/vishrudh-raj-rs-14/lox/expr"

type WhileStmt struct{
	Expression expr.Expr
	Body Stmt
}

func (whileStmt *WhileStmt) Accept(visitor StmtVisitor) interface{}{
	return visitor.VisitWhileStmt(whileStmt);
}