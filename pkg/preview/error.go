package preview

import (
	"github.com/xh3b4sd/tracer"
)

var containerImageFormatError = &tracer.Error{
	Description: "This critical error indicates that the provided container image was unrecognizable, which means that the operator does not know how to proceed safely.",
}
