package runtime

import (
	"github.com/xh3b4sd/tracer"
)

var ExecutionFailedError = &tracer.Error{
	Description: "This internal error implies a severe malfunction of the system.",
}
