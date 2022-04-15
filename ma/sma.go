// Influenced by https://github.com/mxmCherry/movavg
package ma

// TODO: make method set thread-safe
type sma struct {
	// SMA window size
	window int
	// Stored values - circular slice
	values []float64
	// Last set circular slice index
	i int
	// Current average
	average float64
}

// NewSMA constructs a Simple Moving Average
//
// The SMA window size must be greater than zero
// It is not safe to use the returned MovingAverage from multiple goroutines.
func NewSMA(window int) MovingAverage {
	// TODO: validate window size?
	return &sma{
		window: window,
		values: make([]float64, 0, window),
	}
}

// Add adds the value to the series and updates the Simple Moving Average.
func (s *sma) Add(value float64) {
	n := len(s.values)
	switch {
	case n == s.window:
		s.i = (s.i + 1) % s.window
		s.average += (value - s.values[s.i]) / float64(s.window)
		s.values[s.i] = value
	case n != 0:
		s.i++
		s.average = (value + float64(len(s.values))*s.average) / float64(len(s.values)+1)
		s.values = append(s.values, value)
	default:
		// initial insert
		s.average = value
		s.values = append(s.values, value)
	}
}

// Average returns the current value of the average,
// or 0.0 if the series hasn't reached the window size.
// The operation is O(1) since the Simple Moving Average
// is updated in Add.
func (s *sma) Average() float64 {
	if len(s.values) < s.window {
		return 0.0
	}
	return s.average
}

// SetAvg sets the Simple Moving Average's value.
func (s *sma) SetAverage(value float64) {
	s.values = make([]float64, s.window)
	for i := 0; i < s.window; i++ {
		s.values[i] = value
	}
	s.average = value
}
