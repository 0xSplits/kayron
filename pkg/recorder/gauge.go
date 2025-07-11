package recorder

import (
	"context"

	"github.com/xh3b4sd/tracer"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type GaugeConfig struct {
	Des string
	Lab map[string][]string
	Met metric.Meter
	Nam string
}

type Gauge struct {
	gau metric.Float64Gauge
	lab map[string][]string
}

func NewGauge(c GaugeConfig) *Gauge {
	gau, err := c.Met.Float64Gauge(c.Nam, metric.WithDescription(c.Des))
	if err != nil {
		tracer.Panic(tracer.Mask(err))
	}

	return &Gauge{
		lab: c.Lab,
		gau: gau,
	}
}

func (g *Gauge) Labels() map[string][]string {
	return g.lab
}

func (g *Gauge) Record(val float64, lab ...attribute.KeyValue) {
	g.gau.Record(context.Background(), val, metric.WithAttributes(lab...))
}
