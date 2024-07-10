package main

import (
	"log"

	"github.com/mojtabafarzaneh/tolling/aggregator/client"
)

const kafkaTopic = "obudata"

//if you want to use the http transport uncomment the code below
//const aggregatorNewClient = "http://127.0.0.1:3000/aggregate"

func main() {
	var (
		err error
		svc CalculateServicer
	)
	svc = NewCalService()
	svc = NewLogMiddleware(svc)
	//and uncomment this code below as well
	//httpClient := client.NewHTTPClient(aggregatorNewClient)

	//if you want to use the http transport comment the line below
	grpcClient, err := client.NewGRPCClien(":5000")
	if err != nil {
		log.Fatal(err)
	}
	kafkaConsumer, err := NewKafkaConsumer(kafkaTopic, svc, grpcClient)
	if err != nil {
		log.Fatal(err)
	}
	kafkaConsumer.Start()
}
