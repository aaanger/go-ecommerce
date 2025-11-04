package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aaanger/ecommerce/internal/payment/model"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type Client struct {
	ShopID      string
	SecretKey   string
	APIEndpoint string
	HTTPClient  *http.Client
}

func NewClient(shopID, secretKey string) *Client {
	return &Client{
		ShopID:      shopID,
		SecretKey:   secretKey,
		APIEndpoint: "https://api.yookassa.ru/v3/",
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) CreatePayment(ctx context.Context, req *model.CreatePaymentReq) (*model.CreatePaymentRes, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("create payment: %w", err)
	}

	r, err := http.NewRequestWithContext(ctx, "POST", c.APIEndpoint+"payments", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create payment: %w", err)
	}

	idempotenceKey := uuid.New().String()

	r.Header.Set("Idempotence-Key", idempotenceKey)
	r.Header.Set("Content-Type", "application/json")
	r.SetBasicAuth(c.ShopID, c.SecretKey)

	res, err := c.HTTPClient.Do(r)
	if err != nil {
		return nil, fmt.Errorf("create payment: %w", err)
	}

	var paymentRes model.CreatePaymentRes
	if err = json.NewDecoder(res.Body).Decode(&paymentRes); err != nil {
		return nil, fmt.Errorf("create payment: %w", err)
	}

	return &paymentRes, nil
}
