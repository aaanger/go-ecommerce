package model

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Amount      int     `json:"amount"`
	InStock     bool    `json:"in_stock"`
}

type UpdateProduct struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Price       *int    `json:"price"`
	Amount      *int    `json:"amount"`
	InStock     *bool   `json:"in_stock"`
}

type ProductReq struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Price       float64 `json:"price" binding:"required"`
	Amount      int     `json:"amount" binding:"required"`
	InStock     bool    `json:"in_stock" binding:"required"`
}

type ProductSearchReq struct {
	Category string
	Search   string
	MinPrice float64
	MaxPrice float64
	Page     int
	Limit    int
}
