package cloudformation

import (
	"github.com/0xSplits/otelgo/recorder"
	"github.com/0xSplits/otelgo/registry"
	"github.com/xh3b4sd/logger"
	"go.opentelemetry.io/otel/metric"
)

func newRegistry(env string, log logger.Interface, met metric.Meter) registry.Interface {
	cou := map[string]recorder.Interface{}

	gau := map[string]recorder.Interface{}

	{
		gau[Metric] = recorder.NewGauge(recorder.GaugeConfig{
			Des: "the timestamp of a deployment event for the released service",
			Lab: map[string][]string{
				"service": {"kayron", "server", "specta"},
			},
			Met: met,
			Nam: Metric,
		})
	}

	his := map[string]recorder.Interface{}

	return registry.New(registry.Config{
		Env: env,
		Log: log,

		Cou: cou,
		Gau: gau,
		His: his,
	})
}
