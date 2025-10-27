// Package preview injects preview deployments into our internal release
// artifact cache. This operator function enables us to expose additional
// development services within the testing environment only.
package preview

import (
	"fmt"

	"github.com/0xSplits/kayron/pkg/cache"
	"github.com/0xSplits/kayron/pkg/envvar"
	"github.com/0xSplits/kayron/pkg/preview"
	"github.com/google/go-github/v76/github"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"
)

type Config struct {
	Cac *cache.Cache
	Env envvar.Env
	Git *github.Client
	Log logger.Interface
}

type Preview struct {
	cac *cache.Cache
	env envvar.Env
	log logger.Interface
	pre *preview.Preview
}

func New(c Config) *Preview {
	if c.Cac == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Cac must not be empty", c)))
	}
	if c.Env.Environment == "" {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Env must not be empty", c)))
	}
	if c.Git == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Git must not be empty", c)))
	}
	if c.Log == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Log must not be empty", c)))
	}

	var pre *preview.Preview
	{
		pre = preview.New(preview.Config{
			Env: c.Env,
			Git: c.Git,
			Inp: []byte{},
		})
	}

	return &Preview{
		cac: c.Cac,
		env: c.Env,
		log: c.Log,
		pre: pre,
	}
}
