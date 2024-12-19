package stmt

import "github.com/vishrudh-raj-rs-14/lox/expr"

type IfStmt struct{
	Expression expr.Expr
	ThenBlock Stmt
	ElseBlock Stmt
}

func (ifStmt *IfStmt) Accept(visitor StmtVisitor) interface{}{
	return	visitor.VisitIfStmt(ifStmt);
}