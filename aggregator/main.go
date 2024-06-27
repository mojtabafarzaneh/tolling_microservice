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
	store := NewMemoryStore()

	var (
		svc = NewInvoiceAggregator(store)
	)
	makeHttpTransport(*listenAdder, svc)
}

func makeHttpTransport(listenAdder string, svc Aggregator) {
	fmt.Println("http transport running on the port", listenAdder)
	http.HandleFunc("/aggregate", handleaggregate(svc))
	http.ListenAndServe(listenAdder, nil)

}

func handleaggregate(_ Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var distance types.Distance
		if err := json.NewDecoder(r.Body).Decode(&distance); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
}
