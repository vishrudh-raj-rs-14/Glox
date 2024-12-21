package operations

import (
	errorhandler "github.com/vishrudh-raj-rs-14/lox/errorHandler"
	"github.com/vishrudh-raj-rs-14/lox/token"
)

type Class struct{
	Name token.Token
	Methods map[string]*Function
	SuperClass *Class
}


func (klass *Class) Call(interpreterVal interface{}, arguments []interface{}) (res interface{}){

	instance := CreateInstance(*klass);
	initializer := klass.FindMethod("init");
	if(initializer!=nil){
		initializer.Bind(instance).Call(interpreterVal, arguments);
	}
	return instance;
}

func (klass *Class) FindMethod(name string) *Function{
	val, ok:= klass.Methods[name];
	if(!ok){
		if(klass.SuperClass==nil){
			return nil;
		}
		return klass.SuperClass.FindMethod(name);
	}
	return val;
}

func (klass *Class) Arity() int{
	initializer := klass.FindMethod("init");
	if(initializer==nil){
		return 0;
	}
	return initializer.Arity();
}

func (klass *Class) String() string{
	return "<Class " + klass.Name.Lexeme + ">"
}


type Instance struct{
	InstanceOf Class
	fields map[string]interface{}
}

func CreateInstance(class Class) *Instance{
	return &Instance{
		InstanceOf: class,
		fields: make(map[string]interface{}),
	}
}

func (klass *Instance) Get(name token.Token) interface{}{
	val, ok := klass.fields[name.Lexeme];
	if(!ok){
		method := klass.InstanceOf.FindMethod(name.Lexeme);
		if(method==nil){
			errorhandler.ErrorToken(name, "Property does not exist on class");
			return nil;
		}
		return method.Bind(klass);
	}
	return val
}




func (klass *Instance) Set(name token.Token, value interface{}){
	// _, ok := klass.fields[name.Lexeme];
	// if(!ok){
	// 	errorhandler.ErrorToken(name, "Property does not exist on class");
	// }
	klass.fields[name.Lexeme] = value
}


func (instance *Instance) String() string{
	return "<instance of " + instance.InstanceOf.String() + ">"
}