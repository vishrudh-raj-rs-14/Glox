package operations

import (
	"fmt"

	"github.com/vishrudh-raj-rs-14/lox/expr"
)

type PrettyPrint struct {
}

// VisitThisxpr implements expr.ExprVisitor.
func (p PrettyPrint) VisitThisxpr(expr *expr.This) interface{} {
	return "this"
}

// VisitGetxpr implements expr.ExprVisitor.
func (p PrettyPrint) VisitGetxpr(expr *expr.GetExpr) interface{} {
	return p.parenthezise(expr.Name.Lexeme+" of ", expr.Object)
}

func (p PrettyPrint) VisitSetxpr(expr *expr.SetExpr) interface{} {
	return p.parenthezise(expr.Name.Lexeme+" of ", expr.Object, expr.Value)
}

func (p PrettyPrint) Print(expr expr.Expr) string {
	return fmt.Sprintf("%v", expr.Accept(p))
}

func (p PrettyPrint) VisitCallxpr(expr *expr.Call) interface{} {
	return p.parenthezise("functionCall", append(expr.Arguments, expr.Callee)...)
}

func (p PrettyPrint) VisitLogicalxpr(expr *expr.Logical) interface{} {
	return p.parenthezise(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (p PrettyPrint) VisitAssignxpr(expr *expr.Assign) interface{} {
	return fmt.Sprintf("(=  %s, %v)", expr.Name, expr.Value)
}

func (p PrettyPrint) VisitBinaryExpr(expr *expr.Binary) interface{} {
	return p.parenthezise(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (p PrettyPrint) VisitGroupingExpr(expr *expr.Grouping) interface{} {
	return p.parenthezise("group", expr.Expression)
}

func (p PrettyPrint) VisitLiteralExpr(expr *expr.Literal) interface{} {
	return fmt.Sprintf("%v", expr.Value)
}

func (p PrettyPrint) VisitUnaryExpr(expr *expr.Unary) interface{} {
	return p.parenthezise(expr.Operator.Lexeme, expr.Right)
}

func (p PrettyPrint) VisitVariablexpr(expr *expr.Variable) interface{} {
	return ("Variable " + expr.Name)
}

func (p PrettyPrint) parenthezise(name string, exprs ...expr.Expr) string {
	output := "(" + name
	for _, expr := range exprs {
		output += " "
		output += (fmt.Sprintf("%v", expr.Accept(p)))
	}

	output += ")"
	return output

}
