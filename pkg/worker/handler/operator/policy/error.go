package policy

import (
	"errors"

	"github.com/xh3b4sd/tracer"
)

var Cancel = &tracer.Error{
	Description: "This error is used as control flow signal to instruct the worker engine that the current reconciliation loop should not execute any further.",
}

func IsCancel(err error) bool {
	return errors.Is(err, Cancel)
}
