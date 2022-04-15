// Heavily influenced by https://github.com/prometheus/client_golang
package metric

import (
	"math"
	"sync/atomic"
)

type Gauge interface {
	Metric
	Collector

	// Set sets the Gauge to an arbitrary value.
	Set(float64)
}

func NewGauge(name string) Gauge {
	g := &gauge{desc: NewDesc(name, "", MetricTypeGauge)}
	// Init self-collection
	g.init(g)
	return g
}

type gauge struct {
	valBits uint64

	selfCollector

	desc *Desc
}

func (g *gauge) Desc() *Desc {
	return g.desc
}

func (g *gauge) Write() (float64, error) {
	return math.Float64frombits(atomic.LoadUint64(&g.valBits)), nil
}

func (g *gauge) Set(val float64) {
	atomic.StoreUint64(&g.valBits, math.Float64bits(val))
}
