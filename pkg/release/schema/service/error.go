package service

import (
	"errors"

	"github.com/xh3b4sd/tracer"
)

var serviceDefinitionEmptyError = &tracer.Error{
	Description: "The release configuration requires at least one service definition to be provided.",
	Context: []tracer.Context{
		{Key: "reason", Value: "the loader could not find any service definitions"},
	},
}

func IsServiceDefinitionEmpty(err error) bool {
	return errors.Is(err, serviceDefinitionEmptyError)
}

//
//
//

var serviceDeployEmptyError = &tracer.Error{
	Description: "The service configuration requires a deployment strategy to be provided.",
}

func IsServiceDeployEmpty(err error) bool {
	return errors.Is(err, serviceDeployEmptyError)
}

//
//
//

var serviceGithubEmptyError = &tracer.Error{
	Description: "The service configuration requires a github repository to be provided.",
}

func IsServiceGithubEmpty(err error) bool {
	return errors.Is(err, serviceGithubEmptyError)
}

//
//
//

var serviceLabelsEmptyError = &tracer.Error{
	Description: "The service configuration requires all internal labels to be provided.",
}

func IsServiceLabelsEmpty(err error) bool {
	return errors.Is(err, serviceLabelsEmptyError)
}

//
//
//

var serviceProviderEmptyError = &tracer.Error{
	Description: "The service configuration requires a docker repository or provider setting to be provided.",
	Context: []tracer.Context{
		{Key: "reason", Value: "the loader could neither find a docker repository nor provider setting"},
	},
}

func IsServiceProviderEmpty(err error) bool {
	return errors.Is(err, serviceProviderEmptyError)
}
