// Package registry verifies whether the desired service version is readily
// available inside the configured container registry. Only existing container
// images can actually be deployed.
package registry

import (
	"fmt"

	"github.com/0xSplits/kayron/pkg/cache"
	"github.com/0xSplits/kayron/pkg/envvar"
	"github.com/0xSplits/kayron/pkg/release/schema/service"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"
)

type Config struct {
	Art cache.Interface[string, string]
	Aws aws.Config
	Env envvar.Env
	Log logger.Interface
	Ser cache.Interface[int, service.Service]
}

type Registry struct {
	art cache.Interface[string, string]
	ecr *ecr.Client
	env envvar.Env
	log logger.Interface
	ser cache.Interface[int, service.Service]
}

func New(c Config) *Registry {
	if c.Art == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Art must not be empty", c)))
	}
	if c.Aws.Region == "" {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Aws must not be empty", c)))
	}
	if c.Log == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Log must not be empty", c)))
	}
	if c.Ser == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Ser must not be empty", c)))
	}

	return &Registry{
		art: c.Art,
		ecr: ecr.NewFromConfig(c.Aws),
		env: c.Env,
		log: c.Log,
		ser: c.Ser,
	}
}
