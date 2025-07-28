package resolver

import (
	"fmt"

	"github.com/google/go-github/v73/github"
	"github.com/xh3b4sd/tracer"
)

type Config struct {
	Git *github.Client
	Own string
	Rep string
}

type Resolver struct {
	git *github.Client
	own string
	rep string
}

func New(c Config) *Resolver {
	if c.Git == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Git must not be empty", c)))
	}
	if c.Own == "" {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Own must not be empty", c)))
	}
	if c.Rep == "" {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Rep must not be empty", c)))
	}

	return &Resolver{
		git: c.Git,
		own: c.Own,
		rep: c.Rep,
	}
}
