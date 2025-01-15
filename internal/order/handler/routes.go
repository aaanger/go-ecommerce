package handler

import (
	"database/sql"
	"github.com/aaanger/ecommerce/internal/order/repository"
	"github.com/aaanger/ecommerce/internal/order/service"
	repository2 "github.com/aaanger/ecommerce/internal/product/repository"
	"github.com/aaanger/ecommerce/pkg/middleware"
	"github.com/gin-gonic/gin"
)

func OrderRoutes(r *gin.Engine, db *sql.DB) {
	repo := repository.NewOrderRepository(db)
	productRepo := repository2.NewProductRepository(db)
	svc := service.NewOrderService(repo, productRepo)
	h := NewOrderHandler(svc)

	order := r.Group("/orders", middleware.UserIdentity)

	order.POST("/create", h.CreateOrder)
	order.GET("/:id", h.GetOrderByID)
	order.GET("/all", h.GetAllOrders)
	order.PUT("/cancel/:id", h.CancelOrder)

	updateStatus := order.Group("/update-status", middleware.ModeratorIdentity)
	{
		updateStatus.POST("/delivering/:id", h.OrderStatusDelivering)
		updateStatus.POST("/delivered/:id", h.OrderStatusDelivered)
	}
}
