package main

import (
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
}

const kafkaTopic = "obudata"

func main() {
	
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
	if err != nil {
		panic(err)
	}
	
	defer p.Close()
	
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

			// Produce messages to topic (asynchronously)
			topic := "myTopic"
			p.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte("test producing"),
			}, nil)
			return 

			
	recv := NewDataReceiver()
	http.HandleFunc("/ws", recv.handleWS)
	http.ListenAndServe(":30000", nil)
	
}

func NewDataReceiver() *DataReceiver {
	return &DataReceiver{
		msgch: make(chan types.ObuData, 128),
	}
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
		fmt.Printf("received OBU data [%d] :: <lat %.2f, long %.2f> \n", data.OBUID, data.Lat, data.Long)
		dr.msgch <- data

	}
}
