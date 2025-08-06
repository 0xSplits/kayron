// Package cloudformation triggers CloudFormation stack updates once a valid
// state drift was detected. This operator function effectively applies any
// infrastructure and service changes according to the configured desired state.
package cloudformation

import (
	"fmt"

	"github.com/0xSplits/kayron/pkg/cache"
	"github.com/0xSplits/kayron/pkg/envvar"
	"github.com/0xSplits/kayron/pkg/release/schema/service"
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
	Art cache.Interface[string, string]
	Aws aws.Config
	Env envvar.Env
	Log logger.Interface
	Met metric.Meter
	Ser cache.Interface[int, service.Service]
}

type CloudFormation struct {
	art cache.Interface[string, string]
	cfc *cloudformation.Client
	log logger.Interface
	reg registry.Interface
	ser cache.Interface[int, service.Service]
}

func New(c Config) *CloudFormation {
	if c.Art == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Art must not be empty", c)))
	}
	if c.Aws.Region == "" {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Aws must not be empty", c)))
	}
	if c.Log == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Log must not be empty", c)))
	}
	if c.Met == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Met must not be empty", c)))
	}
	if c.Ser == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Ser must not be empty", c)))
	}

	var reg registry.Interface
	{
		reg = newRegistry(c.Env.Environment, c.Log, c.Met)
	}

	return &CloudFormation{
		art: c.Art,
		cfc: cloudformation.NewFromConfig(c.Aws),
		log: c.Log,
		reg: reg,
		ser: c.Ser,
	}
}
