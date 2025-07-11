package registry

// Interface provides a mechanism for the server and worker handlers to record
// the metrics they are responsible for. This abstract metrics registry defines
// the whitelist of available metrics. The consumers, here Kayron's server and
// worker handlers, use this registry to manage their desired metrics on demand.
// It is important for the registry to only provide predefined metrics for
// timeseries that we want to track deliberately. E.g. it should not be possible
// for the outside world to create metrics or labels arbitrarily based on user
// input.
type Interface interface {
	// Counter performs a metrics update for the counter metric matching the
	// provided metric name, if the given input is considered valid.
	//
	//     inp[0] the name of the counter metric to fetch
	//     inp[1] the value of the counter metric to record
	//     inp[2] the labels of the counter metric to associate
	//
	//     out[0] the error describing why the operation failed, if any
	//
	Counter(string, float64, map[string]string) error

	// Gauge performs a metrics update for the gauge metric matching the provided
	// metric name, if the given input is considered valid.
	//
	//     inp[0] the name of the gauge metric to fetch
	//     inp[1] the value of the gauge metric to record
	//     inp[2] the labels of the gauge metric to associate
	//
	//     out[0] the error describing why the operation failed, if any
	//
	Gauge(string, float64, map[string]string) error

	// Histogram performs a metrics update for the histogram metric matching the
	// provided metric name, if the given input is considered valid.
	//
	//     inp[0] the name of the histogram metric to fetch
	//     inp[1] the value of the histogram metric to record
	//     inp[2] the labels of the histogram metric to associate
	//
	//     out[0] the error describing why the operation failed, if any
	//
	Histogram(string, float64, map[string]string) error
}
