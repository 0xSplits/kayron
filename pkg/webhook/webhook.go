package webhook

import (
	"fmt"
	"sync"

	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"
)

const (
	// branch is the ref prefix for commits pushed to branches.
	branch = "refs/heads/"
)

type Config struct {
	Log logger.Interface
}

type Webhook struct {
	cac map[Key]Commit
	log logger.Interface
	mut sync.Mutex
}

func New(c Config) *Webhook {
	if c.Log == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Log must not be empty", c)))
	}

	return &Webhook{
		cac: map[Key]Commit{},
		log: c.Log,
		mut: sync.Mutex{},
	}
}
