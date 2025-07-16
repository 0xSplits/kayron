package roghfs

import (
	"github.com/xh3b4sd/tracer"
)

var fileAlreadyCachedError = &tracer.Error{
	Kind: "fileAlreadyCachedError",
	Desc: "This critical error indicates that the cache logic of the file system is broken, because we ended up caching a file that was supposed to already be cached.",
}
