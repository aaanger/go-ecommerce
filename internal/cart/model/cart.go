package model

import (
	"github.com/aaanger/ecommerce/internal/product/model"
	"time"
)

type Cart struct {
	ID         int        `json:"id"`
	UserID     int        `json:"user_id"`
	Lines      []CartLine `json:"lines"`
	TotalPrice float64    `json:"total_price"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

type CartLine struct {
	ProductID int            `json:"product_id"`
	Product   *model.Product `json:"product"`
	Quantity  int            `json:"quantity"`
}

type AddProductReq struct {
	ProductID int `json:"product_id" binding:"required"`
	Quantity  int `json:"quantity" binding:"required"`
}

type DeleteProductReq struct {
	ProductID int `json:"product_id" binding:"required"`
}
