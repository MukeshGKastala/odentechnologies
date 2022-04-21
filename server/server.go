package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

func MakeHTTPHandler() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/", health).Methods("GET")
	r.HandleFunc("/health", health).Methods("GET")

	// TODO: add authentication and authorization middleware
	r.HandleFunc("/cable-diameter", cablediameter).Methods("GET")

	return r
}
