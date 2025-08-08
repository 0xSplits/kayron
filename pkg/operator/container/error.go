package container

import (
	"github.com/xh3b4sd/tracer"
)

var invalidAmazonResourceNameError = &tracer.Error{
	Description: "The exporter expected the ARN format to be [arn:aws:ecs:<region>:<account>:service/<cluster>/<service>], but a different format was found.",
}

var invalidEcsServiceError = &tracer.Error{
	Description: "This critical error indicates that the query for a single ECS service did not result in exactly one service result in the response.",
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
