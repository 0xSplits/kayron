// Package policy inspects and verifies any potential state drift in order to
// either allow the reconciliation loop to continue, or cancel it. Only valid
// state drifts can be applied within the underlying infrastructure.
package policy

import (
	"fmt"

	"github.com/0xSplits/kayron/pkg/context"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"
)

type Config struct {
	Ctx *context.Context
	Log logger.Interface
}

type Policy struct {
	ctx *context.Context
	log logger.Interface
}

func New(c Config) *Policy {
	if c.Ctx == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Ctx must not be empty", c)))
	}
	if c.Log == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Log must not be empty", c)))
	}

	return &Policy{
		ctx: c.Ctx,
		log: c.Log,
	}
}
