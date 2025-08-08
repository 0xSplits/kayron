package policy

import (
	"errors"

	"github.com/xh3b4sd/tracer"
)

var cancelError = &tracer.Error{
	Description: "This error is used as control flow signal to instruct the worker engine that the current reconciliation loop should not execute any further.",
}

// IsCancel is used to prevent cancel errors to be propagated inside the worker
// engine's logs and metrics.
func IsCancel(err error) bool {
	return errors.Is(err, cancelError)
}

//
//
//

var cacheStateEmptyError = &tracer.Error{
	Description: "This critical error indicates that some cached state of the current reconciliation loop was missing, which means that the operator does not know how to proceed safely.",
}
