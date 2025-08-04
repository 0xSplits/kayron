package infrastructure

import (
	"fmt"

	"github.com/0xSplits/kayron/pkg/cache"
	"github.com/0xSplits/kayron/pkg/envvar"
	"github.com/0xSplits/kayron/pkg/release/schema/service"
	"github.com/0xSplits/kayron/pkg/roghfs"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/go-github/v73/github"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"
)

const (
	Bucket     = "splits-cf-templates"
	Directory  = "cloudformation"
	Repository = "infrastructure"
)

type Config struct {
	Art cache.Interface[string, string]
	Aws aws.Config
	Env envvar.Env
	Log logger.Interface
	Ser cache.Interface[int, service.Service]
}

type Infrastructure struct {
	art cache.Interface[string, string]
	as3 *s3.Client
	env string
	git *github.Client
	log logger.Interface
	own string
	ser cache.Interface[int, service.Service]
}

func New(c Config) *Infrastructure {
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

	return &Infrastructure{
		art: c.Art,
		as3: s3.NewFromConfig(c.Aws),
		env: c.Env.Environment,
		git: github.NewClient(nil).WithAuthToken(c.Env.GithubToken),
		log: c.Log,
		own: own,
		ser: c.Ser,
	}
}
