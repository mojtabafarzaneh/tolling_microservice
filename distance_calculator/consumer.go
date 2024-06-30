package main

import (
	"encoding/json"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/mojtabafarzaneh/tolling/aggregator/client"
	"github.com/mojtabafarzaneh/tolling/types"
	"github.com/sirupsen/logrus"
)

type KafkaConsumer struct {
	consumer   *kafka.Consumer
	isRunning  bool
	calservice CalculateServicer
	aggClient  *client.Client
}

func NewKafkaConsumer(topic string, svc CalculateServicer, aggClient *client.Client) (*KafkaConsumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, err
	}

	c.SubscribeTopics([]string{topic}, nil)

	return &KafkaConsumer{
		consumer:   c,
		calservice: svc,
		aggClient:  aggClient,
	}, nil
}

func (c *KafkaConsumer) Start() {
	logrus.Info("kafka transport started")
	c.isRunning = true
	c.ReadMessageLoop()
}

func (c *KafkaConsumer) ReadMessageLoop() {
	for c.isRunning {
		msg, err := c.consumer.ReadMessage(-1)
		if err != nil {
			logrus.Errorf("kafka consume error %s", err)
			continue
		}
		var data types.ObuData
		if err = json.Unmarshal(msg.Value, &data); err != nil {
			logrus.Errorf("JSON Unmarshal error %s", err)
			continue
		}
		distance, err := c.calservice.CalculateDistance(data)
		if err != nil {
			logrus.Errorf("calculation distance had and error of %s\n", err)
		}
		req := &types.Distance{
			Unix:  time.Now().Unix(),
			OBUID: data.OBUID,
			Value: distance,
		}
		if err := c.aggClient.AggregatorInvoicer(*req); err != nil {
			logrus.Errorf("aggregator error %s", err)
			continue
		}

	}

}
