// Package template fetches the current state of the currently deployed template
// version parameter. This operator function caches information about the
// currently deployed infrastructure versions.
package template

import (
	"fmt"

	"github.com/0xSplits/kayron/pkg/cache"
	"github.com/0xSplits/kayron/pkg/envvar"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
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

type Template struct {
	cfc *cloudformation.Client
	cac *cache.Cache
	env envvar.Env
	log logger.Interface
	tag *resourcegroupstaggingapi.Client
}

func New(c Config) *Template {
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

	return &Template{
		cfc: cloudformation.NewFromConfig(c.Aws),
		cac: c.Cac,
		env: c.Env,
		log: c.Log,
		tag: resourcegroupstaggingapi.NewFromConfig(c.Aws),
	}
}
