package canceler

import (
	"github.com/xh3b4sd/tracer"
)

var invalidStackStatusError = &tracer.Error{
	Description: "This critical error indicates that the operator found the configured CloudFormation stack to be in a stack status that was unrecognizable, which means that the operator does not know how to proceed safely.",
}
