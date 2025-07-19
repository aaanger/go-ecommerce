package service

import (
	"encoding/json"
	"fmt"
	"github.com/aaanger/ecommerce/internal/order/model"
	"github.com/aaanger/ecommerce/pkg/email"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type OrderConsumer struct {
	emailService *email.EmailService
	log          *zap.Logger
}

func NewOrderConsumer(es *email.EmailService, log *zap.Logger) *OrderConsumer {
	return &OrderConsumer{
		emailService: es,
		log:          log,
	}
}

func (c *OrderConsumer) HandleOrderCreated(msg kafka.Message) error {
	var order model.Order

	err := json.Unmarshal(msg.Value, &order)
	if err != nil {
		c.log.Error("Order consumer: error unmarshalling kafka message", zap.Error(err), zap.String("message_value", string(msg.Value)))
		return fmt.Errorf("consumer order created: %w", err)
	}

	err = c.emailService.CreateOrder(order.UserEmail, order)
	if err != nil {
		c.log.Error("Error sending email for creating order", zap.Error(err), zap.Any("order", order))
		return fmt.Errorf("consumer order created: %w", err)
	}

	c.log.Info("Create order consumed and sent email successfully", zap.Any("order", order), zap.String("email", order.UserEmail))

	return nil
}
