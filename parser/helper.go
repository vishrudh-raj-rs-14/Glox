package parser

import (
	"fmt"

	errorhandler "github.com/vishrudh-raj-rs-14/lox/errorHandler"
	"github.com/vishrudh-raj-rs-14/lox/token"
)



func (parser *Parser) match(tokens ...token.TokenType) bool{
	for _, tokenType := range tokens{
		if(parser.check(tokenType)){
			parser.advance();
			return true;
		}
	}
	return false;
}


func (parser *Parser) consume(tokenType token.TokenType, errMsg string) error{
	if(parser.check(tokenType)){
		parser.advance();
		return nil;
	}
	errorhandler.ErrorToken(parser.peek(), errMsg);
	return fmt.Errorf("unexpected token: %v", parser.peek())
}

func (parser *Parser) check(tokenType token.TokenType) bool{
	if(parser.endOfFile()){
		return false;
	}
	return parser.peek().TokenType==tokenType;
}

func (parser *Parser) peek() token.Token{
	return parser.tokens[parser.current];
}

func (parser *Parser) advance() token.Token{
	if(!parser.endOfFile()){
		parser.current++;
		return parser.tokens[parser.current-1];
	}
	return parser.tokens[parser.current];
}

func (parser *Parser) previous() token.Token{
	
	return parser.tokens[parser.current-1];
}

func (parser *Parser) endOfFile() bool{
	return parser.current>=(len(parser.tokens)-1);
}