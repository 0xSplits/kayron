package operator

import (
	"fmt"

	"github.com/0xSplits/kayron/pkg/envvar"
	"github.com/0xSplits/kayron/pkg/recorder"
	"github.com/0xSplits/kayron/pkg/registry"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/google/go-github/v73/github"
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
	acf *cloudformation.Client
	env envvar.Env
	git *github.Client
	log logger.Interface
	reg registry.Interface
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

	cou := map[string]recorder.Interface{}

	gau := map[string]recorder.Interface{}

	{
		gau[Metric] = recorder.NewGauge(recorder.GaugeConfig{
			Des: "the timestamp of a deployment event for the released service",
			Lab: map[string][]string{
				"service": {"kayron", "server", "specta"},
			},
			Met: c.Met,
			Nam: Metric,
		})
	}

	his := map[string]recorder.Interface{}

	var reg registry.Interface
	{
		reg = registry.New(registry.Config{
			Env: c.Env,
			Log: c.Log,

			Cou: cou,
			Gau: gau,
			His: his,
		})
	}

	return &Handler{
		acf: cloudformation.NewFromConfig(c.Aws),
		env: c.Env,
		git: github.NewClient(nil).WithAuthToken(c.Env.GithubToken),
		log: c.Log,
		reg: reg,
	}
}
