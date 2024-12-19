package environment

import (
	"fmt"

	errorhandler "github.com/vishrudh-raj-rs-14/lox/errorHandler"
	"github.com/vishrudh-raj-rs-14/lox/token"
)

type Environment struct{
	Values map[string]interface{}
	Enclosing *Environment
}

func (env *Environment) Define(key string, val interface{}){
	env.Values[key] = val;
}

func (env *Environment) Get(key string) (interface{}, error){
	val, ok := env.Values[key];
	if(!ok){
		if(env.Enclosing==nil){
			return nil, fmt.Errorf("Value not defined");
		}
		return env.Enclosing.Get(key);
	}
	return val, nil;
}

func (env *Environment) Assign(key token.Token, value interface{}) (error){
	_, ok := env.Values[key.Lexeme];
	if(!ok){
		if(env.Enclosing!=nil){
			return env.Enclosing.Assign(key, value);
		}
		errorhandler.ErrorToken(key, "Variable not declared");
		return fmt.Errorf("Value not defined");
	}
	env.Values[key.Lexeme] = value;
	return nil
}

func (env *Environment) GetAt(dist int, name string) (interface{}, error){
	val, ok := env.ancestor(dist).Values[name];
	if(!ok){
		return nil, fmt.Errorf("Variable not declared")		
	}
	return val, nil;
}


func (env *Environment) AssignAt(dist int, name string, value interface{}){
	env.ancestor(dist).Values[name] = value;
}

func (env *Environment) ancestor(dist int) *Environment{
	cur := env;
	for i:=0;i<dist;i++{
		cur = cur.Enclosing;
	}
	return cur;
}