package recorder

import (
	"context"

	"github.com/xh3b4sd/tracer"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type HistogramConfig struct {
	Buc []float64
	Des string
	Lab map[string][]string
	Met metric.Meter
	Nam string
}

type Histogram struct {
	his metric.Float64Histogram
	lab map[string][]string
}

func NewHistogram(c HistogramConfig) *Histogram {
	his, err := c.Met.Float64Histogram(c.Nam, metric.WithDescription(c.Des), metric.WithExplicitBucketBoundaries(c.Buc...))
	if err != nil {
		tracer.Panic(tracer.Mask(err))
	}

	return &Histogram{
		his: his,
		lab: c.Lab,
	}
}

func (h *Histogram) Labels() map[string][]string {
	return h.lab
}

func (h *Histogram) Record(val float64, lab ...attribute.KeyValue) {
	h.his.Record(context.Background(), val, metric.WithAttributes(lab...))
}
