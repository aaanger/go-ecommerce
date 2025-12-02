package model

type Webhook struct {
	Type   string `json:"type"`
	Event  string `json:"event"`
	Object struct {
		ID       string            `json:"id"`
		Status   string            `json:"status"`
		Paid     bool              `json:"paid"`
		Metadata map[string]string `json:"metadata"`
	} `json:"object"`
}
