package resolver

import (
	"errors"
	"net/http"

	"github.com/google/go-github/v76/github"
	"github.com/xh3b4sd/tracer"
)

var releaseNotFoundError = &tracer.Error{
	Description: "This critical error indicates that neither a production branch nor a Github release exists for the production environment, which means that the operator does not know how to proceed safely.",
	Context: []tracer.Context{
		{Key: "suggestion", Value: "create either a production branch or a Github release for the production environment"},
	},
}

func IsReleaseNotFound(err error) bool {
	return errors.Is(err, releaseNotFoundError)
}

//
//
//

func isNotFound(res *github.Response) bool {
	if res == nil {
		return false
	}

	return res.StatusCode == http.StatusNotFound || res.StatusCode == http.StatusUnprocessableEntity
}
