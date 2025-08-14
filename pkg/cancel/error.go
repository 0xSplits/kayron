package cancel

import (
	"errors"

	"github.com/xh3b4sd/tracer"
)

var Error = &tracer.Error{
	Description: "This error is used as control flow signal to instruct the worker engine that the current reconciliation loop should not execute any further.",
}

// Is is used to prevent cancel errors to be propagated inside the worker
// engine's logs and metrics.
func Is(err error) bool {
	return errors.Is(err, Error)
}
