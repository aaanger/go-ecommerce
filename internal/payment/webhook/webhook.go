package webhook

import (
	"github.com/aaanger/ecommerce/internal/order/service"
	"github.com/aaanger/ecommerce/internal/payment/model"
	"github.com/aaanger/ecommerce/pkg/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type Handler struct {
	orderService service.IOrderService
	log          *zap.Logger
}

func NewWebhookHandler(orderService service.IOrderService, log *zap.Logger) *Handler {
	return &Handler{
		orderService: orderService,
		log:          log,
	}
}

func (h *Handler) Handle(c *gin.Context) {
	var webhook model.Webhook

	if err := c.BindJSON(&webhook); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if _, ok := webhook.Object.Metadata["order_id"]; !ok {
		c.Status(http.StatusNotFound)
		return
	}
	orderID, err := strconv.Atoi(webhook.Object.Metadata["order_id"])
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	switch webhook.Event {
	case "payment.succeeded":
		if err := h.orderService.ConfirmOrder(c.Request.Context(), orderID); err != nil {
			response.Error(c, http.StatusInternalServerError, "failed to confirm order")
			return
		}
	case "payment.canceled":
		if err := h.orderService.CancelOrder(c.Request.Context(), orderID); err != nil {
			response.Error(c, http.StatusInternalServerError, "failed to cancel order")
			return
		}
	}

	c.Status(http.StatusOK)
}
