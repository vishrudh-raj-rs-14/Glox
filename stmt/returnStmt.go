package stmt

import (
	"github.com/vishrudh-raj-rs-14/lox/expr"
	"github.com/vishrudh-raj-rs-14/lox/token"
)


type Return struct{
	Keyword token.Token
	Expression expr.Expr
}

func (returnStmt *Return) Accept(visitor StmtVisitor) interface{}{
	return	visitor.VisitReturnStmt(returnStmt);
}