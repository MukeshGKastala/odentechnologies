// Influenced by https://github.com/VividCortex/ewma
package ma

// TODO: make method set thread-safe
type ema struct {
	// EMA window size
	window int
	// The multiplier factor by which the previous samples decay
	decay float64
	// The number of samples added to this instance
	n int
	// Current average
	average float64
}

// TODO: export function once ema's method set has been tested
//
// NewSMA constructs an Exponential Moving Average
// The EMA window size must be greater than zero
//
// // It is not safe to use the returned MovingAverage from multiple goroutines.
func newEMA(window int) MovingAverage {
	// TODO: validate window size?
	return &ema{
		window: window,
		decay:  2 / (float64(window) + 1),
	}
}

// Add updates Exponential Moving Average.
func (e *ema) Add(value float64) {
	switch {
	case e.n < e.window:
		e.average += value
		e.n++
	case e.n == e.window:
		// Calculate the Simple Moving Average(SMA) for the initial EMA value
		e.average = e.average / float64(e.window)
		e.average = (value * e.decay) + (e.average * (1 - e.decay))
		e.n++
	default:
		e.average = (value * e.decay) + (e.average * (1 - e.decay))
	}
}

// Average returns the current value of the average,
// or 0.0 if the series hasn't reached the window size.
// The operation is O(1) since the Exponential Moving Average
// is updated in Add.
func (e *ema) Average() float64 {
	if e.n <= e.window {
		return 0.0
	}
	return e.average
}

// SetAvg sets the Exponential Moving Average's value.
func (e *ema) SetAverage(value float64) {
	e.average = value
	if e.n <= e.window {
		e.n = e.window + 1
	}
}
