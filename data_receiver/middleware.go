package main

import (
	"time"

	"github.com/mojtabafarzaneh/tolling/types"
	"github.com/sirupsen/logrus"
)

type LogMiddleWare struct {
	next DataProducer
}

func NewLogMiddleWare(next DataProducer) *LogMiddleWare {
	return &LogMiddleWare{
		next: next,
	}
}

func (l *LogMiddleWare) ProduceData(data types.ObuData) error {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"obuID": data.OBUID,
			"lat": data.Lat,
			"long":  data.Long,
			"took":  time.Since(start),
		}).Info("producing")
	}(time.Now())

	return l.next.ProduceData(data)
}
