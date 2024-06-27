package main

import (
	"log"
)

const kafkaTopic = "obudata"

func main() {
	var (
		err error
		svc CalculateServicer
	)
	svc = NewCalService()
	kafkaConsumer, err := NewKafkaConsumer(kafkaTopic, svc)
	if err != nil {
		log.Fatal(err)
	}
	kafkaConsumer.Start()
}
