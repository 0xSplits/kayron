package operator

import (
	"fmt"

	"github.com/0xSplits/kayron/pkg/cache"
	"github.com/0xSplits/kayron/pkg/envvar"
	"github.com/0xSplits/kayron/pkg/release/schema/service"
	"github.com/0xSplits/kayron/pkg/worker/handler/operator/container"
	"github.com/0xSplits/kayron/pkg/worker/handler/operator/reference"
	"github.com/0xSplits/kayron/pkg/worker/handler/operator/release"
	"github.com/0xSplits/otelgo/registry"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"
	"go.opentelemetry.io/otel/metric"
)

const (
	Metric = "deployment_event"
)

type Config struct {
	Aws aws.Config
	Env envvar.Env
	Log logger.Interface
	Met metric.Meter
}

type Handler struct {
	con *container.Container
	log logger.Interface
	ref *reference.Reference
	reg registry.Interface
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

	var reg registry.Interface
	{
		reg = newRegistry(c.Env.Environment, c.Log, c.Met)
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

	var con *container.Container
	var ref *reference.Reference
	var rel *release.Release
	{
		con = container.New(container.Config{Aws: c.Aws, Art: art, Env: c.Env, Log: c.Log, Ser: ser})
		ref = reference.New(reference.Config{Art: art, Env: c.Env, Log: c.Log, Ser: ser})
		rel = release.New(release.Config{Art: art, Env: c.Env, Log: c.Log, Ser: ser})
	}

	return &Handler{
		con: con,
		log: c.Log,
		ref: ref,
		reg: reg,
		rel: rel,
	}
}
