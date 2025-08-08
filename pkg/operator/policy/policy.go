// Package policy inspects and verifies any potential state drift in order to
// either allow the reconciliation loop to continue, or cancel it. Only valid
// state drifts can be applied within the underlying infrastructure.
package policy

import (
	"fmt"

	"github.com/0xSplits/kayron/pkg/cache"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"
)

type Config struct {
	Cac *cache.Cache
	Log logger.Interface
}

type Policy struct {
	cac *cache.Cache
	log logger.Interface
}

func New(c Config) *Policy {
	if c.Cac == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Cac must not be empty", c)))
	}
	if c.Log == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Log must not be empty", c)))
	}

	return &Policy{
		cac: c.Cac,
		log: c.Log,
	}
}
