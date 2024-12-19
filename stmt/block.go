package stmt

type BlockStmt struct{
	Statements []Stmt
}

func (expression *BlockStmt) Accept(visitor StmtVisitor) interface{}{
	return visitor.VisitBlockStmt(expression);
}