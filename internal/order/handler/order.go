package handler

import (
	"github.com/aaanger/ecommerce/internal/order/model"
	"github.com/aaanger/ecommerce/internal/order/service"
	"github.com/aaanger/ecommerce/pkg/middleware"
	"github.com/aaanger/ecommerce/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type OrderHandler struct {
	service service.IOrderService
}

func NewOrderHandler(service service.IOrderService) *OrderHandler {
	return &OrderHandler{
		service: service,
	}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req model.CreateOrderReq

	err := c.BindJSON(&req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Failed to parse request body")
		return
	}

	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "user id not found")
		return
	}

	order, err := h.service.CreateOrder(userID, &req)
	if err != nil {
		logrus.Error(err)
		response.Error(c, http.StatusInternalServerError, "Failed to create order")
		return
	}

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

	order, err := h.service.GetOrderByID(userID, orderID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Order not found")
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

	order, err := h.service.UpdateOrderStatus(req.UserID, orderID, req.Status)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update order status")
		return
	}

	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) CancelOrder(c *gin.Context) {
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

	order, err := h.service.CancelOrder(userID, orderID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to cancel order")
		return
	}

	response.JSON(c, http.StatusOK, order)
}
