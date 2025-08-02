// Package release implements the source business logic that caches all relevant
// release settings for further use across the operator chain. See e.g.
// https://github.com/0xSplits/releases for a reference remote Github
// repository.
package release

import (
	"fmt"

	"github.com/0xSplits/kayron/pkg/cache"
	"github.com/0xSplits/kayron/pkg/envvar"
	"github.com/0xSplits/kayron/pkg/release/schema/service"
	"github.com/0xSplits/kayron/pkg/roghfs"
	"github.com/0xSplits/kayron/pkg/worker/handler/operator/release/resolver"
	"github.com/google/go-github/v73/github"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"
)

type Config struct {
	Art cache.Interface[string, string]
	Env envvar.Env
	Log logger.Interface
	Ser cache.Interface[int, service.Service]
}

type Release struct {
	art cache.Interface[string, string]
	env envvar.Env
	git *github.Client
	log logger.Interface
	own string
	rep string
	res resolver.Interface
	ser cache.Interface[int, service.Service]
}

func New(c Config) *Release {
	if c.Art == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Art must not be empty", c)))
	}
	if c.Log == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Log must not be empty", c)))
	}
	if c.Ser == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Ser must not be empty", c)))
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

	var res resolver.Interface
	{
		res = resolver.New(resolver.Config{
			Git: git,
			Own: own,
			Rep: rep,
		})
	}

	return &Release{
		art: c.Art,
		env: c.Env,
		git: git,
		log: c.Log,
		own: own,
		rep: rep,
		res: res,
		ser: c.Ser,
	}
}
