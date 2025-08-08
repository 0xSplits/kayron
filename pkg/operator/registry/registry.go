// Package registry verifies whether the desired service version is readily
// available inside the configured container registry. Only existing container
// images can actually be deployed.
package registry

import (
	"fmt"

	"github.com/0xSplits/kayron/pkg/cache"
	"github.com/0xSplits/kayron/pkg/envvar"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"
)

type Config struct {
	Aws aws.Config
	Cac *cache.Cache
	Env envvar.Env
	Log logger.Interface
}

type Registry struct {
	cac *cache.Cache
	ecr *ecr.Client
	env envvar.Env
	log logger.Interface
}

func New(c Config) *Registry {
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

	return &Registry{
		cac: c.Cac,
		ecr: ecr.NewFromConfig(c.Aws),
		env: c.Env,
		log: c.Log,
	}
}
