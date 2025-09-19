package preview

import (
	"fmt"

	"github.com/0xSplits/kayron/pkg/envvar"
	"github.com/0xSplits/kayron/pkg/scanner"
	"github.com/0xSplits/roghfs"
	"github.com/google/go-github/v73/github"
	"github.com/xh3b4sd/tracer"
)

type Config struct {
	Env envvar.Env
	Inp []byte
}

type Preview struct {
	git *github.Client
	inp []byte
	own string
	sca *scanner.Scanner
}

func New(c Config) *Preview {
	if c.Env.Environment == "" {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Env must not be empty", c)))
	}
	if c.Inp == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Inp must not be empty", c)))
	}

	var err error

	var git *github.Client
	{
		git = github.NewClient(nil).WithAuthToken(c.Env.GithubToken)
	}

	var own string
	{
		own, _, err = roghfs.Parse(c.Env.ReleaseSource)
		if err != nil {
			tracer.Panic(tracer.Mask(err))
		}
	}

	var sca *scanner.Scanner
	{
		sca = scanner.New(scanner.Config{
			Inp: c.Inp,
		})
	}

	return &Preview{
		git: git,
		inp: c.Inp,
		own: own,
		sca: sca,
	}
}
