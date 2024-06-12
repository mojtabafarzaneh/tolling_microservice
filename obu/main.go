package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

//trying to simulate an on board unit that sends out GPS coordinance to each entervallue
//we are going to use websocket to sends the coordinance from this board to our first microservice

var sendInterVall = time.Second

type ObuData struct {
	OBUID int     `json:"obuID"`
	Lat   float64 `json:"lat"`
	Long  float64 `json:"long"`
}

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
		ids[i] = rand.Intn(math.MaxInt)
	}
	return ids
}
func main() {
	obIDs := genObuIDs(20)
	for {
		for i := 0; i < len(obIDs); i++ {
			lat, long := genLocation()
			data := ObuData{
				OBUID: obIDs[i],
				Lat:   lat,
				Long:  long,
			}
			fmt.Printf("%+v\n", data)
		}

		time.Sleep(sendInterVall)
	}
}

func init() {
	rand.Seed(time.Now().Unix())
}
