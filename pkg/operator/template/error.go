package template

import (
	"github.com/xh3b4sd/tracer"
)

var missingRootStackError = &tracer.Error{
	Description: "The operator expected to find exactly one root stack, but no stack was found for the given environment.",
}
