package roghfs

import (
	"errors"

	"github.com/xh3b4sd/tracer"
)

var fileAlreadyCachedError = &tracer.Error{
	Description: "This critical error indicates that the cache logic of the file system is broken, because we ended up caching a file that was supposed to already be cached.",
}

func IsFileAlreadyCached(err error) bool {
	return errors.Is(err, fileAlreadyCachedError)
}

//
//
//

var invalidRepositoryUrlError = &tracer.Error{
	Description: "This runtime error indicates that the provided Github URL was malformed, because we could not parse the repository owner and repository name.",
	Context: []tracer.Context{
		{Key: "expected", Value: "https://github.com/owner/repo"},
	},
}

func IsInvalidRepositoryUrl(err error) bool {
	return errors.Is(err, invalidRepositoryUrlError)
}
