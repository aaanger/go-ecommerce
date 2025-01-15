package broker

import (
	"errors"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Producer struct {
	Producer *kafka.Producer
}

func NewProducer(address string) (*Producer, error) {
	conf := &kafka.ConfigMap{
		"bootstrap.servers": address,
	}

	p, err := kafka.NewProducer(conf)
	if err != nil {
		return nil, err
	}

	return &Producer{
		Producer: p,
	}, nil
}

func (p *Producer) Produce(topic string, value []byte) error {
	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: value,
	}

	deliveryChan := make(chan kafka.Event)

	err := p.Producer.Produce(msg, deliveryChan)
	if err != nil {
		return err
	}

	e := <-deliveryChan

	switch e.(type) {
	case *kafka.Message:
		return nil
	case kafka.Error:
		return err
	default:
		return errors.New("unknown kafka event type")
	}
}
