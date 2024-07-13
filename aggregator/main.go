package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/mojtabafarzaneh/tolling/types"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

func main() {

	var (
		store           = NewMemoryStore()
		svc             = NewInvoiceAggregator(store)
		grpcListenAdder = os.Getenv("AGG_GRPC_ENDPOINT")
		httpListenAdder = os.Getenv("AGG_HTTP_ENDPOINT")
	)
	svc = NewMetricsMiddleware(svc)
	svc = NewLogMiddleware(svc)
	go func() {
		log.Fatal(makeGRPCTransport(grpcListenAdder, svc))
	}()
	log.Fatal(makeHttpTransport(httpListenAdder, svc))
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
	aggMetricsHandler := newHTTPMetricsHandler("aggregate")
	invMetricsHandler := newHTTPMetricsHandler("invoice")
	fmt.Println("http transport running on the port", listenAdder)
	http.HandleFunc("/aggregate", aggMetricsHandler.instrument(handleaggregate(svc)))
	http.HandleFunc("/invoice", invMetricsHandler.instrument(handleIGetInvoice(svc)))
	http.Handle("/metrics", promhttp.Handler())

	return http.ListenAndServe(listenAdder, nil)

}
func writeJSON(rw http.ResponseWriter, status int, v any) error {
	rw.WriteHeader(status)
	rw.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(rw).Encode(v)

}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}
