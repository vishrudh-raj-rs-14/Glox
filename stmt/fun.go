package stmt

import "github.com/vishrudh-raj-rs-14/lox/token"

type FunStmt struct{
	Name token.Token
	Parameters []token.Token
	Body BlockStmt
}

func (expression *FunStmt) Accept(visitor StmtVisitor) interface{}{
	return visitor.VisitFunStmt(expression);
}