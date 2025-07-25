package resolver

import (
	"errors"

	"github.com/xh3b4sd/tracer"
)

var releaseNotFoundError = &tracer.Error{
	Description: "This critical error indicates that neither a production branch nor a Github release exists for the production environment.",
	Context: []tracer.Context{
		{Key: "suggestions", Value: "create either a production branch or a Github release for the production environment"},
	},
}

func IsReleaseNotFound(err error) bool {
	return errors.Is(err, releaseNotFoundError)
}
