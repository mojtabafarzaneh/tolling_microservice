package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/mojtabafarzaneh/tolling/types"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type HTTPMetricsHandler struct {
	reqCounter prometheus.Counter
}

func newHTTPMetricsHandler(reqname string) *HTTPMetricsHandler {
	reqCounter := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: fmt.Sprintf("http_%s_%s", reqname, "request_counter"),
		Name:      "aggregator",
	})
	return &HTTPMetricsHandler{
		reqCounter: reqCounter,
	}

}

func (h HTTPMetricsHandler) instrument(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.reqCounter.Inc()
		next(w, r)
	}
}

func handleIGetInvoice(svc Aggregator) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		value, ok := r.URL.Query()["obu"]
		if !ok {
			writeJSON(rw, http.StatusBadRequest, map[string]string{"error": "obu id is missing"})
			return
		}
		obuID, err := strconv.Atoi(value[0])
		if err != nil {
			writeJSON(rw, http.StatusBadRequest, map[string]string{"error": "invalid obuid"})
			return
		}
		invoice, err := svc.CalculateInvoice(obuID)
		if err != nil {
			writeJSON(rw, http.StatusBadRequest, map[string]string{"error": err.Error()})

		}
		writeJSON(rw, http.StatusOK, invoice)
	}
}

func handleaggregate(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "post method is not alowed"})
			return
		}

		var distance types.Distance
		if err := json.NewDecoder(r.Body).Decode(&distance); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		if err := svc.DistanceAggregator(distance); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
	}
}
