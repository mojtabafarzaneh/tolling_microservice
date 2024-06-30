package main

import (
	"log"

	"github.com/mojtabafarzaneh/tolling/aggregator/client"
)

const kafkaTopic = "obudata"
const aggregatorNewClient = "http://127.0.0.1:3000/aggregate"

func main() {
	var (
		err error
		svc CalculateServicer
	)
	svc = NewCalService()
	svc = NewLogMiddleware(svc)
	kafkaConsumer, err := NewKafkaConsumer(kafkaTopic, svc, client.NewClient(aggregatorNewClient))
	if err != nil {
		log.Fatal(err)
	}
	kafkaConsumer.Start()
}
