package service

import (
	"errors"

	"github.com/xh3b4sd/tracer"
)

var serviceDockerEmptyError = &tracer.Error{
	Kind: "serviceDockerEmptyError",
	Desc: "The service configuration requires a docker repository to be provided.",
}

func IsServiceDockerEmpty(err error) bool {
	return errors.Is(err, serviceDockerEmptyError)
}

var serviceGithubEmptyError = &tracer.Error{
	Kind: "serviceGithubEmptyError",
	Desc: "The service configuration requires a github repository to be provided.",
}

func IsServiceGithubEmpty(err error) bool {
	return errors.Is(err, serviceGithubEmptyError)
}
