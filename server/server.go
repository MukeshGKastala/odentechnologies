package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/odentechnologies/metric/metric"
)

func MakeHTTPHandler(gatherer metric.Gatherer) http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/", health).Methods("GET")
	r.HandleFunc("/health", health).Methods("GET")

	// TODO: add authentication and authorization middleware
	r.HandleFunc("/cable-diameter", cablediameter(gatherer)).Methods("GET")

	return r
}
