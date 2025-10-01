package release

import (
	"github.com/xh3b4sd/tracer"
)

var invalidReleaseFormatError = &tracer.Error{
	Description: "This critical error indicates that the provided release tag does not comply with the required format [v.MAJOR.MINOR.PATCH(-SUFFIX)], which means that the operator does not know how to proceed safely.",
}

//
//
//

func isErr(err error) bool {
	return err != nil
}

func isNil(err error) bool {
	return err == nil
}
