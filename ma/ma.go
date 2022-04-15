// Package ma implements moving averages.
package ma

type MovingAverage interface {
	// Add adds a value to the series.
	Add(float64)
	// Average returns the current value of the moving average.
	Average() float64
	// SetAvg sets the Moving Average's value. The method is intended
	// to seed the current average, shortly after initialization, with a
	// previously calculated and persisted average.
	SetAverage(float64)
}
