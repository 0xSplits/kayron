package policy

import (
	"fmt"

	"github.com/0xSplits/kayron/pkg/cache"
	"github.com/0xSplits/kayron/pkg/release/schema/service"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"
)

type Config struct {
	Art cache.Interface[string, string]
	Log logger.Interface
	Ser cache.Interface[int, service.Service]
}

type Policy struct {
	art cache.Interface[string, string]
	log logger.Interface
	ser cache.Interface[int, service.Service]
}

func New(c Config) *Policy {
	if c.Art == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Art must not be empty", c)))
	}
	if c.Log == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Log must not be empty", c)))
	}
	if c.Ser == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Ser must not be empty", c)))
	}

	return &Policy{
		art: c.Art,
		log: c.Log,
		ser: c.Ser,
	}
}
