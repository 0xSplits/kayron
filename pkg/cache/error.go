package cache

import (
	"github.com/xh3b4sd/tracer"
)

var invalidInfrastructureError = &tracer.Error{
	Description: "This critical error indicates that the resolved release settings do not provide a single infrastructure release, which means that the operator does not know how to proceed safely.",
}
