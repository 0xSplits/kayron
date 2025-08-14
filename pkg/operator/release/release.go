// Package release implements the source business logic that caches all relevant
// service release settings for further use across the operator chain. See e.g.
// https://github.com/0xSplits/releases for a reference of the remote Github
// repository.
package release

import (
	"fmt"

	"github.com/0xSplits/kayron/pkg/cache"
	"github.com/0xSplits/kayron/pkg/envvar"
	"github.com/0xSplits/kayron/pkg/operator/release/canceler"
	"github.com/0xSplits/kayron/pkg/operator/release/resolver"
	"github.com/0xSplits/kayron/pkg/roghfs"
	"github.com/0xSplits/kayron/pkg/stack"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/google/go-github/v73/github"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"
)

type Config struct {
	Aws aws.Config
	Cac *cache.Cache
	Env envvar.Env
	Log logger.Interface
	Sta stack.Interface
}

type Release struct {
	cac *cache.Cache
	can canceler.Interface
	env envvar.Env
	git *github.Client
	log logger.Interface
	own string
	rep string
	res resolver.Interface
	sta stack.Interface
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
	if c.Log == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Log must not be empty", c)))
	}
	if c.Sta == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Sta must not be empty", c)))
	}

	var err error

	var git *github.Client
	{
		git = github.NewClient(nil).WithAuthToken(c.Env.GithubToken)
	}

	var own string
	var rep string
	{
		own, rep, err = roghfs.Parse(c.Env.ReleaseSource)
		if err != nil {
			tracer.Panic(tracer.Mask(err))
		}
	}

	var can canceler.Interface
	{
		can = canceler.New(canceler.Config{
			Aws: c.Aws,
			Env: c.Env,
			Sta: c.Sta,
		})
	}

	var res resolver.Interface
	{
		res = resolver.New(resolver.Config{
			Git: git,
			Own: own,
			Rep: rep,
		})
	}

	return &Release{
		cac: c.Cac,
		can: can,
		env: c.Env,
		git: git,
		log: c.Log,
		own: own,
		rep: rep,
		res: res,
		sta: c.Sta,
	}
}
