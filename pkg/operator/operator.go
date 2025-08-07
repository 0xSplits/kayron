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
	Env envvar.Env
	Log logger.Interface
	Met metric.Meter
}

type Operator struct {
	CloudFormation *cloudformation.CloudFormation
	Container      *container.Container
	Infrastructure *infrastructure.Infrastructure
	Policy         *policy.Policy
	Reference      *reference.Reference
	Release        *release.Release
	Registry       *registry.Registry
	Template       *template.Template
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
		CloudFormation: cloudformation.New(cloudformation.Config{Aws: c.Aws, Ctx: c.Ctx, Env: c.Env, Log: c.Log, Met: c.Met}),
		Container:      container.New(container.Config{Aws: c.Aws, Ctx: c.Ctx, Env: c.Env, Log: c.Log}),
		Infrastructure: infrastructure.New(infrastructure.Config{Aws: c.Aws, Ctx: c.Ctx, Env: c.Env, Log: c.Log}),
		Policy:         policy.New(policy.Config{Ctx: c.Ctx, Log: c.Log}),
		Reference:      reference.New(reference.Config{Ctx: c.Ctx, Env: c.Env, Log: c.Log}),
		Release:        release.New(release.Config{Ctx: c.Ctx, Env: c.Env, Log: c.Log}),
		Registry:       registry.New(registry.Config{Aws: c.Aws, Ctx: c.Ctx, Env: c.Env, Log: c.Log}),
		Template:       template.New(template.Config{Aws: c.Aws, Ctx: c.Ctx, Env: c.Env, Log: c.Log}),
	}
}
