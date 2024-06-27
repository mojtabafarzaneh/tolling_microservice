package main

import (
	"math"

	"github.com/mojtabafarzaneh/tolling/types"
)

// end your interfaces with "er"
type CalculateServicer interface {
	CalculateDistance(data types.ObuData) (float64, error)
}

type CalService struct {
	prevPoint []float64
}

func NewCalService() CalculateServicer {
	return &CalService{}
}

func (s *CalService) CalculateDistance(data types.ObuData) (float64, error) {
	distance := 0.0

	if len(s.prevPoint) > 0 {
		distance = calculateDistance(s.prevPoint[0], s.prevPoint[1], data.Lat, data.Long)
	}
	s.prevPoint = []float64{data.Lat, data.Long}
	return distance, nil

}

func calculateDistance(x1, x2, y1, y2 float64) float64 {
	return math.Sqrt(math.Pow(x2-x1, 2) + math.Pow(y2-y1, 2))
}
