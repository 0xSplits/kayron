package webhook

import (
	"fmt"
	"sync"

	"github.com/0xSplits/kayron/pkg/envvar"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"
)

const (
	// branch is the ref prefix for commits pushed to branches.
	branch = "refs/heads/"
)

type Config struct {
	Env envvar.Env
	Log logger.Interface
}

type Webhook struct {
	cac map[string]Commit
	env envvar.Env
	log logger.Interface
	mut sync.Mutex
}

func New(c Config) *Webhook {
	if c.Env.Environment == "" {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Env must not be empty", c)))
	}
	if c.Log == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Log must not be empty", c)))
	}

	return &Webhook{
		cac: map[string]Commit{},
		env: c.Env,
		log: c.Log,
		mut: sync.Mutex{},
	}
}
