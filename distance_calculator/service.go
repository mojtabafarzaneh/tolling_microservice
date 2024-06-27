package main

import (
	"fmt"
	"math"

	"github.com/mojtabafarzaneh/tolling/types"
)

// end your interfaces with "er"
type CalculateServicer interface {
	CalculateDistance(data types.ObuData) (float64, error)
}

type CalService struct {
	points [][]float64
}

func NewCalService() CalculateServicer {
	return &CalService{
		points: make([][]float64, 0),
	}
}

func (s *CalService) CalculateDistance(data types.ObuData) (float64, error) {
	distance := 0.0
	s.points = append(s.points, []float64{data.Lat, data.Long})

	if len(s.points) > 0 {
		prevPoints := s.points[len(s.points)-1]
		fmt.Println(prevPoints)
		distance = calculateDistance(prevPoints[0], prevPoints[1], data.Lat, data.Long)
		fmt.Println(data.Lat, data.Long)
	}
	fmt.Println("calculating the distance")
	return distance, nil

}

func calculateDistance(x1, x2, y1, y2 float64) float64 {
	return math.Sqrt(math.Pow(x2-x1, 2) + math.Pow(y2-y1, 2))
}
