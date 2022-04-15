package server

import (
	"io"
	"net/http"
)

func health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `{"alive": true}`)
}
