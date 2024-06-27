package main

import (
	"time"

	"github.com/mojtabafarzaneh/tolling/types"
	"github.com/sirupsen/logrus"
)

type LogMiddleware struct {
	next CalculateServicer
}

func NewLogMiddleware(next CalculateServicer) CalculateServicer {
	return &LogMiddleware{
		next: next,
	}
}

func (m *LogMiddleware) CalculateDistance(data types.ObuData) (dist float64, err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took": time.Since(start),
			"err":  err,
			"dist": dist,
		}).Info()
	}(time.Now())
	dist, err = m.next.CalculateDistance(data)
	return
}
