package stmt

import (
	"github.com/vishrudh-raj-rs-14/lox/expr"
	"github.com/vishrudh-raj-rs-14/lox/token"
)

type ClassStmt struct{
	Name token.Token
	SuperClass *expr.Variable
	Methods []FunStmt
}

func (expression *ClassStmt) Accept(visitor StmtVisitor) interface{}{
	return visitor.VisitClassStmt(expression);
}