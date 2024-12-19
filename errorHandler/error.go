package errorhandler

import (
	"fmt"
	"strconv"

	"github.com/vishrudh-raj-rs-14/lox/token"
)


var HadError bool = false;

func Error(line int, msg string){
	HadError = true;
	Report(line, "", msg);
}

func ErrorToken(tok token.Token, errMsg string){
	if(tok.TokenType==token.EOF){
		Report(tok.Line, " at end" ,errMsg);
	}else{
		Report(tok.Line, " at '" + tok.Lexeme + "'" ,errMsg);
	}

}

func Report(line int, where string, msg string){
	HadError = true;
	fmt.Println("[line " + strconv.Itoa(line) + "] Error" + where + ": " + msg);
}
