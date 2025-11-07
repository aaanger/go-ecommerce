package model

import (
	payment "github.com/aaanger/ecommerce/internal/payment/model"
	"github.com/aaanger/ecommerce/internal/product/model"
	"time"
)

const (
	StatusPending    = "Pending"
	StatusCreated    = "Created"
	StatusDelivering = "Delivering"
	StatusDelivered  = "Delivered"
	StatusCanceled   = "Canceled"
)

type Order struct {
	ID         int         `json:"id"`
	UserID     int         `json:"user_id"`
	UserEmail  string      `json:"user_email"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
	Lines      []OrderLine `json:"lines"`
	Status     string      `json:"status"`
	TotalPrice float64     `json:"total_price"`
}

type OrderLine struct {
	ID        int            `json:"id"`
	ProductID int            `json:"product_id"`
	Product   *model.Product `json:"product"`
	Quantity  int            `json:"quantity"`
	Price     float64        `json:"price"`
}

type OrderLineReq struct {
	ProductID int `json:"product_id" binding:"required"`
	Quantity  int `json:"quantity" binding:"required"`
}

type CreateOrderReq struct {
	Lines []OrderLineReq `json:"lines" binding:"required,dive,required"`
}

type CreateOrderRes struct {
	Order   *Order
	Payment *payment.CreatePaymentRes
}

type UpdateOrderStatusReq struct {
	UserID int    `json:"user_id" binding:"required"`
	Status string `json:"status" binding:"required"`
}

type GetAllOrdersRes struct {
	ID         int       `json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Status     string    `json:"status"`
	TotalPrice float64   `json:"total_price"`
}
