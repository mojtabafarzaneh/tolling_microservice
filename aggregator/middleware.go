package main

import (
	"time"

	"github.com/mojtabafarzaneh/tolling/types"
	"github.com/sirupsen/logrus"
)

type LogMiddleware struct {
	next Aggregator
}

// CalculateInvoice implements Aggregator.

func NewLogMiddleware(next Aggregator) *LogMiddleware {
	return &LogMiddleware{
		next: next,
	}
}

func (l *LogMiddleware) DistanceAggregator(data types.Distance) (err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"err":  err,
			"took": time.Since(start),
		}).Info()
	}(time.Now())

	err = l.next.DistanceAggregator(data)
	return
}

func (l *LogMiddleware) CalculateInvoice(obuID int) (inv *types.Invoice, err error) {
	defer func(start time.Time) {
		var (
			distance float64
			amount   float64
		)
		if inv != nil {
			distance = inv.TotalDistance
			amount = inv.TotalAmount
		}
		logrus.WithFields(logrus.Fields{
			"error":         err,
			"took":          time.Since(start),
			"obuID":         obuID,
			"totalDistance": distance,
			"totalAmount":   amount,
		}).Info()
	}(time.Now())

	inv, err = l.next.CalculateInvoice(obuID)
	return
}
