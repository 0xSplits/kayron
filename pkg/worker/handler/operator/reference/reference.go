package reference

import (
	"fmt"

	"github.com/0xSplits/kayron/pkg/cache"
	"github.com/0xSplits/kayron/pkg/envvar"
	"github.com/0xSplits/kayron/pkg/release/artifact"
	"github.com/0xSplits/kayron/pkg/release/schema/service"
	"github.com/0xSplits/kayron/pkg/roghfs"
	"github.com/google/go-github/v73/github"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"
)

type Config struct {
	Art cache.Interface[int, artifact.Artifact]
	Env envvar.Env
	Log logger.Interface
	Ser cache.Interface[int, service.Service]
}

type Reference struct {
	art cache.Interface[int, artifact.Artifact]
	git *github.Client
	log logger.Interface
	own string
	ser cache.Interface[int, service.Service]
}

func New(c Config) *Reference {
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

	var own string
	{
		own, _, err = roghfs.Parse(c.Env.ReleaseSource)
		if err != nil {
			tracer.Panic(tracer.Mask(err))
		}
	}

	return &Reference{
		art: c.Art,
		git: github.NewClient(nil).WithAuthToken(c.Env.GithubToken),
		log: c.Log,
		own: own,
		ser: c.Ser,
	}
}
