// Heavily influenced by https://github.com/prometheus/client_golang
package metric

import (
	"fmt"
	"log"
	"sync"
)

// Registerer is the interface for the part of a registry in charge of
// registering.
type Registerer interface {
	// Register registers a new Collector to be included in metrics
	// collection.
	//
	// Register should return an error if the descriptors provided by the
	// Collector are invalid.
	Register(Collector) error
	// MustRegister works like Register but registers any number of
	// Collectors and panics upon the first registration that causes an
	// error.
	MustRegister(...Collector)
	// TODO: add ability to unregister a Collector
	// Unregister(Collector) bool
}

// Gatherer is the interface for the part of a registry in charge of gathering
// the collected metrics.
type Gatherer interface {
	// Gather calls the Collect method of the registered Collectors and then
	// gathers the collected metrics.
	//
	// TODO: extend the return value to handle Metrics other than Gauges
	Gather() (map[string]float64, error)
}

var (
	defaultRegistry              = NewRegistry()
	DefaultRegisterer Registerer = defaultRegistry
	DefaultGatherer   Gatherer   = defaultRegistry
)

// MustRegister registers the provided Collectors with the DefaultRegisterer and
// panics if any error occurs.
func MustRegister(cs ...Collector) {
	DefaultRegisterer.MustRegister(cs...)
}

// Gather gathers the registered Collector's Metrics with the DefaultGatherer.
func Gather() (map[string]float64, error) {
	return DefaultGatherer.Gather()
}

// NewRegistry creates a new vanilla Registry without any Collectors
// pre-registered.
func NewRegistry() *Registry {
	return &Registry{
		collectorsByName: map[string]Collector{},
	}
}

// Registry registers collectors, collects their metrics, and gathers
// them into a map of metric name to metric value.  It implements Registerer.
// The zero value is not usable. Create instances with NewRegistry
type Registry struct {
	mtx sync.RWMutex
	// map descriptor names to Collector
	collectorsByName map[string]Collector
}

// MustRegister implements Registerer.
func (r *Registry) MustRegister(cs ...Collector) {
	for _, c := range cs {
		if err := r.Register(c); err != nil {
			panic(err)
		}
	}
}

// Register implements Registerer.
func (r *Registry) Register(c Collector) error {
	descChan := make(chan *Desc, 10)
	go func() {
		c.Describe(descChan)
		close(descChan)
	}()
	r.mtx.Lock()
	defer func() {
		// Drain channel in case of premature return to not leak a goroutine.
		for range descChan {
		}
		r.mtx.Unlock()
	}()

	for desc := range descChan {
		// Is the descriptor valid at all?
		if desc.err != nil {
			return fmt.Errorf("descriptor %s is invalid: %s", desc, desc.err)
		}
		if desc.name == "" {
			return fmt.Errorf("descriptor %s is invalid: %s", desc, "no name")
		}
		// TODO: more error handling -  i.e. descriptor uniqueness

		r.collectorsByName[desc.name] = c

		// Print metadata to stdout. Not idiomatic but ok for this challenge
		log.Printf("# TYPE %s %s\n", desc.name, desc.t)
	}

	return nil
}

// Gather implements Gatherer.
func (r *Registry) Gather() (map[string]float64, error) {
	var (
		metricChan = make(chan Metric, 1)
		wg         sync.WaitGroup
	)

	r.mtx.RLock()
	goroutineBudget := len(r.collectorsByName)
	metricByName := make(map[string]float64, len(r.collectorsByName))
	collectors := make(chan Collector, len(r.collectorsByName))
	for _, collector := range r.collectorsByName {
		collectors <- collector
	}
	r.mtx.RUnlock()

	wg.Add(goroutineBudget)

	collectWorker := func() {
		for {
			select {
			case collector := <-collectors:
				collector.Collect(metricChan)
			default:
				return
			}
			wg.Done()
		}
	}

	// Start the first worker now to make sure at least one is running.
	go collectWorker()
	goroutineBudget--

	// Close metricChan once all collectors are collected.
	go func() {
		wg.Wait()
		close(metricChan)
	}()

	// Drain metricChan in case of premature return.
	defer func() {
		if metricChan != nil {
			for range metricChan {
			}
		}
	}()

	// Copy the channel reference so we can nil them out later to remove
	// them from the select statements below.
	mc := metricChan

	for {
		select {
		case metric, ok := <-mc:
			if !ok {
				mc = nil
				break
			}
			processMetric(metric, metricByName)
		default:
			if goroutineBudget <= 0 || len(collectors) == 0 {
				// All collectors are already being worked on or
				// we have already as many goroutines started as
				// there are collectors. Do the same as above,
				// just without the default.
				select {
				case metric, ok := <-mc:
					if !ok {
						mc = nil
						break
					}
					processMetric(metric, metricByName)
				}
				break
			}
			// Start more workers.
			go collectWorker()
			goroutineBudget--
		}
		// Once both metricChan is closed and drained,
		// the contraption above will nil out mc,
		// and then we can leave the collect loop here.
		if mc == nil {
			break
		}
	}

	return metricByName, nil
}

// processMetric is an internal helper method only used by the Gather metric.
func processMetric(metric Metric, metricByName map[string]float64) {
	val, _ := metric.Write()
	metricByName[metric.Desc().name] = val
}
