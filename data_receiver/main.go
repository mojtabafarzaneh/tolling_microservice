package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/mojtabafarzaneh/tolling/types"
)

type DataReceiver struct {
	msgch chan types.ObuData
	conn  *websocket.Conn
	prod  DataProducer
}

var kafkaTopic = "obudata"

func main() {
	recv, err := NewDataReceiver()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/ws", recv.handleWS)
	http.ListenAndServe(":30000", nil)

}

func NewDataReceiver() (*DataReceiver, error) {
	var (
		p          DataProducer
		err        error
		kafkaTopic = "obudata"
	)

	p, err = NewKafkaProducer(kafkaTopic)
	if err != nil {
		return nil, err
	}

	p = NewLogMiddleWare(p)
	return &DataReceiver{
		msgch: make(chan types.ObuData, 128),
		prod:  p,
	}, nil
}

func (dr *DataReceiver) ProduceData(data types.ObuData) error {
	return dr.prod.ProduceData(data)
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
		if err := dr.ProduceData(data); err != nil {
			log.Println("kafka produce error: ", err)
		}
	}
}
