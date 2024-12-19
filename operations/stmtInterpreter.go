package operations

import (
	"fmt"

	"github.com/vishrudh-raj-rs-14/lox/callable"
	"github.com/vishrudh-raj-rs-14/lox/environment"
	errorhandler "github.com/vishrudh-raj-rs-14/lox/errorHandler"
	"github.com/vishrudh-raj-rs-14/lox/expr"
	"github.com/vishrudh-raj-rs-14/lox/stmt"
)

type Interpreter struct {
	Env environment.Environment
	vars map[expr.Expr]int
}



func CreateInterpreter() Interpreter {
	interpreter := Interpreter{
		Env: environment.Environment{
			Values: make(map[string]interface{}),
			
		},
		vars: map[expr.Expr]int{},
	}
	interpreter.Env.Define("clock", &callable.ClockFunction{})

	return interpreter
}

func (interpret Interpreter) Interpret(statements []stmt.Stmt) {
	for _, statement := range statements {
		interpret.Execute(statement)
		if errorhandler.HadError {
			break
		}
	}
}

func (interpreter Interpreter) Resolve(expr expr.Expr, level int) {
	interpreter.vars[expr] = level;
}

func (interpreter Interpreter) Execute(statement stmt.Stmt) {
	if(statement==nil){
		return;
	}
	statement.Accept(interpreter)
}

func (interpret Interpreter) VisitReturnStmt(stmt *stmt.Return) interface{} {
	val := interpret.Evaluate(stmt.Expression);
	panic(Return{Value: val});
}

func (interpret Interpreter) VisitFunStmt(stmt *stmt.FunStmt) interface{} {
	function := Function{
		FunctionDeclaration: *stmt,
		Closure: interpret.Env,
		parameterCount:      len(stmt.Parameters),
	}
	interpret.Env.Define(stmt.Name.Lexeme, &function)
	return nil
}

func (interpreter Interpreter) VisitExpressionStmt(statement *stmt.ExpressionStmt) interface{} {
	interpreter.Evaluate(statement.Expression)
	return nil
}

func (interpreter Interpreter) VisitPrintStmt(statement *stmt.PrintStmt) interface{} {
	val := interpreter.Evaluate(statement.Expression)
	fmt.Println(stringify(val))
	return nil
}

func (interpreter Interpreter) VisitVarStmt(statement *stmt.VarStatement) interface{} {
	var val interface{} = nil
	if statement.Value != nil {
		val = interpreter.Evaluate(statement.Value)
	}
	interpreter.Env.Define(statement.Name.Lexeme, val)
	return nil
}

func (interpret Interpreter) VisitIfStmt(stmt *stmt.IfStmt) interface{} {
	if isTruthy(interpret.Evaluate(stmt.Expression)) {
		interpret.Execute(stmt.ThenBlock)
	} else {
		interpret.Execute(stmt.ElseBlock)
	}
	return nil
}

func (interpret Interpreter) VisitWhileStmt(stmt *stmt.WhileStmt) interface{} {
	for isTruthy(interpret.Evaluate(stmt.Expression)) {
		interpret.Execute(stmt.Body)
	}
	return nil
}

func (interpret Interpreter) VisitBlockStmt(stmt *stmt.BlockStmt) interface{} {
	newEnv := environment.Environment{
		Values:    map[string]interface{}{},
		Enclosing: &interpret.Env,
	}
	interpret.ExecuteBlock(stmt, &newEnv)
	return nil
}

func (interpreter Interpreter) ExecuteBlock(stmt *stmt.BlockStmt, env *environment.Environment) {
	prev := interpreter.Env
	interpreter.Env = *env
	interpreter.Interpret(stmt.Statements)
	interpreter.Env = prev
}
