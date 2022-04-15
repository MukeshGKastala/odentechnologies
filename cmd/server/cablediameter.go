package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/odentechnologies/metric/ma"
	"github.com/odentechnologies/metric/metric"
	"github.com/odentechnologies/metric/pkg/poll"
)

func calculateMovingAverage(window int, interval time.Duration, done <-chan struct{}) {
	sma := ma.NewSMA(window)

	// Calculate and store the moving average
	poll.Until(func() {
		// TODO: URL should come from environment variables
		resp, _ := http.Get("http://takehome-backend.oden.network/?metric=cable-diameter")
		defer resp.Body.Close()

		metric := struct {
			Diameter float64 `json:"value"`
		}{}
		decoder := json.NewDecoder(resp.Body)
		decoder.Decode(&metric)

		sma.Add(metric.Diameter)
		cableDiameterMovingAverage.Set(sma.Average())
	}, interval, done)
}

func gatherMetrics(interval time.Duration, done <-chan struct{}) {
	poll.Until(func() {
		metrics, _ := metric.Gather()
		log.Printf("%s %.2f\n", cdma, metrics[cdma])
	}, interval, done)
}
