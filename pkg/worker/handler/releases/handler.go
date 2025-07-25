package releases

import (
	"fmt"

	"github.com/0xSplits/kayron/pkg/cache"
	"github.com/0xSplits/kayron/pkg/envvar"
	"github.com/0xSplits/kayron/pkg/release/schema/service"
	"github.com/0xSplits/kayron/pkg/roghfs"
	"github.com/0xSplits/kayron/pkg/worker/handler/releases/resolver"
	"github.com/google/go-github/v73/github"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"
)

type Config struct {
	Env envvar.Env
	Log logger.Interface
	Rel cache.Interface[int, service.Service]
}

type Handler struct {
	env envvar.Env
	git *github.Client
	log logger.Interface
	own string
	rep string
	rel cache.Interface[int, service.Service]
	res resolver.Interface
}

func New(c Config) *Handler {
	if c.Log == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Log must not be empty", c)))
	}
	if c.Rel == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Rel must not be empty", c)))
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

	return &Handler{
		env: c.Env,
		git: git,
		log: c.Log,
		own: own,
		rep: rep,
		rel: c.Rel,
		res: res,
	}
}
