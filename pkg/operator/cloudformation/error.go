package cloudformation

import (
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
)

func isNoStateDriftError(err error) bool {
	ioe, mat := err.(*types.InvalidOperationException)
	if !mat {
		return false
	}

	if ioe.ErrorCode() == "ValidationError" && ioe.ErrorMessage() == "No updates are to be performed." {
		return true
	}

	return false
}
