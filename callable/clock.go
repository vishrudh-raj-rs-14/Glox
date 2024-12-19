package callable

import (
	"time"
)

type ClockFunction struct{}

func (c *ClockFunction) Arity() int {
	return 0
}

func (c *ClockFunction) Call(interpreter interface{}, arguments []interface{}) interface{} {
	return float64(time.Now().UnixNano()) / 1e9
}

func (c *ClockFunction) String() string {
	return "<native fn>"
}