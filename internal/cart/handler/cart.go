package handler

import (
	"github.com/aaanger/ecommerce/internal/cart/model"
	"github.com/aaanger/ecommerce/internal/cart/service"
	"github.com/aaanger/ecommerce/pkg/cookie"
	"github.com/aaanger/ecommerce/pkg/middleware"
	"github.com/aaanger/ecommerce/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CartHandler struct {
	service service.ICartService
}

func NewCartHandler(service service.ICartService) *CartHandler {
	return &CartHandler{
		service: service,
	}
}

func (h *CartHandler) GetCart(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	session, err := cookie.ReadCookie(c.Request, cookie.CookieSession)
	if err != nil {
		response.Error(c, http.StatusBadGateway, "Try to visit page later")
		return
	}

	cart, err := h.service.GetCartByUserID(userID, session)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "get cart error")
		return
	}

	response.JSON(c, http.StatusOK, cart)
}

func (h *CartHandler) AddProduct(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	session, err := cookie.ReadCookie(c.Request, cookie.CookieSession)
	if err != nil {
		response.Error(c, http.StatusBadGateway, "Try to visit page later")
		return
	}

	var input model.AddProductReq

	err = c.ShouldBindJSON(&input)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid input parameters")
		return
	}

	cart, err := h.service.AddProduct(userID, input.ProductID, input.Quantity, session)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Add product to cart error")
		return
	}

	response.JSON(c, http.StatusOK, cart)
}

func (h *CartHandler) DeleteProduct(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	session, err := cookie.ReadCookie(c.Request, cookie.CookieSession)
	if err != nil {
		response.Error(c, http.StatusBadGateway, "Try to visit page later")
		return
	}

	var input model.DeleteProductReq

	err = c.ShouldBindJSON(&input)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid input parameters")
		return
	}

	cart, err := h.service.DeleteProduct(userID, input.ProductID, session)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to delete product from the cart")
		return
	}

	response.JSON(c, http.StatusOK, cart)
}
