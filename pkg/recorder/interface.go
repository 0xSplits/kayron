package recorder

import "go.opentelemetry.io/otel/attribute"

type Interface interface {
	// Labels returns the whitelisted set of labels that is allowed to be recorded
	// for the underlying metric.
	Labels() map[string][]string

	// Record tracks the given value for the underlying metric.
	Record(val float64, lab ...attribute.KeyValue)
}
