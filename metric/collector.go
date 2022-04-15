// Heavily influenced by https://github.com/prometheus/client_golang
package metric

type Collector interface {
	Describe(chan<- *Desc)
	Collect(chan<- Metric)
}

// selfCollector implements Collector for a single Metric so that the Metric
// collects itself. Add it as an anonymous field to a struct that implements
// Metric, and call init with the Metric itself as an argument.
type selfCollector struct {
	self Metric
}

// init provides the selfCollector with a reference to the metric it is supposed
// to collect.
func (c *selfCollector) init(self Metric) {
	c.self = self
}

// Describe implements Collector.
func (c *selfCollector) Describe(ch chan<- *Desc) {
	ch <- c.self.Desc()
}

// Collect implements Collector.
func (c *selfCollector) Collect(ch chan<- Metric) {
	ch <- c.self
}
