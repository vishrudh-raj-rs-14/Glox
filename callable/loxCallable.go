package callable

type LoxCallable interface{
	Call(interpreter interface{}, arguments []interface{}) (res interface{})
	Arity() int
	String() string
}
