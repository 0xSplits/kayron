package policy

import (
	"github.com/xh3b4sd/tracer"
)

var cacheStateEmptyError = &tracer.Error{
	Description: "This critical error indicates that some cached state of the current reconciliation loop was missing, which means that the operator does not know how to proceed safely.",
}
