package server

import (
	"encoding/json"
	"math"
	"net/http"

	"github.com/odentechnologies/metric/metric"
)

type cableDiameterResponse struct {
	MovingAverage float64 `json:"moving_average"`
}

func cablediameter(w http.ResponseWriter, r *http.Request) {
	metrics, _ := metric.Gather()
	cdma := math.Floor(metrics["cable_diameter_moving_average"]*100) / 100
	b, _ := json.Marshal(&cableDiameterResponse{MovingAverage: cdma})

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)

}
