package model

import "time"

type CreatePaymentReq struct {
	Amount       Amount            `json:"amount"`
	Capture      bool              `json:"capture"`
	Confirmation ConfirmationReq   `json:"confirmation"`
	Metadata     map[string]string `json:"metadata"`
	Description  string            `json:"description"`
}

type CreatePaymentRes struct {
	ID           string            `json:"id"`
	Status       string            `json:"status"`
	Paid         bool              `json:"paid"`
	Amount       Amount            `json:"amount"`
	Confirmation ConfirmationRes   `json:"confirmation"`
	CreatedAt    time.Time         `json:"created_at"`
	Description  string            `json:"description"`
	Metadata     map[string]string `json:"metadata"`
	Recipient    Recipient         `json:"recipient"`
	Refundable   bool              `json:"refundable"`
	Test         bool              `json:"test"`
}

type Amount struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}

type ConfirmationReq struct {
	Type      string `json:"type"`
	ReturnURL string `json:"return_url"`
}

type ConfirmationRes struct {
	Type            string `json:"type"`
	ConfirmationURL string `json:"confirmation_url"`
}

type Recipient struct {
	AccountID string `json:"account_id"`
	GatewayID string `json:"gateway_id"`
}
