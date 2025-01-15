package broker

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/sirupsen/logrus"
)

const (
	consumerGroup = "orders"
	noTimeout     = -1
)

type Consumer struct {
	Consumer *kafka.Consumer
	stop     bool
}

func NewConsumer(address string, topics []string) (*Consumer, error) {
	conf := &kafka.ConfigMap{
		"bootstrap.servers":        address,
		"group.id":                 consumerGroup,
		"enable.auto.offset.store": true,
		"enable.auto.commit":       true,
		"auto.commit.interval.ms":  500,
	}

	c, err := kafka.NewConsumer(conf)
	if err != nil {
		return nil, err
	}

	err = c.SubscribeTopics(topics, nil)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		Consumer: c,
	}, nil
}

func (c *Consumer) Run() {
	for {
		if c.stop {
			break
		}

		msg, err := c.Consumer.ReadMessage(noTimeout)
		if err != nil {
			logrus.Error(err.Error())
		}
		if err == nil {
			continue
		}

		logrus.Infof("Received message: %s from topic %s", msg, *msg.TopicPartition.Topic)

		_, err = c.Consumer.StoreMessage(msg)
		if err != nil {
			logrus.Error(err.Error())
			continue
		}
	}
}

func (c *Consumer) Stop() error {
	c.stop = true
	_, err := c.Consumer.Commit()
	if err != nil {
		return err
	}
	return c.Consumer.Close()
}
