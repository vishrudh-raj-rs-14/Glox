package operations

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/vishrudh-raj-rs-14/lox/callable"
	errorhandler "github.com/vishrudh-raj-rs-14/lox/errorHandler"
	"github.com/vishrudh-raj-rs-14/lox/expr"
	"github.com/vishrudh-raj-rs-14/lox/token"
)



func (interpret Interpreter) Evaluate(expr expr.Expr) interface{}{
	return expr.Accept(interpret);
}


func (interpret Interpreter) VisitThisxpr(expr *expr.This) interface{} {
	val, err := interpret.variableLookup(expr.Keyword, expr);
	if(err!=nil){
		return nil;
	}
	return val;
}


func (interpret Interpreter) VisitCallxpr(expr *expr.Call) interface{} {
	callee := interpret.Evaluate(expr.Callee);
	var arguments []interface{};
	for _, argument := range expr.Arguments{
		arguments = append(arguments, interpret.Evaluate(argument));
		if(errorhandler.HadError){
			return nil;
		}
	}
	function, ok := callee.(callable.LoxCallable);
    if !ok {
        errorhandler.ErrorToken(expr.CloseParen, "Only functions can be called")
		return nil;
    }

	if(function.Arity()!=len(arguments)){
        errorhandler.ErrorToken(expr.CloseParen, "Arguments to the function dont match")
		return nil;
	}
	
	return function.Call(interpret, arguments);;
}

func(interpret Interpreter) VisitBinaryExpr(expr *expr.Binary) interface{}{
	if(expr==nil){
		return nil;
	}

	left := interpret.Evaluate(expr.Left);
	right := interpret.Evaluate(expr.Right);
	switch expr.Operator.TokenType {
    case token.MINUS:
		if(!checkNumberOperands(expr.Operator, left, right)){
			return nil
		}
        return left.(float64) - right.(float64)
    case token.SLASH:
		if(!checkNumberOperands(expr.Operator, left, right)){
			return nil
		}
		if(right.(float64)==0){
			errorhandler.ErrorToken(expr.Operator, "Cannot divide by zero")
			return nil
		}
        return left.(float64) / right.(float64)
    case token.STAR:
		if(!checkNumberOperands(expr.Operator, left, right)){
			return nil
		}
        return left.(float64) * right.(float64)
	case token.MOD:
		if(!checkNumberOperands(expr.Operator, left, right)){
			return nil
		}
		leftVal := left.(float64);
		rightVal := right.(float64);
		if((leftVal-float64(int(leftVal)))!=0 || (rightVal-float64(int(rightVal)))!=0){
			errorhandler.ErrorToken(expr.Operator, "Mod can be performed only with integers");
		}
        return float64(int(leftVal) % int(rightVal))
	
	case token.PLUS:
		switch left.(type) {
		case float64:
			if rightValue, ok := right.(float64); ok {
				return left.(float64) + rightValue
			}else if rightValue, ok := right.(string); ok {
				return strconv.FormatFloat(left.(float64), 'f', -1, 64) + rightValue
			}
			return nil
		case string:
			if rightValue, ok := right.(string); ok {
				return left.(string) + rightValue
			}else if rightValue, ok := right.(float64); ok {
				return left.(string) + strconv.FormatFloat(rightValue, 'f', -1, 64)
			}
			return nil
		default:
			errorhandler.ErrorToken(expr.Operator, "Both the operand must be strings")
			return nil
		}
	case token.GREATER:
		if(!checkNumberOperands(expr.Operator, left, right)){
			return nil
		}
		return left.(float64) > right.(float64)
	case token.LESS:
		if(!checkNumberOperands(expr.Operator, left, right)){
			return nil
		}
		return left.(float64) < right.(float64)
	case token.GREATER_EQUAL:
		if(!checkNumberOperands(expr.Operator, left, right)){
			return nil
		}
		return left.(float64) >= right.(float64)
	case token.LESS_EQUAL:
		if(!checkNumberOperands(expr.Operator, left, right)){
			return nil
		}
		return left.(float64) <= right.(float64)
	case token.EQUAL_EQUAL:
		return isEqual(left, right);
	case token.BANG_EQUAL:
		return !(isEqual(left, right));

    }

	return nil;
}

func(interpret Interpreter) VisitGroupingExpr(expr *expr.Grouping) interface{}{
	if(expr==nil){
		return nil;
	}
	return interpret.Evaluate(expr.Expression);
}

func(interpret Interpreter) VisitLiteralExpr(expr *expr.Literal) interface{}{
	if(expr==nil){
		return nil;
	}
	return expr.Value;
}

func(interpret Interpreter) VisitUnaryExpr(expr *expr.Unary) interface{}{
	if(expr==nil){
		return nil;
	}
	right := interpret.Evaluate(expr.Right);

	switch (expr.Operator.TokenType){
	case token.MINUS:
		return -(right.(float64));
	case token.BANG:
		return !(isTruthy(right));
	}
	return nil;
}

func(interpret Interpreter) VisitVariablexpr(expr *expr.Variable) interface{}{
	val, err := interpret.variableLookup(expr.Token, expr);
	if(err!=nil){
		return nil;
	}
	return val;
}

func (interpret Interpreter) VisitAssignxpr(expr *expr.Assign) interface{} {
	interpret.Env.AssignAt(interpret.vars[expr] ,expr.Name.Lexeme, interpret.Evaluate(expr.Value));
	return interpret.Evaluate(expr.Value);

}

func (interpret Interpreter) VisitLogicalxpr(expr *expr.Logical) interface{} {
	if(expr.Operator.TokenType == token.OR){
		if(isTruthy(interpret.Evaluate(expr.Left))){
			return interpret.Evaluate(expr.Left);
		}else{
			return interpret.Evaluate(expr.Right)
		}
	}else{
		if(!isTruthy(interpret.Evaluate(expr.Left))){
			return interpret.Evaluate(expr.Left);
		}else{
			return interpret.Evaluate(expr.Right);
		}
	}
}

func (interpret Interpreter) VisitGetxpr(expr *expr.GetExpr) interface{} {
	val := interpret.Evaluate(expr.Object);
	ClassVal, ok := val.(*Instance);
	if(!ok || ClassVal==nil){
		errorhandler.ErrorToken(expr.Name, "Only instances can have properties");
		return nil;
	}
	return ClassVal.Get(expr.Name);
}

func (interpret Interpreter) VisitSetxpr(expr *expr.SetExpr) interface{} {
	val := interpret.Evaluate(expr.Object);
	ClassVal, ok := val.(*Instance);
	if(!ok || ClassVal==nil){
		errorhandler.ErrorToken(expr.Name, "Only instances can have properties");
		return nil;
	}
	ClassVal.Set(expr.Name, interpret.Evaluate(expr.Value));
	return interpret.Evaluate(expr.Value)
}

func (interpreter Interpreter) variableLookup(name token.Token, expr expr.Expr) (interface{}, error){
	dist, ok := interpreter.vars[expr];
	if(!ok){
		errorhandler.ErrorToken(name, "Cant find variable");
		return nil, fmt.Errorf("Cant find variable");
	}
	val, err := interpreter.Env.GetAt(dist, name.Lexeme);
	if(err!=nil){
		errorhandler.ErrorToken(name, "Cant find variable");
		return nil, fmt.Errorf("Cant find variable");

	}
	return val, nil;

}

func isTruthy(object interface{}) bool {
    if object == nil {
        return false
    }

    switch v := object.(type) {
    case bool:
        return v
    case int, int8, int16, int32, int64:
        return v != 0
    case uint, uint8, uint16, uint32, uint64:
        return v != 0
    case float32:
        return v != 0.0
    case float64:
        return v != 0.0
    }

    return true
}

func isEqual(a, b interface{}) bool {
    if a == nil && b == nil {
        return true
    }
    if a == nil || b == nil {
        return false
    }

    return reflect.DeepEqual(a ,b)
}

func checkNumberOperands(token token.Token ,operands ...interface{}) bool{
	for _, val := range operands{
		if _, ok:= val.(float64); !ok{
			errorhandler.ErrorToken(token, "Operand must be a number;");
			return false;
		}
	}
	return true;
}

func stringify(object interface{}) string {
    if object == nil {
        return "nil"
    }

    switch v := object.(type) {
    case float64:
        text := fmt.Sprintf("%g", v)
        return text
    default:
        return fmt.Sprintf("%v", object)
    }
}
