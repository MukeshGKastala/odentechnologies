package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/odentechnologies/metric/metric"
	"github.com/odentechnologies/metric/server"
)

const cdma = "cable_diameter_moving_average"

// Modeled metric collection after the Prometheus Go client library.
// Took this as an opportunity to learn about the inner workings of ^
var cableDiameterMovingAverage = metric.NewGauge(cdma)

func init() {
	metric.MustRegister(cableDiameterMovingAverage)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port

	// Idiomatic way to tell the goroutines in the upstream stages to abandon what they're doing.
	done := make(chan struct{})
	defer close(done)

	// One minute moving average that updates every second.
	// There is a ~1.5 minute ramp-up period - average is 0.0 until then.
	// Why ~1.5 minute and not ~1 minute(window * interval -> 60 * 1s)?
	// The method that polls the Oden backend waits the interval after function execution.
	// Meaning the elapsed time to calculate an average is
	// ~1 minute + total execution time of function calls - left as a TODO in the poll package.
	go calculateMovingAverage(60, 1000*time.Millisecond, done)
	// Offset metric gathering to avoid converging on periodic behavior.
	// This prints to STDOUT.
	go gatherMetrics(1500*time.Millisecond, done)

	srv := &http.Server{
		Addr:    addr,
		Handler: server.MakeHTTPHandler(),
	}

	signalCh := make(chan os.Signal, 1)
	// Interrupt handles Ctrl+C locally on all operating systems.
	// SIGTERM handles a cloud environment's termination signal.
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("http server running on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server exited: %s", err)
		}
	}()

	sig := <-signalCh
	log.Printf("server shutting down: %s signal caught", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("failed to shutdown server: %s", err)
	}
}
