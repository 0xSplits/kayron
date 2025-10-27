// Package status manages deployment status updates so that humans can better
// understand the current state of the underlying systems. E.g. this status
// operator function may create and update pull request comments for preview
// deployments in Github.
package status

import (
	"fmt"

	"github.com/0xSplits/kayron/pkg/cache"
	"github.com/0xSplits/kayron/pkg/envvar"
	"github.com/0xSplits/kayron/pkg/policy"
	"github.com/0xSplits/roghfs"
	"github.com/google/go-github/v76/github"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"
)

type Config struct {
	Cac *cache.Cache
	Env envvar.Env
	Git *github.Client
	Log logger.Interface
	Pol *policy.Policy
}

type Status struct {
	cac *cache.Cache
	env envvar.Env
	git *github.Client
	log logger.Interface
	own string
	pol *policy.Policy
}

func New(c Config) *Status {
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
	{
		own, _, err = roghfs.Parse(c.Env.ReleaseSource)
		if err != nil {
			tracer.Panic(tracer.Mask(err))
		}
	}

	return &Status{
		cac: c.Cac,
		env: c.Env,
		git: c.Git,
		log: c.Log,
		own: own,
		pol: c.Pol,
	}
}
