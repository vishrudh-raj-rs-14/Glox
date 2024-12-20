package operations

import (
	"github.com/vishrudh-raj-rs-14/lox/environment"
	"github.com/vishrudh-raj-rs-14/lox/stmt"
)

type Function struct{
	FunctionDeclaration stmt.FunStmt
	Closure environment.Environment
	parameterCount int
	IsInitializer bool
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
			if(function.IsInitializer){
				val, _ := function.Closure.GetAt(0, "this");
				res = val;
			}
			// fmt.Println(r);
			if returnValue, ok := r.(Return); ok {
				res = returnValue.Value
			}
		}
	}()

	interpreter.ExecuteBlock(&function.FunctionDeclaration.Body, &newEnvironment);
	if(function.IsInitializer){
		val, _ := function.Closure.GetAt(0, "this");
		return val;
	}
	return;

}

func (fun *Function) Bind(instance *Instance) *Function{
	environment := environment.Environment{
		Values: make(map[string]interface{}),
		Enclosing: &fun.Closure,
	}
	environment.Define("this", instance);
	return &Function{
		FunctionDeclaration: fun.FunctionDeclaration,
		parameterCount: len(fun.FunctionDeclaration.Parameters),
		Closure: environment,
		IsInitializer: fun.IsInitializer,
	}
}


func (function *Function) Arity() int{
	return function.parameterCount
}


func (function *Function) String() string{
	return "< fn " + function.FunctionDeclaration.Name.Lexeme + " >"
}