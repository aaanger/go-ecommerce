package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Producer struct {
	writer *kafka.Writer
	log    *zap.Logger
}

func NewProducer(writer *kafka.Writer, log *zap.Logger) *Producer {
	return &Producer{
		writer: writer,
		log:    log,
	}
}

func (p *Producer) Produce(ctx context.Context, topic, key string, value interface{}, retries int) error {
	p.log.Info("Producing kafka message", zap.String("topic", topic))

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	for i := 0; i < retries; i++ {
		err = p.writer.WriteMessages(ctx, kafka.Message{
			Topic: topic,
			Key:   []byte(key),
			Value: data,
		})
		if err == nil {
			return nil
		}
	}

	p.log.Error("Error producing kafka message", zap.String("topic", topic), zap.Error(err), zap.Any("value", value))
	return fmt.Errorf("error producing kafka message after %d retries: %w", retries, err)
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
