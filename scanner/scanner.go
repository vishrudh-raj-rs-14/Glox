package scanner

import (
	"strconv"

	errorhandler "github.com/vishrudh-raj-rs-14/lox/errorHandler"
	"github.com/vishrudh-raj-rs-14/lox/token"
)

var Keywords map[string]token.TokenType = map[string]token.TokenType{
	"and":    token.AND,
	"class":  token.CLASS,
	"else":   token.ELSE,
	"false":  token.FALSE,
	"for":    token.FOR,
	"fun":    token.FUN,
	"if":     token.IF,
	"nil":    token.NIL,
	"or":     token.OR,
	"print":  token.PRINT,
	"return": token.RETURN,
	"super":  token.SUPER,
	"this":   token.THIS,
	"true":   token.TRUE,
	"var":    token.VAR,
	"while":  token.WHILE,
}

type Scanner struct {
	source string
	tokens []token.Token
	start int
	current int
	line int
}

func (scanner *Scanner) ScanTokens() []token.Token{
	for ;!scanner.endOfFile();{
		scanner.start = scanner.current;
		scanner.ScanToken();
	}
	scanner.tokens = append(scanner.tokens, token.NewToken(token.EOF, "", nil, scanner.line));
	return scanner.tokens;
}

func (scanner *Scanner) ScanToken(){
	c := scanner.advance();
	switch c {
	case '{':
		scanner.addToken(token.LEFT_BRACE, nil);
	case '}':
		scanner.addToken(token.RIGHT_BRACE, nil);
	case '(':
		scanner.addToken(token.LEFT_PAREN, nil);
	case ')':
		scanner.addToken(token.RIGHT_PAREN, nil);
	case ',':
		scanner.addToken(token.COMMA, nil);
	case '.':
		scanner.addToken(token.DOT, nil);
	case '+':
		scanner.addToken(token.PLUS, nil);
	case '-':
		scanner.addToken(token.MINUS, nil);
	case '%':
		scanner.addToken(token.MOD, nil);
	case ';':
		scanner.addToken(token.SEMICOLON, nil);
	case '*':
		scanner.addToken(token.STAR, nil);
	case '!':
		if(scanner.match('=')){
			scanner.addToken(token.BANG_EQUAL, nil);
		}else{
			scanner.addToken(token.BANG, nil);
		}
	case '=':
		if(scanner.match('=')){
			scanner.addToken(token.EQUAL_EQUAL, nil);
		}else{
			scanner.addToken(token.EQUAL, nil);
		}
	case '<':
		if(scanner.match('=')){
			scanner.addToken(token.LESS_EQUAL, nil);
		}else{
			scanner.addToken(token.LESS, nil);
		}
	case '>':
		if(scanner.match('=')){
			scanner.addToken(token.GREATER_EQUAL, nil);
		}else{
			scanner.addToken(token.GREATER, nil);
		}
	case '/':
		if(scanner.match('/')){
			for ;scanner.peek()!='\n' && !scanner.endOfFile();{
				scanner.advance();
			}
		}else if(scanner.match('*')){
			var depth int = 1;
			for ;!scanner.endOfFile();{
				if(scanner.peek()=='\n'){
					scanner.line++;
				}
				if(scanner.peek()=='/'&&scanner.nextPeek()=='*'){
					scanner.advance();
					scanner.advance();
					depth++;
					continue;
				}
				if(scanner.peek()=='*'&&scanner.nextPeek()=='/'){
					scanner.advance();
					scanner.advance();
					depth--;
					if(depth)==0{
						break;
					}
				}
				scanner.advance();
			}
		}else{
			scanner.addToken(token.SLASH, nil);
		}
	case '\n':
		scanner.line++;
	case ' ':
	case '\r':
	case '\t':
	case '"':
		scanner.getString();
	
	default:
		if(isDigit(c)){
			scanner.getNumber();
		}else if(isAlpha(c)){
			scanner.getIdentifier();
		}else{
			errorhandler.Error(scanner.line, "Unexpected symbol " + string(c));
		}
	}

}

func (scanner *Scanner) getString(){
	for;scanner.peek()!='"' && !scanner.endOfFile();{
		s := scanner.advance();
		if(s=='\n'){
			scanner.line++;
		}
	}
	if(scanner.endOfFile()) {
		errorhandler.Error(scanner.line, "Unexpected termination of string");
	}
	scanner.advance();
	scanner.addToken(token.STRING, scanner.source[scanner.start+1: scanner.current-1]);
}

func (scanner *Scanner) getNumber(){
	for;isDigit(scanner.peek()) && !scanner.endOfFile();{
		scanner.advance()
	}
	if(scanner.peek()=='.' && isDigit(scanner.nextPeek())){
		scanner.advance();
		for;isDigit(scanner.peek()) && !scanner.endOfFile();{
			scanner.advance()
		}
	}
	num, _ :=  strconv.ParseFloat(scanner.source[scanner.start:scanner.current], 64);
	scanner.addToken(token.NUMBER, num);
}

func (scanner *Scanner) getIdentifier(){
	for;isAlphaNum(scanner.peek()) && !scanner.endOfFile();{
		scanner.advance();
	}
	text := scanner.source[scanner.start:scanner.current];
	typeVal, ok := Keywords[text];
	if(!ok){
		typeVal = token.IDENTIFIER;
	}
	scanner.addToken(typeVal, nil);
}

func (scanner *Scanner) endOfFile() bool {
	return scanner.current >= len(scanner.source);
}

// consume current char and return it.
func (scanner *Scanner) advance() byte {
	curChar := scanner.source[scanner.current];
	scanner.current++;
	return (curChar);
}

// checks if current char matches next byte
func (scanner *Scanner) match(next byte) bool {
	if(scanner.endOfFile()) {
		return false;
	}
	if(scanner.source[scanner.current] != next){
		return false;
	}
	scanner.current++;
	return true;
}

// only returns current char without consuming it.
func (scanner *Scanner) peek() byte {
	if(scanner.endOfFile()){
		return byte(0);
	}
	return scanner.source[scanner.current];
}

// peeks at the next character(current+1)
func (scanner *Scanner) nextPeek() byte {
	if(scanner.current+1>=len(scanner.source)){
		return byte(0);
	}
	return scanner.source[scanner.current+1];
}

// adds token to the list of tokens
func (scanner *Scanner) addToken(tokenType token.TokenType, literal interface{}) {
	text := scanner.source[scanner.start:scanner.current];
	scanner.tokens = append(scanner.tokens, token.NewToken(tokenType, text, literal, scanner.line));
}

func isDigit(c byte) bool {
	return c>='0' && c<='9';
}

func isAlpha(c byte) bool {
	return (c>='a' && c<='z')||(c>='A' && c<='Z')||c=='_';
}

func isAlphaNum(c byte) bool {
	return isAlpha(c) || isDigit(c);
}

func CreateScanner (src string)  *Scanner{
	return &Scanner{
		source: src,
		start: 0,
		current: 0,
		line: 1,
	}
}