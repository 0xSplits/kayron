package operator

import (
	"fmt"

	"github.com/0xSplits/kayron/pkg/cache"
	"github.com/0xSplits/kayron/pkg/envvar"
	"github.com/0xSplits/kayron/pkg/release/schema/service"
	"github.com/0xSplits/kayron/pkg/worker/handler/operator/cloudformation"
	"github.com/0xSplits/kayron/pkg/worker/handler/operator/container"
	"github.com/0xSplits/kayron/pkg/worker/handler/operator/infrastructure"
	"github.com/0xSplits/kayron/pkg/worker/handler/operator/policy"
	"github.com/0xSplits/kayron/pkg/worker/handler/operator/reference"
	"github.com/0xSplits/kayron/pkg/worker/handler/operator/registry"
	"github.com/0xSplits/kayron/pkg/worker/handler/operator/release"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"
	"go.opentelemetry.io/otel/metric"
)

type Config struct {
	Aws aws.Config
	Env envvar.Env
	Log logger.Interface
	Met metric.Meter
}

type Handler struct {
	clo *cloudformation.CloudFormation
	con *container.Container
	inf *infrastructure.Infrastructure
	log logger.Interface
	pol *policy.Policy
	ref *reference.Reference
	reg *registry.Registry
	rel *release.Release
}

func New(c Config) *Handler {
	if c.Aws.Region == "" {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Aws must not be empty", c)))
	}
	if c.Log == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Log must not be empty", c)))
	}
	if c.Met == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Met must not be empty", c)))
	}

	// The cache implementations used here are the dumb pipes connecting our
	// business logic across different boundary conditions by providing critical
	// information required to make the entire operator chain work as expected.
	// The operator's worker handler is executed iteratively all the time, but the
	// real work may only be done if the cached information provided deeper down
	// the stack is sufficient. E.g. an empty cache result may cause some business
	// logic to be skipped temporarily. Note that the cache keys setup here are
	// cross references to the respective cache values within the various cache
	// implementations.

	var art cache.Interface[string, string]
	{
		art = cache.New[string, string]()
	}

	var ser cache.Interface[int, service.Service]
	{
		ser = cache.New[int, service.Service]()
	}

	var clo *cloudformation.CloudFormation
	var con *container.Container
	var inf *infrastructure.Infrastructure
	var pol *policy.Policy
	var ref *reference.Reference
	var rel *release.Release
	var reg *registry.Registry
	{
		clo = cloudformation.New(cloudformation.Config{Art: art, Aws: c.Aws, Env: c.Env, Log: c.Log, Met: c.Met, Ser: ser})
		con = container.New(container.Config{Art: art, Aws: c.Aws, Env: c.Env, Log: c.Log, Ser: ser})
		inf = infrastructure.New(infrastructure.Config{Art: art, Aws: c.Aws, Env: c.Env, Log: c.Log, Ser: ser})
		pol = policy.New(policy.Config{Art: art, Log: c.Log, Ser: ser})
		ref = reference.New(reference.Config{Art: art, Env: c.Env, Log: c.Log, Ser: ser})
		rel = release.New(release.Config{Art: art, Env: c.Env, Log: c.Log, Ser: ser})
		reg = registry.New(registry.Config{Art: art, Aws: c.Aws, Env: c.Env, Log: c.Log, Ser: ser})
	}

	return &Handler{
		clo: clo,
		con: con,
		inf: inf,
		log: c.Log,
		pol: pol,
		ref: ref,
		reg: reg,
		rel: rel,
	}
}
