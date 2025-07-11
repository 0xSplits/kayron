package recorder

import (
	"context"

	"github.com/xh3b4sd/tracer"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type CounterConfig struct {
	Des string
	Lab map[string][]string
	Met metric.Meter
	Nam string
}

type Counter struct {
	cou metric.Float64Counter
	lab map[string][]string
}

func NewCounter(c CounterConfig) *Counter {
	cou, err := c.Met.Float64Counter(c.Nam, metric.WithDescription(c.Des))
	if err != nil {
		tracer.Panic(tracer.Mask(err))
	}

	return &Counter{
		lab: c.Lab,
		cou: cou,
	}
}

func (c *Counter) Labels() map[string][]string {
	return c.lab
}

func (c *Counter) Record(val float64, lab ...attribute.KeyValue) {
	c.cou.Add(context.Background(), val, metric.WithAttributes(lab...))
}
