package operator

import (
	"fmt"

	"github.com/0xSplits/kayron/pkg/cache"
	"github.com/0xSplits/kayron/pkg/envvar"
	"github.com/0xSplits/kayron/pkg/operator/cloudformation"
	"github.com/0xSplits/kayron/pkg/operator/container"
	"github.com/0xSplits/kayron/pkg/operator/infrastructure"
	"github.com/0xSplits/kayron/pkg/operator/preview"
	"github.com/0xSplits/kayron/pkg/operator/reference"
	"github.com/0xSplits/kayron/pkg/operator/registry"
	"github.com/0xSplits/kayron/pkg/operator/release"
	"github.com/0xSplits/kayron/pkg/operator/status"
	"github.com/0xSplits/kayron/pkg/operator/template"
	"github.com/0xSplits/kayron/pkg/policy"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"
	"go.opentelemetry.io/otel/metric"
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

type Operator struct {
	cloudFormation *cloudformation.CloudFormation
	container      *container.Container
	infrastructure *infrastructure.Infrastructure
	preview        *preview.Preview
	reference      *reference.Reference
	release        *release.Release
	registry       *registry.Registry
	status         *status.Status
	template       *template.Template
}

func New(c Config) *Operator {
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

	return &Operator{
		cloudFormation: cloudformation.New(cloudformation.Config{Aws: c.Aws, Cac: c.Cac, Dry: c.Dry, Env: c.Env, Log: c.Log, Met: c.Met, Pol: c.Pol}),
		container:      container.New(container.Config{Aws: c.Aws, Cac: c.Cac, Env: c.Env, Log: c.Log}),
		infrastructure: infrastructure.New(infrastructure.Config{Aws: c.Aws, Cac: c.Cac, Dry: c.Dry, Env: c.Env, Log: c.Log, Pol: c.Pol}),
		preview:        preview.New(preview.Config{Cac: c.Cac, Env: c.Env, Log: c.Log}),
		reference:      reference.New(reference.Config{Cac: c.Cac, Env: c.Env, Log: c.Log}),
		release:        release.New(release.Config{Aws: c.Aws, Cac: c.Cac, Env: c.Env, Log: c.Log, Pol: c.Pol}),
		registry:       registry.New(registry.Config{Aws: c.Aws, Cac: c.Cac, Env: c.Env, Log: c.Log}),
		status:         status.New(status.Config{Env: c.Env, Log: c.Log, Pol: c.Pol}),
		template:       template.New(template.Config{Cac: c.Cac, Log: c.Log, Pol: c.Pol}),
	}
}
