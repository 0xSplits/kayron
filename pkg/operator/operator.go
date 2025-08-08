package operator

import (
	"fmt"

	"github.com/0xSplits/kayron/pkg/context"
	"github.com/0xSplits/kayron/pkg/envvar"
	"github.com/0xSplits/kayron/pkg/operator/cloudformation"
	"github.com/0xSplits/kayron/pkg/operator/container"
	"github.com/0xSplits/kayron/pkg/operator/infrastructure"
	"github.com/0xSplits/kayron/pkg/operator/policy"
	"github.com/0xSplits/kayron/pkg/operator/reference"
	"github.com/0xSplits/kayron/pkg/operator/registry"
	"github.com/0xSplits/kayron/pkg/operator/release"
	"github.com/0xSplits/kayron/pkg/operator/template"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"
	"go.opentelemetry.io/otel/metric"
)

type Config struct {
	Aws aws.Config
	Ctx *context.Context
	Dry bool
	Env envvar.Env
	Log logger.Interface
	Met metric.Meter
}

type Operator struct {
	cloudFormation *cloudformation.CloudFormation
	container      *container.Container
	infrastructure *infrastructure.Infrastructure
	policy         *policy.Policy
	reference      *reference.Reference
	release        *release.Release
	registry       *registry.Registry
	template       *template.Template
}

func New(c Config) *Operator {
	if c.Aws.Region == "" {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Aws must not be empty", c)))
	}
	if c.Ctx == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Ctx must not be empty", c)))
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

	return &Operator{
		cloudFormation: cloudformation.New(cloudformation.Config{Aws: c.Aws, Ctx: c.Ctx, Dry: c.Dry, Env: c.Env, Log: c.Log, Met: c.Met}),
		container:      container.New(container.Config{Aws: c.Aws, Ctx: c.Ctx, Env: c.Env, Log: c.Log}),
		infrastructure: infrastructure.New(infrastructure.Config{Aws: c.Aws, Ctx: c.Ctx, Dry: c.Dry, Env: c.Env, Log: c.Log}),
		policy:         policy.New(policy.Config{Ctx: c.Ctx, Log: c.Log}),
		reference:      reference.New(reference.Config{Ctx: c.Ctx, Env: c.Env, Log: c.Log}),
		release:        release.New(release.Config{Ctx: c.Ctx, Env: c.Env, Log: c.Log}),
		registry:       registry.New(registry.Config{Aws: c.Aws, Ctx: c.Ctx, Env: c.Env, Log: c.Log}),
		template:       template.New(template.Config{Aws: c.Aws, Ctx: c.Ctx, Env: c.Env, Log: c.Log}),
	}
}
