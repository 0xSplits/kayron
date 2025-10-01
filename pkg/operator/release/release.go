// Package release implements the source business logic that caches all relevant
// service release settings for further use across the operator chain. See e.g.
// https://github.com/0xSplits/releases for a reference of the remote Github
// repository.
package release

import (
	"fmt"

	"github.com/0xSplits/kayron/pkg/cache"
	"github.com/0xSplits/kayron/pkg/envvar"
	"github.com/0xSplits/kayron/pkg/operator/release/resolver"
	"github.com/0xSplits/kayron/pkg/policy"
	"github.com/0xSplits/roghfs"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/google/go-github/v75/github"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"
)

type Config struct {
	Aws aws.Config
	Cac *cache.Cache
	Env envvar.Env
	Git *github.Client
	Log logger.Interface
	Pol *policy.Policy
}

type Release struct {
	cac *cache.Cache
	env envvar.Env
	git *github.Client
	log logger.Interface
	own string
	rep string
	res resolver.Interface
	pol *policy.Policy
}

func New(c Config) *Release {
	if c.Aws.Region == "" {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Aws must not be empty", c)))
	}
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
	if c.Pol == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Pol must not be empty", c)))
	}

	var err error

	var own string
	var rep string
	{
		own, rep, err = roghfs.Parse(c.Env.ReleaseSource)
		if err != nil {
			tracer.Panic(tracer.Mask(err))
		}
	}

	var res resolver.Interface
	{
		res = resolver.New(resolver.Config{
			Git: c.Git,
			Own: own,
			Rep: rep,
		})
	}

	return &Release{
		cac: c.Cac,
		env: c.Env,
		git: c.Git,
		log: c.Log,
		own: own,
		rep: rep,
		res: res,
		pol: c.Pol,
	}
}
