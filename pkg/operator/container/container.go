// Package container fetches the current state of the currently deployed
// container image tags. This operator function caches information about the
// currently deployed service versions.
package container

import (
	"fmt"

	"github.com/0xSplits/kayron/pkg/cache"
	"github.com/0xSplits/kayron/pkg/envvar"
	"github.com/0xSplits/kayron/pkg/release/schema/service"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
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

type Container struct {
	art cache.Interface[string, string]
	ecs *ecs.Client
	env envvar.Env
	log logger.Interface
	ser cache.Interface[int, service.Service]
	tag *resourcegroupstaggingapi.Client
}

func New(c Config) *Container {
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

	return &Container{
		art: c.Art,
		ecs: ecs.NewFromConfig(c.Aws),
		env: c.Env,
		log: c.Log,
		ser: c.Ser,
		tag: resourcegroupstaggingapi.NewFromConfig(c.Aws),
	}
}
