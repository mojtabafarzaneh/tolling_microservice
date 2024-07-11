package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
	"github.com/mojtabafarzaneh/tolling/types"
)

//trying to simulate an on board unit that sends out GPS coordinance to each entervallue
//we are going to use websocket to sends the coordinance from this board to our first microservice

var sendInterVall = time.Second * 6

const wsEndPoint = "ws://127.0.0.1:30000/ws"

func genCoord() float64 {
	n := float64(rand.Intn(100) + 1)
	f := rand.Float64()
	return n + f

}

func genLocation() (float64, float64) {
	return genCoord(), genCoord()
}

func genObuIDs(n int) []int {
	ids := make([]int, n)
	for i := 0; i < n; i++ {
		ids[i] = rand.Intn(999999)
	}
	return ids
}
func main() {
	obIDs := genObuIDs(20)
	conn, _, err := websocket.DefaultDialer.Dial(wsEndPoint, nil)
	if err != nil {
		log.Fatal(err)
	}
	for {
		for i := 0; i < len(obIDs); i++ {
			lat, long := genLocation()
			data := types.ObuData{
				OBUID: obIDs[i],
				Lat:   lat,
				Long:  long,
			}
			if err := conn.WriteJSON(data); err != nil {
				log.Fatal(err)
			}
		}

		time.Sleep(sendInterVall)
		fmt.Printf("%d data has been sent\n", len(obIDs))
	}
}

func init() {
	rand.Seed(time.Now().Unix())
}
