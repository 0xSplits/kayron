package container

import (
	"github.com/xh3b4sd/tracer"
)

var invalidAmazonResourceNameError = &tracer.Error{
	Description: "The exporter expected the ARN format to be [arn:aws:ecs:<region>:<account>:service/<cluster>/<service>], but a different format was found.",
}

var invalidEcsServiceError = &tracer.Error{
	Description: "This critical error indicates that the query for a single ECS service did not yield exactly one service result in the response, which means that the operator does not know how to proceed safely.",
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
