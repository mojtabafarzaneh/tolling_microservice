package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gorilla/websocket"
	"github.com/mojtabafarzaneh/tolling/types"
)

type DataReceiver struct {
	msgch chan types.ObuData
	conn  *websocket.Conn
	prod *kafka.Producer
}

var  kafkaTopic = "obudata"

func main() {
	recv, err := NewDataReceiver()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/ws", recv.handleWS)
	http.ListenAndServe(":30000", nil)
	
}

func NewDataReceiver() (*DataReceiver, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
	if err != nil {
		panic(err)
	}

		
	// Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
					} else {
						fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
					}
				}
			}
			}()


	return &DataReceiver{
		msgch: make(chan types.ObuData, 128),
		prod: p,
	}, nil
}

func (dr *DataReceiver) ProduceData(data types.ObuData) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	

	err = dr.prod.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic: &kafkaTopic, 
			Partition: kafka.PartitionAny,
		},
			Value: b,
			}, nil)

	return err
}


func (dr *DataReceiver) handleWS(w http.ResponseWriter, r *http.Request) {
	u := websocket.Upgrader{
		ReadBufferSize:  1028,
		WriteBufferSize: 1028,
	}
	conn, err := u.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	dr.conn = conn

	go dr.wsReceiverLoop()

}

func (dr *DataReceiver) wsReceiverLoop() {
	fmt.Println("new OBU connected, client connected")
	for {
		var data types.ObuData
		if err := dr.conn.ReadJSON(&data); err != nil {
			log.Println("read error:", err)
			continue
		}
		if err := dr.ProduceData(data); err != nil{
			log.Println("kafka produce error: ", err)
		}
	}
}
