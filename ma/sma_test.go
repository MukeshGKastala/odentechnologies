package ma_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/odentechnologies/metric/ma"
	"github.com/stretchr/testify/require"
)

func TestSMA(t *testing.T) {
	tests := map[string]struct {
		window   int
		values   []float64
		averages []float64
	}{
		// Average should return 0.0 until enough values are added
		"window size greater than number of values": {
			window:   5,
			values:   []float64{10, 10, 10, 10},
			averages: []float64{0, 0, 0, 0},
		},
		"window size equals number of values": {
			window:   5,
			values:   []float64{10, 10, 10, 10, 10},
			averages: []float64{0, 0, 0, 0, 10},
		},
		"window size less than number of values": {
			window:   5,
			values:   []float64{10, 10, 10, 10, 10, 10},
			averages: []float64{0, 0, 0, 0, 10, 10},
		},
		"average with precision five": {
			window:   2,
			values:   []float64{1.2345, 6.789},
			averages: []float64{0, 4.01175},
		},
		// Average the most recent window
		"truly a simple moving average": {
			window:   4,
			values:   []float64{10, 10, 5, 5, 20, 20},
			averages: []float64{0, 0, 0, 7.5, 10, 12.5},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			sma := ma.NewSMA(tc.window)
			for i, value := range tc.values {
				sma.Add(value)
				require.Equal(t, tc.averages[i], sma.Average())
			}
		})
	}
}

func TestSMASetAverage(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	window := rand.Intn(10-1) + 1
	value := 1 + rand.Float64()*(50-1)

	sma := ma.NewSMA(window)
	require.Equal(t, 0.0, sma.Average())

	for i := 0; i < window; i++ {
		sma.Add(value)
	}
	require.InDelta(t, value, sma.Average(), .001)

	value = 1 + rand.Float64()*(50-1)
	sma.SetAverage(value)
	require.InDelta(t, value, sma.Average(), .001)

	sma.Add(10)
	expected := (10 + value*(float64(window)-1)) / float64(window)
	require.InDelta(t, expected, sma.Average(), .001)
}
