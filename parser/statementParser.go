package parser

import (
	"fmt"

	errorhandler "github.com/vishrudh-raj-rs-14/lox/errorHandler"
	"github.com/vishrudh-raj-rs-14/lox/expr"
	"github.com/vishrudh-raj-rs-14/lox/stmt"
	"github.com/vishrudh-raj-rs-14/lox/token"
)

func (parser *Parser) declaration() (stmt.Stmt, error){
	var err error;
	var statement stmt.Stmt;
	if(parser.match(token.VAR)){
		statement, err = parser.varDeclaration();
	}else{
		statement, err = parser.statement();
	}
	if(err==nil){
		return statement, nil;
	}else{
		return nil, err;
	}
}


func (parser *Parser) varDeclaration() (stmt.Stmt, error){
		
	if(parser.match(token.IDENTIFIER)){
		varName := parser.previous();
		if(parser.match(token.EQUAL)){
			intialValue, err := parser.expression();
			if(err!=nil){
				return nil, err;
			}
			statement := stmt.VarStatement{
				Name:varName,
				Value: intialValue,
			}
			err = parser.consume(token.SEMICOLON, "Expected ; after declaration");
			if(err!=nil){
				return nil, err;
			}
			return &statement, nil;
		}
		err := parser.consume(token.SEMICOLON, "Expected ; after declaration");
		if(err!=nil){
			return nil, err;
		}

		statement := stmt.VarStatement{
			Name:varName,
			Value: nil,
		}
		return &statement, nil
}
	errorhandler.ErrorToken(parser.peek(), "Expected identifier name");
	return nil, fmt.Errorf("Expected identifier name");

}

func (parser *Parser) block() (stmt.Stmt, error){

	var statements []stmt.Stmt;
	for;(!parser.check(token.RIGHT_BRACE) && !parser.endOfFile());{
		statement, err := parser.declaration();
		if(err!=nil){
			return nil, err;
		}
		statements = append(statements, statement);
	}
	err:= parser.consume(token.RIGHT_BRACE, "Expected } at end of block");
	if(err!=nil){
		return nil, err;
	}

	return &stmt.BlockStmt{
		Statements: statements,
	}, nil

}


func (parser *Parser) statement() (stmt.Stmt, error){
	if(parser.match(token.PRINT)){
		return parser.printStatement();
	}else if(parser.match(token.LEFT_BRACE)){
		val, err := parser.block();
		return val, err;
	}else if(parser.match(token.IF)){
		return parser.ifStatement();
	}else if(parser.match(token.WHILE)){
		return parser.whileStatement();
	}else if(parser.match(token.FOR)){
		return parser.forStatement();
	}else if(parser.match(token.FUN)){
		return parser.function();
	}else if(parser.match(token.RETURN)){
		return parser.returnStatement();
	}

	return parser.expressionStatement()
}

func (parser *Parser) returnStatement() (stmt.Stmt, error){

	keyword := parser.previous();
	var value expr.Expr = nil;
	var err error;
	if(!parser.check(token.SEMICOLON)){
		value, err = parser.expression();
		if(err!=nil){
			return nil, err;
		}
	}
	parser.consume(token.SEMICOLON, "Expected ; after return value");
	return &stmt.Return{
		Expression: value,
		Keyword: keyword,
	}, nil

}


func (parser *Parser) function() (stmt.Stmt, error){
	
	if(parser.check(token.IDENTIFIER)){
		name := parser.advance(); 

		err := parser.consume(token.LEFT_PAREN, "Expected ( after function name");
		if(err!=nil){
			return nil, err;
		}
		var parameters []token.Token;
		if(!parser.check(token.RIGHT_PAREN)){
		for;parser.match(token.IDENTIFIER);{
			parameters = append(parameters, parser.previous());
			if(parser.check(token.RIGHT_PAREN)){
				parser.advance();
				break;
			}
			err := parser.consume(token.COMMA, "Expected ,");
			if(err!=nil){
				return nil, err;
			}
			if(len(parameters)>=255){
				errorhandler.ErrorToken(parser.peek(), "Can't have more than 255 arguments");
				return nil, fmt.Errorf("Can't have more than 255 arguments")
			}
		}}else{
			err = parser.consume(token.RIGHT_PAREN, "Expected ) at the end of function arguments");
		if(err!=nil){
			return nil, err;
		}
		}
		err = parser.consume(token.LEFT_BRACE, "Expected { at the beginning of function block");
		if(err!=nil){
			return nil, err;
		}
		body, err := parser.block();
		if(err!=nil){
			return nil, err;
		}
		blockBody, ok := body.(*stmt.BlockStmt)
		if !ok {
			errorhandler.ErrorToken(parser.peek(), "Error parsing the function body");
			return nil, fmt.Errorf("Error parsing the function body")
		}
		return &stmt.FunStmt{
			Name:name,
			Parameters: parameters,
			Body: *blockBody,

		}, nil



	}else{
		errorhandler.ErrorToken(parser.peek(), "Expected Identifer");
		return nil, fmt.Errorf("Expected Identifer")
	}

}

func (parser *Parser) forStatement() (stmt.Stmt, error){
	err:= parser.consume(token.LEFT_PAREN, "Expect ( after if");
	if(err!=nil){
		return nil, err;
	}
	
	var initializer stmt.Stmt;

	if(parser.match(token.SEMICOLON)){
		initializer=nil;
	}else if(parser.match(token.VAR)){
		initializer, err=parser.varDeclaration();
		if(err!=nil){
			return nil, err;
		}
	}else {
		initializer, err = parser.expressionStatement();
		if(err!=nil){
			return nil, err;
		}
	}
	var condition expr.Expr;

	if(parser.match(token.SEMICOLON)){
		condition=nil;
	}else {
		condition, err = parser.expression();
		if(err!=nil){
			return nil, err;
		}
	}
	
	parser.consume(token.SEMICOLON, "Expected ; after condition");

	var increment expr.Expr;

	if(parser.check(token.RIGHT_PAREN)){
		increment=nil;
	}else {
		increment, err = parser.expression();
		if(err!=nil){
			return nil, err;
		}
	}

	parser.consume(token.RIGHT_PAREN, "Expected ) after condition");

	body, err := parser.statement();
	if(err!=nil){
		return nil, err;
	}

	if increment != nil {
		body = &stmt.BlockStmt{
			Statements: []stmt.Stmt{
				body,
				&stmt.ExpressionStmt{Expression: increment},
			},
		}
	}

	if condition == nil {
		condition = &expr.Literal{
			Value: true,
		}
	}
	body = &stmt.WhileStmt{
		Expression: condition,
		Body: body,
	}
	if initializer != nil {
		body = &stmt.BlockStmt{
			Statements: []stmt.Stmt{initializer, body},
		}
	}

	return body, nil;

}

func (parser *Parser) whileStatement() (stmt.Stmt, error){
	err:= parser.consume(token.LEFT_PAREN, "Expect ( after if");
	if(err!=nil){
		return nil, err;
	}
	expression, err :=parser.expression();
	if(err!=nil){
		return nil, err;
	}
	err = parser.consume(token.RIGHT_PAREN, "Expect ) at end of expression");
	if(err!=nil){
		return nil, err;
	}
	statement, err := parser.statement();
	if(err!=nil){
		return nil, err;
	}

	return &stmt.WhileStmt{
		Expression: expression,
		Body: statement,
	}, nil


}

func (parser *Parser) ifStatement() (stmt.Stmt, error){
	err:= parser.consume(token.LEFT_PAREN, "Expect ( after if");
	if(err!=nil){
		return nil, err;
	}
	expression, err :=parser.expression();
	if(err!=nil){
		return nil, err;
	}
	err = parser.consume(token.RIGHT_PAREN, "Expect ) at end of expression");
	if(err!=nil){
		return nil, err;
	}
	statement, err := parser.statement();
	if(err!=nil){
		return nil, err;
	}
	var elseBlock stmt.Stmt;
	if(parser.match(token.ELSE)){
		elseBlock , err = parser.statement();
		if(err!=nil){
			return nil, err;
		}
	}

	return &stmt.IfStmt{
		Expression: expression,
		ThenBlock: statement,
		ElseBlock: elseBlock,
	}, nil


}

func (parser *Parser) printStatement() (stmt.Stmt, error){
	expr, err := parser.expression();
	if(err!=nil){
		return nil, err;
	}
	err = parser.consume(token.SEMICOLON, "Expected ; after value");
	if(err!=nil){
		return nil, err;
	}
	statement := &stmt.PrintStmt{
		Expression: expr,
	}
	return statement, nil;

}

func (parser *Parser) expressionStatement() (stmt.Stmt, error){
	expr, err := parser.expression();
	if(err!=nil){
		return nil, err;
	}
	err = parser.consume(token.SEMICOLON, "Expected ; after value");
	if(err!=nil){
		return nil, err;
	}
	statement := &stmt.ExpressionStmt{
		Expression: expr,
	}
	return statement, nil;
}

