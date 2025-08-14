package stack

import (
	"github.com/xh3b4sd/tracer"
)

var invalidRootStackError = &tracer.Error{
	Description: "This critical error indicates that the operator could not find the configured root stack for the given environment, which means that the operator does not know how to proceed safely.",
}
