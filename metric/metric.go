// Heavily influenced by https://github.com/prometheus/client_golang
package metric

type Metric interface {
	// Desc returns the descriptor for the Metric. This method idempotently
	// returns the same descriptor throughout the lifetime of the
	// Metric. The returned descriptor is immutable by contract.
	Desc() *Desc
	// Write returns the the Gauge Metric value.
	// TODO: extend the return value to handle Metrics other than Gauges
	Write() (float64, error)
}
