package image

import (
	"github.com/xh3b4sd/tracer"
)

var invalidImageTagError = &tracer.Error{
	Description: "This critical error indicates that the provided image tag does neither represent an image index nor an image itself, which means that the operator does not know how to cleanup the given image tag.",
}
