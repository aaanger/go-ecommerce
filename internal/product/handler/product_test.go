package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/aaanger/ecommerce/internal/product/model"
	"github.com/aaanger/ecommerce/internal/product/service/mocks"
	"github.com/aaanger/ecommerce/pkg/lib"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

type ProductHandlerSuite struct {
	suite.Suite
	service *mocks.IProductService
	handler *ProductHandler
}

func (suite *ProductHandlerSuite) SetupTest() {
	suite.service = mocks.NewIProductService(suite.T())
	suite.handler = NewProductHandler(suite.service)
}

func TestProductHandlerSuite(t *testing.T) {
	suite.Run(t, new(ProductHandlerSuite))
}

func (suite *ProductHandlerSuite) TestHandler_CreateProductSuccess() {
	req := &model.ProductReq{
		Name:        "test",
		Description: "test",
		Price:       1,
		Amount:      1,
		InStock:     true,
	}

	res := &model.Product{
		ID:          1,
		Name:        "test",
		Description: "test",
		Price:       1,
		Amount:      1,
		InStock:     true,
	}

	suite.service.On("CreateProduct", req).Return(res, nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", 1)
		c.Set("role", "moderator")
		c.Next()
	})
	router.POST("/create", suite.handler.CreateProduct)

	requestBody, _ := json.Marshal(req)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/create", bytes.NewBuffer(requestBody))

	router.ServeHTTP(w, r)

	var response interface{}
	var productRes model.Product
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	lib.Copy(&productRes, &response)

	suite.Equal(http.StatusOK, w.Code)
	suite.Equal(res, &productRes)
}

func (suite *ProductHandlerSuite) TestHandler_CreateProductEmptyFields() {
	req := &model.ProductReq{
		Name:        "test",
		Description: "test",
		Price:       1,
		Amount:      1,
		InStock:     true,
	}

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", 1)
		c.Set("role", "moderator")
		c.Next()
	})
	router.POST("/create", suite.handler.CreateProduct)

	requestBody, _ := json.Marshal(req)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/create", bytes.NewBuffer(requestBody))

	router.ServeHTTP(w, r)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Equal(`"Invalid input parameters"`, w.Body.String())
}

func (suite *ProductHandlerSuite) TestHandler_CreateProductServiceFailure() {
	req := &model.ProductReq{
		Name:        "test",
		Description: "test",
		Price:       1,
		Amount:      1,
		InStock:     true,
	}

	suite.service.On("CreateProduct", req).Return(nil, errors.New("error"))

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", 1)
		c.Set("role", "moderator")
		c.Next()
	})
	router.POST("/create", suite.handler.CreateProduct)

	requestBody, _ := json.Marshal(req)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/create", bytes.NewBuffer(requestBody))

	router.ServeHTTP(w, r)

	suite.Equal(http.StatusInternalServerError, w.Code)
	suite.Equal(`"Failed to create product"`, w.Body.String())
}

// =====================================================================================================================

func (suite *ProductHandlerSuite) TestHandler_GetAllProductsSuccess() {
	res := []model.Product{
		{
			ID:          1,
			Name:        "1",
			Description: "1",
		},
		{
			ID:          2,
			Name:        "2",
			Description: "2",
		},
	}

	suite.service.On("GetAllProducts").Return(res, nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", 1)
		c.Set("role", "moderator")
		c.Next()
	})
	router.GET("/all", suite.handler.GetAllProducts)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/all", nil)

	router.ServeHTTP(w, r)

	var response interface{}
	var productsRes []model.Product
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	lib.Copy(&productsRes, &response)
	suite.Equal(http.StatusOK, w.Code)
	suite.Equal(res, productsRes)
}

func (suite *ProductHandlerSuite) TestHandler_GetAllProductsServiceFailure() {
	suite.service.On("GetAllProducts").Return(nil, errors.New("error"))

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", 1)
		c.Set("role", "moderator")
		c.Next()
	})
	router.GET("/all", suite.handler.GetAllProducts)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/all", nil)

	router.ServeHTTP(w, r)

	suite.Equal(http.StatusInternalServerError, w.Code)
	suite.Equal(`"Failed to get products"`, w.Body.String())
}

// =====================================================================================================================\

func (suite *ProductHandlerSuite) TestHandler_GetProductByIDSuccess() {
	res := &model.Product{
		ID:          1,
		Name:        "test",
		Description: "test",
		Price:       1,
		Amount:      1,
		InStock:     true,
	}

	suite.service.On("GetProductByID", 1).Return(res, nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", 1)
		c.Set("role", "moderator")
		c.Next()
	})
	router.GET("/:id", suite.handler.GetProductByID)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/1", nil)

	router.ServeHTTP(w, r)

	var response interface{}
	var productRes model.Product
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	lib.Copy(&productRes, &response)

	suite.Equal(http.StatusOK, w.Code)
	suite.Equal(res, &productRes)
}

func (suite *ProductHandlerSuite) TestHandler_GetProductByIDServiceFailure() {
	suite.service.On("GetProductByID", 1).Return(nil, errors.New("error"))

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", 1)
		c.Set("role", "moderator")
		c.Next()
	})
	router.GET("/:id", suite.handler.GetProductByID)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/1", nil)

	router.ServeHTTP(w, r)

	suite.Equal(http.StatusInternalServerError, w.Code)
	suite.Equal(`"Failed to get product"`, w.Body.String())
}

// =====================================================================================================================

func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func boolPtr(b bool) *bool {
	return &b
}

func (suite *ProductHandlerSuite) TestHandler_UpdateProductSuccess() {
	req := model.UpdateProduct{
		Name:        stringPtr("test"),
		Description: stringPtr("test"),
		Price:       intPtr(1),
		Amount:      intPtr(1),
		InStock:     boolPtr(true),
	}

	suite.service.On("UpdateProduct", 1, req).Return(nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", 1)
		c.Set("role", "moderator")
		c.Next()
	})
	router.PUT("/:id", suite.handler.UpdateProduct)

	requestBody, _ := json.Marshal(req)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("PUT", "/1", bytes.NewBuffer(requestBody))

	router.ServeHTTP(w, r)

	var response interface{}
	var productRes model.UpdateProduct
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	lib.Copy(&productRes, &response)

	suite.Equal(http.StatusOK, w.Code)
	suite.Equal(productRes, req)
}

func (suite *ProductHandlerSuite) TestHandler_UpdateProductServiceFailure() {
	req := model.UpdateProduct{
		Name:        stringPtr("test"),
		Description: stringPtr("test"),
		Price:       intPtr(1),
		Amount:      intPtr(1),
		InStock:     boolPtr(true),
	}

	suite.service.On("UpdateProduct", 1, req).Return(errors.New("error"))

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", 1)
		c.Set("role", "moderator")
		c.Next()
	})
	router.PUT("/:id", suite.handler.UpdateProduct)

	requestBody, _ := json.Marshal(req)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("PUT", "/1", bytes.NewBuffer(requestBody))

	router.ServeHTTP(w, r)

	suite.Equal(http.StatusInternalServerError, w.Code)
	suite.Equal(`"Failed to update product"`, w.Body.String())
}

// =====================================================================================================================

func (suite *ProductHandlerSuite) TestHandler_DeleteProductSuccess() {
	suite.service.On("DeleteProduct", 1).Return(nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", 1)
		c.AddParam("id", strconv.Itoa(1))
		c.Set("role", "moderator")
		c.Next()
	})
	router.DELETE("/:id", suite.handler.DeleteProduct)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("DELETE", "/1", nil)

	router.ServeHTTP(w, r)

	suite.Equal(http.StatusOK, w.Code)
	suite.Equal("1", w.Body.String())
}

func (suite *ProductHandlerSuite) TestHandler_DeleteProductServiceFailure() {
	suite.service.On("DeleteProduct", 1).Return(errors.New("error"))

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", 1)
		c.AddParam("id", strconv.Itoa(1))
		c.Set("role", "moderator")
		c.Next()
	})
	router.DELETE("/:id", suite.handler.DeleteProduct)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("DELETE", "/1", nil)

	router.ServeHTTP(w, r)

	suite.Equal(http.StatusInternalServerError, w.Code)
	suite.Equal(`"Failed to delete product"`, w.Body.String())
}
