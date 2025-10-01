// Package infrastructure prepares the desired state of the CloudFormation
// templates provided within the configured infrastructure repository. Those
// templates are fetched from Github and uploaded to S3. See e.g.
// https://github.com/0xSplits/infrastructure for a reference of the remote
// Github repository.
package infrastructure

import (
	"fmt"

	"github.com/0xSplits/kayron/pkg/cache"
	"github.com/0xSplits/kayron/pkg/envvar"
	"github.com/0xSplits/kayron/pkg/policy"
	"github.com/0xSplits/roghfs"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/go-github/v75/github"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"
)

type Config struct {
	Aws aws.Config
	Cac *cache.Cache
	Dry bool
	Env envvar.Env
	Git *github.Client
	Log logger.Interface
	Pol *policy.Policy
}

type Infrastructure struct {
	as3 *s3.Client
	cac *cache.Cache
	dry bool
	env envvar.Env
	git *github.Client
	log logger.Interface
	own string
	pol *policy.Policy
}

func New(c Config) *Infrastructure {
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
	{
		own, _, err = roghfs.Parse(c.Env.ReleaseSource)
		if err != nil {
			tracer.Panic(tracer.Mask(err))
		}
	}

	return &Infrastructure{
		as3: s3.NewFromConfig(c.Aws),
		cac: c.Cac,
		dry: c.Dry,
		env: c.Env,
		git: c.Git,
		log: c.Log,
		own: own,
		pol: c.Pol,
	}
}
