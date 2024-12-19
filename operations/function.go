package operations

import (
	"github.com/vishrudh-raj-rs-14/lox/environment"
	"github.com/vishrudh-raj-rs-14/lox/stmt"
)

type Function struct{
	FunctionDeclaration stmt.FunStmt
	Closure environment.Environment
	parameterCount int
}


type Return struct {
	Value interface{}
}

func (function *Function) Call(interpreterVal interface{}, arguments []interface{}) (res interface{}){
	interpreter, ok := interpreterVal.(Interpreter) // Type assert to *Interpreter
    if !ok {
        panic("Expected Interpreter as the interpreter")
    }
	newEnvironment := environment.Environment{
		Values: make(map[string]interface{}),
		Enclosing: &function.Closure,
	}

	for i:=0;i<len(function.FunctionDeclaration.Parameters);i++{
		newEnvironment.Define(function.FunctionDeclaration.Parameters[i].Lexeme, arguments[i]);
	}
	res = nil;
	defer func() {
		if r := recover(); r != nil {
			// fmt.Println(r);
			if returnValue, ok := r.(Return); ok {
				res = returnValue.Value
			}
		}
	}()

	interpreter.ExecuteBlock(&function.FunctionDeclaration.Body, &newEnvironment);
	return;

}

func (function *Function) Arity() int{
	return function.parameterCount
}


func (function *Function) String() string{
	return "< fn " + function.FunctionDeclaration.Name.Lexeme + " >"
}