package email

import (
	"fmt"
	"github.com/aaanger/ecommerce/internal/order/model"
	"github.com/go-mail/mail/v2"
	"strconv"
)

type SMTPConfig struct {
	Host     string
	Port     string
	Username string
	Password string
}

type EmailService struct {
	DefaultSender string
	dialer        *mail.Dialer
}

type Email struct {
	From      string
	To        string
	Subject   string
	Plaintext string
}

func NewEmailService(sender string, cfg SMTPConfig) (*EmailService, error) {
	port, err := strconv.Atoi(cfg.Port)
	if err != nil {
		return nil, fmt.Errorf("new email service: %w", err)
	}

	es := EmailService{
		DefaultSender: sender,
		dialer:        mail.NewDialer(cfg.Host, port, cfg.Username, cfg.Password),
	}

	return &es, nil
}

func (es *EmailService) Send(email Email) error {
	msg := mail.NewMessage()
	msg.SetHeader("To", email.To)
	msg.SetHeader("From", email.From)
	msg.SetHeader("Subject", email.Subject)
	msg.SetBody("text/plain", email.Plaintext)

	err := es.dialer.DialAndSend(msg)
	if err != nil {
		return fmt.Errorf("email send: %w", err)
	}
	return nil
}

func (es *EmailService) CreateOrder(to string, order model.Order) error {
	email := Email{
		From:      es.DefaultSender,
		To:        to,
		Subject:   "Ваш заказ принят в обработку",
		Plaintext: fmt.Sprintf("Детали заказа: %v, Итого: %f рублей", order.Lines, order.TotalPrice),
	}

	err := es.Send(email)
	if err != nil {
		return fmt.Errorf("email create order: %w", err)
	}
	return nil
}
