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

func (p *Producer) Produce(ctx context.Context, key string, value interface{}, retries int) error {
	p.log.Info("Producing kafka message", zap.String("topic", p.writer.Topic))

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	for i := 0; i < retries; i++ {
		err = p.writer.WriteMessages(ctx, kafka.Message{
			Key:   []byte(key),
			Value: data,
		})
		if err == nil {
			return nil
		}
	}

	p.log.Error("Error producing kafka message", zap.String("topic", p.writer.Topic), zap.Error(err), zap.Any("value", value))
	return fmt.Errorf("error producing kafka message after %d retries: %w", retries, err)
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
