package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"

	"github.com/mojtabafarzaneh/tolling/types"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

func main() {
	httpListenAdder := flag.String("httplistenAdder", ":3000", "the listen adder of the http server")
	grpcListenAdder := flag.String("grpclistenAdder", ":5000", "the listen adder of the grpc server")
	flag.Parse()

	var (
		store = NewMemoryStore()
		svc   = NewInvoiceAggregator(store)
	)
	svc = NewMetricsMiddleware(svc)
	svc = NewLogMiddleware(svc)
	go func() {
		log.Fatal(makeGRPCTransport(*grpcListenAdder, svc))
	}()
	log.Fatal(makeHttpTransport(*httpListenAdder, svc))
}

func DistanceAgg(svc Aggregator) {
	panic("unimplemented")
}

func AggregateDistance(svc Aggregator) {
	panic("unimplemented")
}

func makeGRPCTransport(listenAdder string, svc Aggregator) error {
	fmt.Println("grpc transport runnig on port ", listenAdder)
	ln, err := net.Listen("tcp", listenAdder)
	if err != nil {
		return err
	}
	defer ln.Close()
	server := grpc.NewServer([]grpc.ServerOption{}...)
	types.RegisterAggregatorServer(server, NewGRPCServer(svc))

	return server.Serve(ln)
}

func makeHttpTransport(listenAdder string, svc Aggregator) error {
	fmt.Println("http transport running on the port", listenAdder)
	http.HandleFunc("/aggregate", handleaggregate(svc))
	http.HandleFunc("/invoice", handleIGetInvoice(svc))
	http.Handle("/metrics", promhttp.Handler())

	return http.ListenAndServe(listenAdder, nil)

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
