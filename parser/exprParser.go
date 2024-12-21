package parser

import (
	"fmt"

	errorhandler "github.com/vishrudh-raj-rs-14/lox/errorHandler"
	"github.com/vishrudh-raj-rs-14/lox/expr"
	"github.com/vishrudh-raj-rs-14/lox/token"
)

func (parser *Parser) expression() (expr.Expr, error) {
	finalExpr, err := parser.assignment()
	if err != nil {
		return nil, err
	}
	return finalExpr, nil
}

func (parser *Parser) assignment() (expr.Expr, error) {
	expression, err := parser.or()
	if(err!=nil){
		return nil, err;
	}
	if(parser.match(token.EQUAL)){
		equals := parser.previous();
		value, err := parser.assignment();
		if(err!=nil){
			return nil, err
		}
		if variable, ok := expression.(*expr.Variable); ok {
            return &expr.Assign{
                Name:  variable.Token, 
                Value: value,         
            }, nil
        }else{
			if variable, ok := expression.(*expr.GetExpr); ok {
				return &expr.SetExpr{
					Name: variable.Name,
					Object: variable.Object,
					Value: value,      
				}, nil
			}
		}
		errorhandler.ErrorToken(equals, "Invalid assignment type");
		return nil, fmt.Errorf("Invalid assignment type");
	}
	return expression, nil
}

func (parser *Parser) or() (expr.Expr, error) {
	expression, err := parser.and()
	if err != nil {
		return nil, err
	}
	for parser.match(token.OR) {
		operator := parser.previous()
		right, err := parser.and()
		if err != nil {
			return nil, err
		}
		expression = &expr.Logical{
			Left:     expression,
			Operator: operator,
			Right:    right,
		}
	}
	return expression, nil
}

func (parser *Parser) and() (expr.Expr, error) {
	expression, err := parser.equality()
	if err != nil {
		return nil, err
	}
	for parser.match(token.AND) {
		operator := parser.previous()
		right, err := parser.equality()
		if err != nil {
			return nil, err
		}
		expression = &expr.Logical{
			Left:     expression,
			Operator: operator,
			Right:    right,
		}
	}
	return expression, nil
}

func (parser *Parser) equality() (expr.Expr, error) {
	expression, err := parser.comparison()
	if err != nil {
		return nil, err
	}
	for parser.match(token.BANG_EQUAL, token.EQUAL_EQUAL) {
		operator := parser.previous()
		right, err := parser.comparison()
		if err != nil {
			return nil, err
		}
		expression = &expr.Binary{
			Left:     expression,
			Operator: operator,
			Right:    right,
		}
	}
	return expression, nil
}

func (parser *Parser) comparison() (expr.Expr, error) {
	expression, err := parser.term()
	if err != nil {
		return nil, err
	}
	for parser.match(token.GREATER_EQUAL, token.LESS_EQUAL, token.GREATER, token.LESS) {
		operator := parser.previous()
		right, err := parser.term()
		if err != nil {
			return nil, err
		}
		expression = &expr.Binary{
			Left:     expression,
			Operator: operator,
			Right:    right,
		}
	}
	return expression, nil
}

func (parser *Parser) term() (expr.Expr, error) {
	expression, err := parser.factor()
	if err != nil {
		return nil, err
	}
	for parser.match(token.MINUS, token.PLUS) {
		operator := parser.previous()
		right, err := parser.factor()
		if err != nil {
			return nil, err
		}
		expression = &expr.Binary{
			Left:     expression,
			Operator: operator,
			Right:    right,
		}
	}
	return expression, nil
}

func (parser *Parser) factor() (expr.Expr, error) {
	expression, err := parser.unary()
	if err != nil {
		return nil, err
	}
	for parser.match(token.STAR, token.SLASH, token.MOD) {
		operator := parser.previous()
		right, err := parser.unary()
		if err != nil {
			return nil, err
		}
		expression = &expr.Binary{
			Left:     expression,
			Operator: operator,
			Right:    right,
		}
	}
	return expression, nil
}

func (parser *Parser) unary() (expr.Expr, error) {
	if parser.match(token.MINUS, token.BANG) {
		operator := parser.previous()
		right, err := parser.call()
		if err != nil {
			return nil, err
		}
		return &expr.Unary{
			Right:    right,
			Operator: operator,
		}, nil
	}
	return parser.call()
}

func (parser *Parser) call() (expr.Expr, error) {
	callee, err := parser.primary();
	if err != nil {
		return nil, err
	}
	for;;{
		if(parser.match(token.LEFT_PAREN)){
			val, err := parser.finishCall(callee);
			if(err!=nil){
				return nil, err;
			}
			callee = val;
		}else if(parser.match(token.DOT)){
			if(!parser.check(token.IDENTIFIER)){
				errorhandler.ErrorToken(parser.peek(), "Expected Identifier")
				return nil, fmt.Errorf("Expected Identifier")
			}
			name := parser.advance();
			callee = &expr.GetExpr{
				Name: name,
				Object: callee,
			}
		}else{
			break;
		}
	}

	return callee, nil;

}


func (parser *Parser) finishCall(callee expr.Expr) (expr.Expr, error) {
    var arguments []expr.Expr

    if !parser.check(token.RIGHT_PAREN) {
        for {
			val, err :=  parser.expression();
			if(err!=nil){
				return nil, err;
			}
            arguments = append(arguments,val)
            if !parser.match(token.COMMA) {
                break
            }
			if(len(arguments)>=255){
				errorhandler.ErrorToken(parser.peek(), "Too many arguments used. max limit is 255");
				return nil, fmt.Errorf("Too many arguments used. max limit is 255")
			}
        }
    }

    err := parser.consume(token.RIGHT_PAREN, "Expect ')' after arguments.")
	if(err!=nil){
		return nil, err
	}
    return &expr.Call{
		Callee: callee,
		Arguments: arguments,
		CloseParen: parser.previous(),
	}, nil
}


func (parser *Parser) primary() (expr.Expr, error) {
	if parser.match(token.FALSE) {
		return &expr.Literal{Value: false}, nil
	}

	if parser.match(token.TRUE) {
		return &expr.Literal{Value: true}, nil
	}

	if parser.match(token.NIL) {
		return &expr.Literal{Value: nil}, nil
	}

	if(parser.match(token.THIS)){
		return &expr.This{Keyword: parser.previous()}, nil;
	}

	if(parser.match(token.SUPER)){
		keyword := parser.previous();
		err := parser.consume(token.DOT, "Super cannot be used alone");
		if(err!=nil){
			return nil, err;
		}
		if(!parser.check(token.IDENTIFIER)){
			errorhandler.ErrorToken(parser.peek(), "Method of super must be a identifier");
			return nil, fmt.Errorf("Method of super must be a identifier");
		}
		method := parser.advance();
		return &expr.Super{
			Keyword: keyword,
			Method: method,
		}, nil
	}

	if parser.match(token.IDENTIFIER) {
		return &expr.Variable{Name: parser.previous().Lexeme, Token: parser.previous()}, nil
	}

	if parser.match(token.STRING, token.NUMBER) {
		return &expr.Literal{Value: parser.previous().Literal}, nil
	}

	if parser.match(token.LEFT_PAREN) {
		expr, err := parser.expression()
		if err != nil {
			return nil, err
		}
		err = parser.consume(token.RIGHT_PAREN, "Expected ')' after expression")
		if err != nil {
			return nil, err
		}
		return expr, nil
	}
	errorhandler.ErrorToken(parser.peek(), "Unexpected token");
	return nil, fmt.Errorf("Expected ')' after expression")
}


