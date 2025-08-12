package template

import (
	"github.com/xh3b4sd/tracer"
)

var invalidRootStackError = &tracer.Error{
	Description: "The operator expected to find exactly one root stack for the configured stack name, but no root stack was found for the given environment.",
}
