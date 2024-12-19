package parser

import (
	errorhandler "github.com/vishrudh-raj-rs-14/lox/errorHandler"
	"github.com/vishrudh-raj-rs-14/lox/stmt"
	"github.com/vishrudh-raj-rs-14/lox/token"
)

type Parser struct {
	tokens []token.Token
	current int	
}

func CreateParser(tokens []token.Token) *Parser{
	return &Parser{
		tokens: tokens,
		current: 0,
	}
}

func (parser *Parser) Parse() []stmt.Stmt{
	statements := make([]stmt.Stmt, 0);
	for;!parser.endOfFile();{
		val, err := parser.declaration()
		if(errorhandler.HadError || err!=nil){
			break;
		}
		statements = append(statements, val);
	}
	return statements;
}

