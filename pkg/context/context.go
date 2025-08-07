package context

import (
	"fmt"
	"sync"

	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"
)

type Config struct {
	Log logger.Interface
}

type Context struct {
	inf []Object
	log logger.Interface
	mut sync.Mutex
	ser []Object
}

// The cache implementations used here are the dumb pipes connecting our
// business logic across different boundary conditions by providing critical
// information required to make the entire operator chain work as expected.
// The operator's worker handler is executed constantly in an interval, but
// the real work may only be done if the cached information provided deeper
// down the stack is sufficient. E.g. an empty cache result may cause some
// business logic to be skipped temporarily. Note that the cache keys setup
// here are cross references to the respective cache values within the various
// cache implementations.
func New(c Config) *Context {
	if c.Log == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Log must not be empty", c)))
	}

	return &Context{
		inf: nil,
		log: c.Log,
		mut: sync.Mutex{},
		ser: nil,
	}
}
