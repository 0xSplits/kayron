package artifact

import (
	"errors"

	"github.com/xh3b4sd/tracer"
)

var serviceNotCachedError = &tracer.Error{
	Description: "This critical error indicates that the cache logic of the worker handler is broken, because we ended up missing a service that was supposed to already be cached.",
}

func IsServiceNotCached(err error) bool {
	return errors.Is(err, serviceNotCachedError)
}
