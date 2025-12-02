package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"sync"
)

type Consumer struct {
	reader *kafka.Reader
	log    *zap.Logger
}

func NewConsumer(reader *kafka.Reader, log *zap.Logger) *Consumer {
	return &Consumer{
		reader: reader,
		log:    log,
	}
}

func (c *Consumer) Consume(ctx context.Context, handler func(msg kafka.Message) error, workers int) {
	defer c.reader.Close()

	msgChan := make(chan kafka.Message, 50)

	var wg sync.WaitGroup
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func(id int) {
			defer wg.Done()
			for msg := range msgChan {
				if err := handler(msg); err != nil {
					c.log.Error("Kafka consumer: error handling message", zap.Error(err), zap.Any("kafka_message", msg), zap.Int("worker_id", id))
					continue
				}
				if err := c.reader.CommitMessages(ctx, msg); err != nil {
					c.log.Error("Kafka consumer: error commiting message", zap.Error(err), zap.Any("kafka_message", msg), zap.Int("worker_id", id))
				}
			}
		}(i)
	}

	go func() {
		<-ctx.Done()
		wg.Wait()
		close(msgChan)
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			msg, err := c.reader.FetchMessage(ctx)
			if err != nil {
				c.log.Error("Kafka consumer: error fetching message", zap.Error(err))
				continue
			}
			msgChan <- msg
		}
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
