package main

import (
	"encoding/json"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/mojtabafarzaneh/tolling/types"
	"github.com/sirupsen/logrus"
)

type KafkaConsumer struct {
	consumer   *kafka.Consumer
	isRunning  bool
	calservice CalculateServicer
}

func NewKafkaConsumer(topic string, svc CalculateServicer) (*KafkaConsumer, error) {
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
			logrus.Errorf("distance calculator error %s", err)
			continue
		}
		fmt.Printf("distance %.2f\n", distance)

	}

}
