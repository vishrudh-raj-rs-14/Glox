package stmt

type Stmt interface{
	Accept(visitor StmtVisitor) interface{}
}


type StmtVisitor interface{
	VisitExpressionStmt(stmt *ExpressionStmt) interface{}
	VisitPrintStmt(stmt *PrintStmt) interface{}
	VisitVarStmt(stmt *VarStatement) interface{}
	VisitBlockStmt(stmt *BlockStmt) interface{}
	VisitIfStmt(stmt *IfStmt) interface{}
	VisitWhileStmt(stmt *WhileStmt) interface{}
	VisitFunStmt(stmt *FunStmt) interface{}
	VisitReturnStmt(stmt *Return) interface{}
}

