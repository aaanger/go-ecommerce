package handler

import (
	"github.com/aaanger/ecommerce/internal/order/model"
	"github.com/aaanger/ecommerce/internal/order/service"
	"github.com/aaanger/ecommerce/pkg/middleware"
	"github.com/aaanger/ecommerce/pkg/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type OrderHandler struct {
	service  service.IOrderService
	consumer *service.OrderConsumer
	log      *zap.Logger
}

func NewOrderHandler(service service.IOrderService, consumer *service.OrderConsumer, log *zap.Logger) *OrderHandler {
	return &OrderHandler{
		service:  service,
		consumer: consumer,
		log:      log,
	}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	log := h.log.With(
		zap.String("service", "order"),
		zap.String("layer", "handler"),
		zap.String("method", "CreateOrder"))

	var req model.CreateOrderReq

	err := c.BindJSON(&req)
	if err != nil {
		log.Error("Create order: failed to parse request",
			zap.Error(err),
			zap.Any("request body", c.Request.Body))
		response.Error(c, http.StatusBadRequest, "Failed to parse request body")
		return
	}

	userID, err := middleware.GetUserID(c)
	if err != nil {
		log.Warn("Create order: user ID not found",
			zap.Error(err))
		response.Error(c, http.StatusUnauthorized, "user id not found")
		return
	}

	email, err := middleware.GetUserEmail(c)
	if err != nil {
		log.Warn("Create order: user email not found", zap.Error(err))
		response.Error(c, http.StatusUnauthorized, "user email not found")
		return
	}

	log.Info("Creating order", zap.Int("userID", userID), zap.Any("request data", req))

	order, err := h.service.CreateOrder(c, userID, email, &req)
	if err != nil {
		log.Error("Failed to create order", zap.Error(err), zap.Any("request data", req))
		response.Error(c, http.StatusInternalServerError, "Failed to create order")
		return
	}

	log.Info("Order created successfully", zap.Any("order", order))
	response.JSON(c, http.StatusOK, order)
}

func (h *OrderHandler) GetOrderByID(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "user id not found")
		return
	}

	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid order id")
		return
	}

	order, err := h.service.GetOrderByID(orderID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Order not found")
		return
	}

	if order.UserID != userID {
		response.Error(c, http.StatusForbidden, "not your order")
		return
	}

	response.JSON(c, http.StatusOK, order)
}

func (h *OrderHandler) GetAllOrders(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "user id not found")
		return
	}

	orders, err := h.service.GetAllOrders(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Orders not found")
		return
	}

	var res []model.GetAllOrdersRes

	for i := range orders {
		res = append(res, model.GetAllOrdersRes{
			ID:         orders[i].ID,
			CreatedAt:  orders[i].CreatedAt,
			UpdatedAt:  orders[i].UpdatedAt,
			Status:     orders[i].Status,
			TotalPrice: orders[i].TotalPrice,
		})
	}

	response.JSON(c, http.StatusOK, res)
}

func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	var req model.UpdateOrderStatusReq

	err := c.BindJSON(&req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid input parameters")
		return
	}

	role, ok := c.Get("role")
	if !ok {
		response.Error(c, http.StatusUnauthorized, "role not found")
		return
	}

	if role != "moderator" {
		response.Error(c, http.StatusForbidden, "moderator role required")
		return
	}

	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid order id")
		return
	}

	order, err := h.service.UpdateOrderStatus(orderID, req.Status)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update order status")
		return
	}

	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) CancelOrder(c *gin.Context) {
	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid order id")
		return
	}

	err = h.service.CancelOrder(c, orderID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to cancel order")
		return
	}

	response.JSON(c, http.StatusOK, "order canceled")
}
