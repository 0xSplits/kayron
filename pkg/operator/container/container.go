// Package container fetches the current state of the currently deployed
// container image tags. This operator function caches information about the
// currently deployed service versions.
package container

import (
	"fmt"

	"github.com/0xSplits/kayron/pkg/cache"
	"github.com/0xSplits/kayron/pkg/envvar"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"
)

type Config struct {
	Aws aws.Config
	Cac *cache.Cache
	Env envvar.Env
	Log logger.Interface
}

type Container struct {
	cac *cache.Cache
	ecs *ecs.Client
	env envvar.Env
	log logger.Interface
	tag *resourcegroupstaggingapi.Client
}

func New(c Config) *Container {
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

	return &Container{
		cac: c.Cac,
		ecs: ecs.NewFromConfig(c.Aws),
		env: c.Env,
		log: c.Log,
		tag: resourcegroupstaggingapi.NewFromConfig(c.Aws),
	}
}
