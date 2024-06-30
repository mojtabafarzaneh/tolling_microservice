package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"

	"github.com/mojtabafarzaneh/tolling/types"
)

func main() {
	listenAdder := flag.String("listenAdder", ":3000", "the listen adder of the http server")
	flag.Parse()

	var (
		store = NewMemoryStore()
		svc   = NewInvoiceAggregator(store)
	)
	svc = NewLogMiddleware(svc)

	makeHttpTransport(*listenAdder, svc)
}

func DistanceAgg(svc Aggregator) {
	panic("unimplemented")
}

func AggregateDistance(svc Aggregator) {
	panic("unimplemented")
}

func makeHttpTransport(listenAdder string, svc Aggregator) {
	fmt.Println("http transport running on the port", listenAdder)
	http.HandleFunc("/aggregate", handleaggregate(svc))
	http.ListenAndServe(listenAdder, nil)

}

func handleaggregate(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

func writeJSON(rw http.ResponseWriter, status int, v any) error {
	rw.WriteHeader(status)
	rw.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(rw).Encode(v)

}
