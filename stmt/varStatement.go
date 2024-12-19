package stmt

import (
	"github.com/vishrudh-raj-rs-14/lox/expr"
	"github.com/vishrudh-raj-rs-14/lox/token"
)

type VarStatement struct{
	Name token.Token
	Value expr.Expr
}

func (varStatement *VarStatement) Accept(visitor StmtVisitor) interface{}{
	return visitor.VisitVarStmt(varStatement);
}