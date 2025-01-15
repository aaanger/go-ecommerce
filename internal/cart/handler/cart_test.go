package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/aaanger/ecommerce/internal/cart/model"
	"github.com/aaanger/ecommerce/internal/cart/service/mocks"
	"github.com/aaanger/ecommerce/pkg/lib"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type CartHandlerSuite struct {
	suite.Suite
	service *mocks.ICartService
	handler *CartHandler
}

func (suite *CartHandlerSuite) SetupTest() {
	suite.service = mocks.NewICartService(suite.T())
	suite.handler = NewCartHandler(suite.service)
}

func TestOrderHandlerSuite(t *testing.T) {
	suite.Run(t, new(CartHandlerSuite))
}

// =====================================================================================================================

func (suite *CartHandlerSuite) TestHandler_GetCartSuccess() {
	suite.service.On("GetCartByUserID", 1).Return(
		&model.Cart{
			ID:     1,
			UserID: 1,
		}, nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", 1)
		c.Next()
	})
	router.GET("/", suite.handler.GetCart)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	router.ServeHTTP(w, r)

	var response interface{}
	var cartRes model.Cart
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	lib.Copy(&cartRes, &response)

	suite.Equal(http.StatusOK, w.Code)
	suite.Equal(1, cartRes.ID)
	suite.Equal(1, cartRes.UserID)
}

func (suite *CartHandlerSuite) TestHandler_GetCartUnauthorized() {
	router := gin.New()
	router.GET("/", suite.handler.GetCart)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	router.ServeHTTP(w, r)

	suite.Equal(http.StatusUnauthorized, w.Code)
	suite.Equal(`"user id not found"`, w.Body.String())
}

func (suite *CartHandlerSuite) TestHandler_GetCartServiceFailure() {
	suite.service.On("GetCartByUserID", 1).Return(nil, errors.New("error"))

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", 1)
		c.Next()
	})
	router.GET("/", suite.handler.GetCart)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	router.ServeHTTP(w, r)

	suite.Equal(http.StatusInternalServerError, w.Code)
	suite.Equal(`"get cart error"`, w.Body.String())
}

func (suite *CartHandlerSuite) TestHandler_AddProductSuccess() {
	req := &model.AddProductReq{
		ProductID: 1,
		Quantity:  1,
	}

	res := &model.Cart{
		ID:     1,
		UserID: 1,
		Lines: []model.CartLine{
			{
				ProductID: 1,
				Quantity:  1,
			},
		},
		TotalPrice: 123,
	}

	suite.service.On("AddProduct", 1, req.ProductID, req.Quantity).Return(res, nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", 1)
		c.Next()
	})
	router.POST("/add", suite.handler.AddProduct)

	requestBody, _ := json.Marshal(req)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/add", bytes.NewBuffer(requestBody))

	router.ServeHTTP(w, r)

	var response interface{}
	var cartRes model.Cart
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	lib.Copy(&cartRes, &response)

	suite.Equal(http.StatusOK, w.Code)
	suite.Equal(float64(123), cartRes.TotalPrice)
	suite.Equal(1, len(cartRes.Lines))
	suite.Equal(1, cartRes.ID)
}

func (suite *CartHandlerSuite) TestHandler_AddProductUnauthorized() {
	req := &model.AddProductReq{
		ProductID: 1,
		Quantity:  1,
	}

	router := gin.New()
	router.POST("/add", suite.handler.AddProduct)

	requestBody, _ := json.Marshal(req)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/add", bytes.NewBuffer(requestBody))

	router.ServeHTTP(w, r)

	suite.Equal(http.StatusUnauthorized, w.Code)
	suite.Equal(`"user id not found"`, w.Body.String())
}

func (suite *CartHandlerSuite) TestHandler_AddProductEmptyFields() {
	req := &model.AddProductReq{
		ProductID: 1,
	}

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", 1)
	})
	router.POST("/add", suite.handler.AddProduct)

	requestBody, _ := json.Marshal(req)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/add", bytes.NewBuffer(requestBody))

	router.ServeHTTP(w, r)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Equal(`"Invalid input parameters"`, w.Body.String())
}

func (suite *CartHandlerSuite) TestHandler_AddProductServiceFailure() {
	req := &model.AddProductReq{
		ProductID: 1,
		Quantity:  1,
	}

	suite.service.On("AddProduct", 1, req.ProductID, req.Quantity).Return(nil, errors.New("error"))

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", 1)
	})
	router.POST("/add", suite.handler.AddProduct)

	requestBody, _ := json.Marshal(req)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/add", bytes.NewBuffer(requestBody))

	router.ServeHTTP(w, r)

	suite.Equal(http.StatusInternalServerError, w.Code)
	suite.Equal(`"Add product to cart error"`, w.Body.String())
}

// =====================================================================================================================

func (suite *CartHandlerSuite) TestHandler_DeleteProductSuccess() {
	req := &model.DeleteProductReq{
		ProductID: 1,
	}

	res := &model.Cart{
		ID:     1,
		UserID: 1,
		Lines: []model.CartLine{
			{
				ProductID: 2,
				Quantity:  1,
			},
		},
	}

	suite.service.On("DeleteProduct", 1, req.ProductID).Return(res, nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", 1)
	})
	router.DELETE("/", suite.handler.DeleteProduct)

	requestBody, _ := json.Marshal(req)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("DELETE", "/", bytes.NewBuffer(requestBody))

	router.ServeHTTP(w, r)

	var response interface{}
	var cartRes model.Cart
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	lib.Copy(&cartRes, &response)

	suite.Equal(http.StatusOK, w.Code)
	suite.Equal([]model.CartLine{{ProductID: 2, Quantity: 1}}, cartRes.Lines)
	suite.Equal(1, cartRes.ID)
}

func (suite *CartHandlerSuite) TestHandler_DeleteProductUnauthorized() {
	req := &model.DeleteProductReq{
		ProductID: 1,
	}

	router := gin.New()
	router.DELETE("/", suite.handler.DeleteProduct)

	requestBody, _ := json.Marshal(req)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("DELETE", "/", bytes.NewBuffer(requestBody))

	router.ServeHTTP(w, r)

	suite.Equal(http.StatusUnauthorized, w.Code)
	suite.Equal(`"user id not found"`, w.Body.String())
}

func (suite *CartHandlerSuite) TestHandler_DeleteProductEmptyField() {
	req := &model.DeleteProductReq{}

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", 1)
	})
	router.DELETE("/", suite.handler.DeleteProduct)

	requestBody, _ := json.Marshal(req)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("DELETE", "/", bytes.NewBuffer(requestBody))

	router.ServeHTTP(w, r)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Equal(`"invalid input parameters"`, w.Body.String())
}

func (suite *CartHandlerSuite) TestHandler_DeleteProductServiceFailure() {
	req := &model.DeleteProductReq{
		ProductID: 1,
	}

	suite.service.On("DeleteProduct", 1, req.ProductID).Return(nil, errors.New("error"))

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", 1)
	})
	router.DELETE("/", suite.handler.DeleteProduct)

	requestBody, _ := json.Marshal(req)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("DELETE", "/", bytes.NewBuffer(requestBody))

	router.ServeHTTP(w, r)

	suite.Equal(http.StatusInternalServerError, w.Code)
	suite.Equal(`"Failed to delete product from the cart"`, w.Body.String())
}
