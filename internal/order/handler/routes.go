package handler

import (
	"database/sql"
	grpcorder "github.com/aaanger/ecommerce/internal/order/handler/grpc/product"
	"github.com/aaanger/ecommerce/internal/order/repository"
	"github.com/aaanger/ecommerce/internal/order/service"
	payment "github.com/aaanger/ecommerce/internal/payment/client"
	"github.com/aaanger/ecommerce/internal/payment/webhook"
	repository2 "github.com/aaanger/ecommerce/internal/product/repository"
	"github.com/aaanger/ecommerce/pkg/kafka"
	"github.com/aaanger/ecommerce/pkg/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func OrderRoutes(r *gin.Engine, db *sql.DB, producer *kafka.Producer, grpcClient *grpcorder.OrderGRPCClient, paymentClient *payment.Client, consumer *service.OrderConsumer, logger *zap.Logger) {
	repo := repository.NewOrderRepository(db, logger)
	productRepo := repository2.NewProductRepository(db)
	svc := service.NewOrderService(repo, productRepo, grpcClient, paymentClient, producer, logger)
	h := NewOrderHandler(svc, consumer, logger)

	webhookHandler := webhook.NewWebhookHandler(svc, logger)

	r.POST("/payment/webhook", webhookHandler.Handle)

	order := r.Group("/orders", middleware.UserIdentity)

	order.POST("/create", h.CreateOrder)
	order.GET("/:id", h.GetOrderByID)
	order.GET("/all", h.GetAllOrders)
	order.PUT("/cancel/:id", h.CancelOrder)

	updateStatus := order.Group("/update-status")
	{
		updateStatus.PUT("/:id", h.UpdateOrderStatus)
	}
}
