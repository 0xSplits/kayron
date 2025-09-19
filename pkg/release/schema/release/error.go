package release

import (
	"errors"

	"github.com/xh3b4sd/tracer"
)

var releaseDefinitionEmptyError = &tracer.Error{
	Description: "The release configuration requires at least one release definition to be provided.",
	Context: []tracer.Context{
		{Key: "reason", Value: "the loader could not find any release definitions"},
	},
}

func IsServiceDefinitionEmpty(err error) bool {
	return errors.Is(err, releaseDefinitionEmptyError)
}

//
//
//

var releaseDeployEmptyError = &tracer.Error{
	Description: "The release configuration requires a deployment strategy to be provided.",
}

func IsServiceDeployEmpty(err error) bool {
	return errors.Is(err, releaseDeployEmptyError)
}

//
//
//

var releaseDeployPreviewError = &tracer.Error{
	Description: "The release configuration does not allow preview deployments for infrastructure providers.",
}

func IsServiceDeployPreview(err error) bool {
	return errors.Is(err, releaseDeployPreviewError)
}

//
//
//

var releaseGithubEmptyError = &tracer.Error{
	Description: "The release configuration requires a github repository to be provided.",
}

func IsServiceGithubEmpty(err error) bool {
	return errors.Is(err, releaseGithubEmptyError)
}

//
//
//

var releaseLabelsEmptyError = &tracer.Error{
	Description: "The release configuration requires all internal labels to be provided.",
}

func IsServiceLabelsEmpty(err error) bool {
	return errors.Is(err, releaseLabelsEmptyError)
}

//
//
//

var releaseProviderEmptyError = &tracer.Error{
	Description: "The release configuration requires a docker repository or provider setting to be provided.",
	Context: []tracer.Context{
		{Key: "reason", Value: "the loader could neither find a docker repository nor provider setting"},
	},
}

func IsServiceProviderEmpty(err error) bool {
	return errors.Is(err, releaseProviderEmptyError)
}
