// Heavily influenced by https://github.com/prometheus/client_golang
package metric

import "fmt"

type MetricType int32

const (
	MetricTypeGauge MetricType = iota
)

var MetricTypeName = map[MetricType]string{
	0: "GAUGE",
}

func (x MetricType) String() string {
	return MetricTypeName[x]
}

type Desc struct {
	// t is the metric type
	t MetricType
	// name is the name of the descriptor - must be globally unique
	name string
	// err is an error that occurred during construction. It is reported on
	// registration time.
	err error
}

func NewDesc(name, help string, t MetricType) *Desc {
	return &Desc{
		t:    t,
		name: name,
	}
}

func (d *Desc) String() string {
	return fmt.Sprintf("Desc{name: %q, type: %q}", d.name, d.t)
}
