// Package template fetches the current state of the currently deployed template
// version parameter. This operator function caches information about the
// currently deployed infrastructure versions.
package template

import (
	"fmt"

	"github.com/0xSplits/kayron/pkg/cache"
	"github.com/0xSplits/kayron/pkg/stack"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"
)

type Config struct {
	Cac *cache.Cache
	Log logger.Interface
	Sta stack.Interface
}

type Template struct {
	cac *cache.Cache
	log logger.Interface
	sta stack.Interface
}

func New(c Config) *Template {
	if c.Cac == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Cac must not be empty", c)))
	}
	if c.Log == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Log must not be empty", c)))
	}
	if c.Sta == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Sta must not be empty", c)))
	}

	return &Template{
		cac: c.Cac,
		log: c.Log,
		sta: c.Sta,
	}
}
