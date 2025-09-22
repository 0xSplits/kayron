// Package cloudformation triggers CloudFormation stack updates once a valid
// state drift was detected. This operator function effectively applies any
// infrastructure and service changes according to the configured desired state.
package cloudformation

import (
	"fmt"

	"github.com/0xSplits/kayron/pkg/cache"
	"github.com/0xSplits/kayron/pkg/envvar"
	"github.com/0xSplits/kayron/pkg/policy"
	"github.com/0xSplits/otelgo/registry"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"
	"go.opentelemetry.io/otel/metric"
)

const (
	Metric = "deployment_event"
)

type Config struct {
	Aws aws.Config
	Cac *cache.Cache
	Dry bool
	Env envvar.Env
	Log logger.Interface
	Met metric.Meter
	Pol *policy.Policy
}

type CloudFormation struct {
	cfc *cloudformation.Client
	cac *cache.Cache
	dry bool
	env envvar.Env
	log logger.Interface
	pol *policy.Policy
	reg registry.Interface
}

func New(c Config) *CloudFormation {
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
	if c.Met == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Met must not be empty", c)))
	}
	if c.Pol == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Pol must not be empty", c)))
	}

	var reg registry.Interface
	{
		reg = newRegistry(c.Env.Environment, c.Log, c.Met)
	}

	return &CloudFormation{
		cfc: cloudformation.NewFromConfig(c.Aws),
		cac: c.Cac,
		dry: c.Dry,
		env: c.Env,
		log: c.Log,
		pol: c.Pol,
		reg: reg,
	}
}
