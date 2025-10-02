package github

import (
	"github.com/xh3b4sd/tracer"
)

var invalidGithubRepositoryError = &tracer.Error{
	Description: "This critical error indicates that the provided repository name does not comply with the required format enforced by Github, which means that the operator does not know how to proceed safely.",
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
