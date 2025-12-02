package handler

import (
	"github.com/aaanger/ecommerce/internal/product/model"
	"github.com/aaanger/ecommerce/internal/product/service"
	"github.com/aaanger/ecommerce/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ProductHandler struct {
	service service.IProductService
}

func NewProductHandler(service service.IProductService) *ProductHandler {
	return &ProductHandler{
		service: service,
	}
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var product model.ProductReq

	err := c.ShouldBindJSON(&product)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid input parameters")
		return
	}

	createdProduct, err := h.service.CreateProduct(&product)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create product")
		return
	}

	response.JSON(c, http.StatusOK, createdProduct)
}

func (h *ProductHandler) GetProducts(c *gin.Context) {
	products, err := h.service.GetAllProducts()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get products")
		return
	}

	response.JSON(c, http.StatusOK, products)
}

func (h *ProductHandler) GetProductByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
	}

	product, err := h.service.GetProductByID(id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get product")
		return
	}

	response.JSON(c, http.StatusOK, product)
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	var input model.UpdateProduct

	err = c.ShouldBindJSON(&input)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid input parameters")
		return
	}

	err = h.service.UpdateProduct(id, input)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update product")
		return
	}

	response.JSON(c, http.StatusOK, input)
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid product id")
		return
	}

	err = h.service.DeleteProduct(id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to delete product")
		return
	}

	response.JSON(c, http.StatusOK, id)
}
