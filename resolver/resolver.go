package resolver

import (
	"fmt"

	errorhandler "github.com/vishrudh-raj-rs-14/lox/errorHandler"
	"github.com/vishrudh-raj-rs-14/lox/expr"
	"github.com/vishrudh-raj-rs-14/lox/operations"
	"github.com/vishrudh-raj-rs-14/lox/stmt"
	"github.com/vishrudh-raj-rs-14/lox/token"
)

type FunctionType int

const (
	NONE FunctionType = iota
	FUNCTION
	INITIALIZER
	METHOD
)

type ClassType int

const (
	NONECLASS ClassType = iota
	CLASS
	SUBCLASS
)

type Resolver struct {
	interpreter     operations.Interpreter
	currentFunction FunctionType
	currentClass    ClassType
	scopes          []map[string]bool
}


func CreateResolver(interpreter operations.Interpreter) *Resolver {
	scopes := make([]map[string]bool, 0)
	scopes = append(scopes, map[string]bool{})
	return &Resolver{
		interpreter:     interpreter,
		currentFunction: NONE,
		currentClass:    NONECLASS,
		scopes:          scopes,
	}
}

func (resolver *Resolver) beginScope() {
	resolver.scopes = append(resolver.scopes, map[string]bool{})
}

func (resolver *Resolver) endScope() {
	resolver.scopes = resolver.scopes[:len(resolver.scopes)-1]

}

func (resolver *Resolver) declare(varName token.Token) {
	if len(resolver.scopes) == 0 {
		return
	}
	_, ok := resolver.peek()[varName.Lexeme]
	if ok {
		errorhandler.ErrorToken(varName, "Variable already declared")
		return
	}
	resolver.scopes[len(resolver.scopes)-1][varName.Lexeme] = false
}

func (resolver *Resolver) define(varName token.Token) {
	if len(resolver.scopes) == 0 {
		return
	}
	resolver.scopes[len(resolver.scopes)-1][varName.Lexeme] = true

}

func (resolver *Resolver) peek() map[string]bool {
	if len(resolver.scopes) == 0 {
		return nil
	}
	return resolver.scopes[len(resolver.scopes)-1]
}

func (resolver *Resolver) isEmpty() bool {
	return len(resolver.scopes) == 0
}


// VisitSuperxpr implements expr.ExprVisitor.
func (resolver *Resolver) VisitSuperxpr(expr *expr.Super) interface{} {
	if(resolver.currentClass==NONECLASS){
		errorhandler.ErrorToken(expr.Keyword, "Can't use 'super' outside of a class.")
		return nil;
	}else if(resolver.currentClass==SUBCLASS){
		errorhandler.ErrorToken(expr.Keyword, "Can't use 'super' in a class with no superclass.")
		return nil;
	}
	resolver.resolveLocal(expr, expr.Keyword);
	return nil;
}

// VisitThisxpr implements expr.ExprVisitor.
func (resolver *Resolver) VisitThisxpr(expr *expr.This) interface{} {
	if resolver.currentClass == NONECLASS {
		errorhandler.ErrorToken(expr.Keyword, "'this' can be used only inside a class")
		return nil
	}
	resolver.resolveLocal(expr, expr.Keyword)
	return nil
}

func (resolver *Resolver) VisitGetxpr(expr *expr.GetExpr) interface{} {
	resolver.resolveExpr(expr.Object)
	return nil
}

// VisitSetxpr implements expr.ExprVisitor.
func (resolver *Resolver) VisitSetxpr(expr *expr.SetExpr) interface{} {
	resolver.resolveExpr(expr.Object)
	resolver.resolveExpr(expr.Value)
	return nil
}

// VisitAssignxpr implements expr.ExprVisitor.
func (resolvar *Resolver) VisitAssignxpr(expr *expr.Assign) interface{} {
	resolvar.resolveExpr(expr.Value)
	resolvar.resolveLocal(expr, expr.Name)
	return nil
}

// VisitBinaryExpr implements expr.ExprVisitor.
func (resolvar *Resolver) VisitBinaryExpr(expr *expr.Binary) interface{} {
	resolvar.resolveExpr(expr.Left)
	resolvar.resolveExpr(expr.Right)
	return nil
}

// VisitCallxpr implements expr.ExprVisitor.
func (resolvar *Resolver) VisitCallxpr(expr *expr.Call) interface{} {
	resolvar.resolveExpr(expr.Callee)

	for _, argument := range expr.Arguments {
		resolvar.resolveExpr(argument)
	}
	return nil

}

// VisitGroupingExpr implements expr.ExprVisitor.
func (resolvar *Resolver) VisitGroupingExpr(expr *expr.Grouping) interface{} {
	resolvar.resolveExpr(expr.Expression)
	return nil
}

// VisitLiteralExpr implements expr.ExprVisitor.
func (resolvar *Resolver) VisitLiteralExpr(expr *expr.Literal) interface{} {
	return nil
}

// VisitLogicalxpr implements expr.ExprVisitor.
func (resolvar *Resolver) VisitLogicalxpr(expr *expr.Logical) interface{} {
	resolvar.resolveExpr(expr.Left)
	resolvar.resolveExpr(expr.Right)
	return nil
}

// VisitUnaryExpr implements expr.ExprVisitor.
func (resolvar *Resolver) VisitUnaryExpr(expr *expr.Unary) interface{} {
	resolvar.resolveExpr(expr.Right)
	return nil
}

// VisitVariablexpr implements expr.ExprVisitor.
func (resolvar *Resolver) VisitVariablexpr(expr *expr.Variable) interface{} {
	val, ok := resolvar.peek()[expr.Name]
	if ok && val == false {
		errorhandler.ErrorToken(expr.Token, "Can't read local variable in its own initializer")
	}
	resolvar.resolveLocal(expr, expr.Token)
	return nil
}

func (resolvar *Resolver) resolveLocal(expr expr.Expr, name token.Token) {
	done := false
	for i := len(resolvar.scopes) - 1; i >= 0; i-- {
		_, ok := resolvar.scopes[i][name.Lexeme]
		if ok {
			resolvar.interpreter.Resolve(expr, len(resolvar.scopes)-1-i)
			done = true
			break
		}
	}
	if !done {
		fmt.Println(name.TokenType)
		errorhandler.ErrorToken(name, "Variable used without declaration")
	}

}

//-----------------------------------------------

func (resolvar *Resolver) resolveFuntion(fun *stmt.FunStmt, functionType FunctionType) {
	enclosingFun := resolvar.currentFunction
	resolvar.currentFunction = functionType
	resolvar.beginScope()
	for _, parameter := range fun.Parameters {
		resolvar.declare(parameter)
		resolvar.define(parameter)
	}
	resolvar.ResolveStatements(fun.Body.Statements)
	resolvar.endScope()
	resolvar.currentFunction = enclosingFun
}

// VisitClassStmt implements stmt.StmtVisitor.
func (resolver *Resolver) VisitClassStmt(stmt *stmt.ClassStmt) interface{} {
	cur := resolver.currentClass
	resolver.currentClass = CLASS
	resolver.declare(stmt.Name)
	resolver.define(stmt.Name)

	if stmt.SuperClass != nil {
		if stmt.Name.Lexeme == stmt.SuperClass.Name {
			errorhandler.ErrorToken(stmt.SuperClass.Token, "A class cant be inherited to itself")
		} else {
			resolver.resolveExpr(stmt.SuperClass)
		}
	}
	if stmt.SuperClass != nil {
		resolver.currentClass=SUBCLASS
		resolver.beginScope()
		resolver.peek()["super"] = true
	}
	resolver.beginScope()
	resolver.peek()["this"] = true
	for _, method := range stmt.Methods {
		funType := METHOD
		if method.Name.Lexeme == "init" {
			funType = INITIALIZER
		}
		resolver.resolveFuntion(&method, funType)
	}
	if stmt.SuperClass != nil {
		resolver.endScope()
	}
	resolver.endScope()
	resolver.currentClass = cur
	return nil
}

// VisitBlockStmt implements stmt.StmtVisitor.
func (resolvar *Resolver) VisitBlockStmt(stmt *stmt.BlockStmt) interface{} {
	resolvar.beginScope()
	resolvar.ResolveStatements(stmt.Statements)
	resolvar.endScope()
	return nil
}

// VisitExpressionStmt implements stmt.StmtVisitor.
func (resolvar *Resolver) VisitExpressionStmt(stmt *stmt.ExpressionStmt) interface{} {
	resolvar.resolveExpr(stmt.Expression)
	return nil
}

// VisitFunStmt implements stmt.StmtVisitor.
func (resolvar *Resolver) VisitFunStmt(stmt *stmt.FunStmt) interface{} {
	resolvar.declare(stmt.Name)
	resolvar.define(stmt.Name)
	resolvar.resolveFuntion(stmt, FUNCTION)
	return nil
}

// VisitIfStmt implements stmt.StmtVisitor.
func (resolvar *Resolver) VisitIfStmt(stmt *stmt.IfStmt) interface{} {
	resolvar.resolveExpr(stmt.Expression)
	resolvar.resolveStmt(stmt.ThenBlock)
	if stmt.ElseBlock != nil {
		resolvar.resolveStmt(stmt.ElseBlock)
	}
	return nil
}

// VisitPrintStmt implements stmt.StmtVisitor.
func (resolvar *Resolver) VisitPrintStmt(stmt *stmt.PrintStmt) interface{} {
	resolvar.resolveExpr(stmt.Expression)
	return nil
}

// VisitReturnStmt implements stmt.StmtVisitor.
func (resolvar *Resolver) VisitReturnStmt(stmt *stmt.Return) interface{} {
	if resolvar.currentFunction == NONE {
		errorhandler.ErrorToken(stmt.Keyword, "Cant return from top level")
		return nil
	}
	if stmt.Expression != nil {
		if resolvar.currentFunction == INITIALIZER {
			errorhandler.ErrorToken(stmt.Keyword, "Cant return from constructor")
			return nil
		}
	}
	resolvar.resolveExpr(stmt.Expression)
	return nil
}

// VisitVarStmt implements stmt.StmtVisitor.
func (resolvar *Resolver) VisitVarStmt(stmt *stmt.VarStatement) interface{} {
	resolvar.declare(stmt.Name)
	if stmt.Value != nil {
		resolvar.resolveExpr(stmt.Value)
	}
	resolvar.define(stmt.Name)
	return nil
}

// VisitWhileStmt implements stmt.StmtVisitor.
func (resolvar *Resolver) VisitWhileStmt(stmt *stmt.WhileStmt) interface{} {
	resolvar.resolveExpr(stmt.Expression)
	resolvar.resolveStmt(stmt.Body)
	return nil
}

func (resolvar *Resolver) ResolveStatements(statements []stmt.Stmt) {
	for _, statement := range statements {
		resolvar.resolveStmt(statement)
	}
}

func (resolvar *Resolver) resolveStmt(statement stmt.Stmt) {
	if statement == nil {
		return
	}
	statement.Accept(resolvar)
}

func (resolvar *Resolver) resolveExpr(expression expr.Expr) {
	if expression == nil {
		return
	}
	expression.Accept(resolvar)
}
