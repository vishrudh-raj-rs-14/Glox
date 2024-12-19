package expr

import "github.com/vishrudh-raj-rs-14/lox/token"


type Variable struct{
	Name string;
	Token token.Token;
}

func (variable *Variable) Accept(obj ExprVisitor) interface{}{
	return obj.VisitVariablexpr(variable);
}