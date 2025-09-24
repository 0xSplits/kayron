// Package logging emits log messages for each and every reconciliation loop
// about the most important aspects of the current state of the system.
package logging

import (
	"fmt"

	"github.com/0xSplits/kayron/pkg/envvar"
	"github.com/0xSplits/kayron/pkg/policy"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"
)

type Config struct {
	Env envvar.Env
	Log logger.Interface
	Pol *policy.Policy
}

type Logging struct {
	env envvar.Env
	log logger.Interface
	pol *policy.Policy
}

func New(c Config) *Logging {
	if c.Env.Environment == "" {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Env must not be empty", c)))
	}
	if c.Log == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Log must not be empty", c)))
	}
	if c.Pol == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Pol must not be empty", c)))
	}

	return &Logging{
		env: c.Env,
		log: c.Log,
		pol: c.Pol,
	}
}
