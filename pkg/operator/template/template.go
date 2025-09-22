// Package template fetches the current state of the currently deployed template
// version parameter. This operator function caches information about the
// currently deployed infrastructure versions.
package template

import (
	"fmt"

	"github.com/0xSplits/kayron/pkg/cache"
	"github.com/0xSplits/kayron/pkg/policy"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"
)

type Config struct {
	Cac *cache.Cache
	Log logger.Interface
	Pol *policy.Policy
}

type Template struct {
	cac *cache.Cache
	log logger.Interface
	pol *policy.Policy
}

func New(c Config) *Template {
	if c.Cac == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Cac must not be empty", c)))
	}
	if c.Log == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Log must not be empty", c)))
	}
	if c.Pol == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Pol must not be empty", c)))
	}

	return &Template{
		cac: c.Cac,
		log: c.Log,
		pol: c.Pol,
	}
}
