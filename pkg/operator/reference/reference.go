// Package reference fetches the desired state of any service release source
// code as Git reference. Those consistent release tags and commit shas define
// the service versions intended to be deployed.
package reference

import (
	"fmt"

	"github.com/0xSplits/kayron/pkg/cache"
	"github.com/0xSplits/kayron/pkg/envvar"
	"github.com/0xSplits/kayron/pkg/webhook"
	"github.com/0xSplits/roghfs"
	"github.com/google/go-github/v75/github"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"
)

type Config struct {
	Cac *cache.Cache
	Env envvar.Env
	Log logger.Interface
	Whk *webhook.Webhook
}

type Reference struct {
	cac *cache.Cache
	env envvar.Env
	git *github.Client
	log logger.Interface
	own string
	whk *webhook.Webhook
}

func New(c Config) *Reference {
	if c.Cac == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Cac must not be empty", c)))
	}
	if c.Env.Environment == "" {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Env must not be empty", c)))
	}
	if c.Log == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Log must not be empty", c)))
	}
	if c.Whk == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Whk must not be empty", c)))
	}

	var err error

	var own string
	{
		own, _, err = roghfs.Parse(c.Env.ReleaseSource)
		if err != nil {
			tracer.Panic(tracer.Mask(err))
		}
	}

	return &Reference{
		cac: c.Cac,
		env: c.Env,
		git: github.NewClient(nil).WithAuthToken(c.Env.GithubToken),
		log: c.Log,
		own: own,
		whk: c.Whk,
	}
}
