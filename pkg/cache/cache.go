package cache

import (
	"fmt"
	"sync"

	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"
)

type Config struct {
	// Frc is to forcefully trigger an infrastructure deployment, regardless of
	// any detectable state drift. This option should be used with caution. E.g.
	// the command "kayron daemon" should never allow to apply this trigger flag
	// within Kayron's normal reconciliation loop.
	Frc bool

	Log logger.Interface
}

// Cache represents the dumb pipes connecting our business logic across
// different boundary conditions by providing critical information required to
// make the entire operator chain work as expected. The operator's worker
// handler is executed constantly in an interval, but the real work may only be
// done if the cached information about the specific release artifacts provided
// deeper down the stack is sufficient. E.g. an empty cache result may cause
// some business logic to be skipped temporarily, or cause the reconciliation
// loop to be cancelled entirely.
type Cache struct {
	frc bool
	inf []Object
	log logger.Interface
	mut sync.Mutex
	ser []Object
}

func New(c Config) *Cache {
	if c.Log == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Log must not be empty", c)))
	}

	return &Cache{
		frc: c.Frc,
		inf: nil,
		log: c.Log,
		mut: sync.Mutex{},
		ser: nil,
	}
}
