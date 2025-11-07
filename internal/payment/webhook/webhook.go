package webhook

import (
	"github.com/aaanger/ecommerce/internal/order/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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

}
