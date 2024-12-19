package expr

import "github.com/vishrudh-raj-rs-14/lox/token"

//implements Expr
type Binary struct{
	Left     Expr
    Operator token.Token
    Right    Expr
}

func (binExp *Binary) Accept(obj ExprVisitor) interface{}{
	return obj.VisitBinaryExpr(binExp);
}