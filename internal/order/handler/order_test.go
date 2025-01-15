package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/aaanger/ecommerce/internal/order/model"
	"github.com/aaanger/ecommerce/internal/order/service/mocks"
	"github.com/aaanger/ecommerce/pkg/lib"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

type OrderHandlerSuite struct {
	suite.Suite
	service *mocks.IOrderService
	handler *OrderHandler
}

func (suite *OrderHandlerSuite) SetupTest() {
	suite.service = mocks.NewIOrderService(suite.T())
	suite.handler = NewOrderHandler(suite.service)
}

func TestOrderHandlerSuite(t *testing.T) {
	suite.Run(t, new(OrderHandlerSuite))
}

// =====================================================================================================================

func (suite *OrderHandlerSuite) TestHandler_CreateOrderOK() {
	req := &model.CreateOrderReq{
		Lines: []model.OrderLineReq{
			{
				ProductID: 1,
				Quantity:  1,
			},
			{
				ProductID: 2,
				Quantity:  2,
			},
		},
	}

	res := &model.Order{
		ID:     1,
		UserID: 1,
		Lines: []model.OrderLine{
			{
				ProductID: 1,
				Quantity:  1,
			},
			{
				ProductID: 2,
				Quantity:  2,
			},
		},
		Status:     model.StatusOrderCreated,
		TotalPrice: 123,
	}

	suite.service.On("CreateOrder", 1, req).Return(res, nil).Times(1)

	requestBody, _ := json.Marshal(req)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", 1)
		c.Next()
	})
	router.POST("/create", suite.handler.CreateOrder)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/create", bytes.NewBuffer(requestBody))
	router.ServeHTTP(w, r)

	var response interface{}
	var orderRes model.Order
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	lib.Copy(&orderRes, &response)

	suite.Equal(http.StatusOK, w.Code)
	suite.Equal(model.StatusOrderCreated, orderRes.Status)
	suite.Equal(float64(123), orderRes.TotalPrice)
	suite.Equal(2, len(orderRes.Lines))
}

func (suite *OrderHandlerSuite) TestHandler_CreateOrderEmptyFields() {
	req := &model.CreateOrderReq{
		Lines: []model.OrderLineReq{
			{
				Quantity: 1,
			},
		},
	}

	requestBody, _ := json.Marshal(req)

	router := gin.New()

	router.POST("/create", suite.handler.CreateOrder)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/create", bytes.NewBuffer(requestBody))
	router.ServeHTTP(w, r)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Equal(`"Failed to parse request body"`, w.Body.String())
}

func (suite *OrderHandlerSuite) TestHandler_CreateOrderUnauthorized() {
	req := &model.CreateOrderReq{
		Lines: []model.OrderLineReq{
			{
				ProductID: 1,
				Quantity:  1,
			},
			{
				ProductID: 2,
				Quantity:  2,
			},
		},
	}

	requestBody, _ := json.Marshal(req)

	router := gin.New()

	router.POST("/create", suite.handler.CreateOrder)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/create", bytes.NewBuffer(requestBody))
	router.ServeHTTP(w, r)

	suite.Equal(http.StatusUnauthorized, w.Code)
	suite.Equal(`"user id not found"`, w.Body.String())
}

func (suite *OrderHandlerSuite) TestHandler_CreateOrderServiceFailure() {
	req := &model.CreateOrderReq{
		Lines: []model.OrderLineReq{
			{
				ProductID: 1,
				Quantity:  1,
			},
			{
				ProductID: 2,
				Quantity:  2,
			},
		},
	}

	suite.service.On("CreateOrder", 1, req).Return(nil, errors.New("error"))

	requestBody, _ := json.Marshal(req)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", 1)
		c.Next()
	})
	router.POST("/create", suite.handler.CreateOrder)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/create", bytes.NewBuffer(requestBody))
	router.ServeHTTP(w, r)

	suite.Equal(http.StatusInternalServerError, w.Code)
	suite.Equal(`"Failed to create order"`, w.Body.String())
}

// =====================================================================================================================

func (suite *OrderHandlerSuite) TestHandler_GetOrderByIDSuccess() {
	suite.service.On("GetOrderByID", 1, 1).Return(&model.Order{
		ID:     1,
		UserID: 1,
		Lines: []model.OrderLine{
			{
				ProductID: 1,
				Quantity:  1,
			},
			{
				ProductID: 2,
				Quantity:  2,
			},
		},
		Status:     model.StatusOrderCreated,
		TotalPrice: 123,
	}, nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", 1)
		c.AddParam("id", strconv.Itoa(1))
		c.Next()
	})
	router.GET("/:id", suite.handler.GetOrderByID)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/1", nil)

	router.ServeHTTP(w, r)

	var response interface{}
	var orderRes model.Order
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	lib.Copy(&orderRes, &response)

	suite.Equal(http.StatusOK, w.Code)
	suite.Equal(model.StatusOrderCreated, orderRes.Status)
	suite.Equal(float64(123), orderRes.TotalPrice)
	suite.Equal(2, len(orderRes.Lines))
}

func (suite *OrderHandlerSuite) TestHandler_GetOrderByIDUnauthorized() {
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.AddParam("id", strconv.Itoa(1))
		c.Next()
	})
	router.GET("/:id", suite.handler.GetOrderByID)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/1", nil)

	router.ServeHTTP(w, r)

	suite.Equal(http.StatusUnauthorized, w.Code)
	suite.Equal(`"user id not found"`, w.Body.String())
}

func (suite *OrderHandlerSuite) TestHandler_GetOrderByIDServiceFailure() {
	suite.service.On("GetOrderByID", 1, 1).Return(nil, errors.New("error"))

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", 1)
		c.AddParam("id", strconv.Itoa(1))
		c.Next()
	})
	router.GET("/:id", suite.handler.GetOrderByID)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/1", nil)

	router.ServeHTTP(w, r)

	suite.Equal(http.StatusInternalServerError, w.Code)
	suite.Equal(`"Order not found"`, w.Body.String())
}

// =====================================================================================================================

func (suite *OrderHandlerSuite) TestHandler_GetAllOrdersSuccess() {
	suite.service.On("GetAllOrders", 1).Return([]model.Order{
		{
			ID:     1,
			UserID: 1,
			Lines: []model.OrderLine{
				{
					ProductID: 1,
					Quantity:  1,
				},
				{
					ProductID: 2,
					Quantity:  2,
				},
			},
			Status:     model.StatusOrderCreated,
			TotalPrice: 123,
		},
	}, nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", 1)
		c.Next()
	})
	router.GET("/all", suite.handler.GetAllOrders)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/all", nil)

	router.ServeHTTP(w, r)

	var response interface{}
	var orderRes []model.Order
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	lib.Copy(&orderRes, &response)

	suite.Equal(http.StatusOK, w.Code)
	suite.Equal(model.StatusOrderCreated, orderRes[0].Status)
	suite.Equal(float64(123), orderRes[0].TotalPrice)
}

func (suite *OrderHandlerSuite) TestHandler_GetAllOrdersUnauthorized() {
	router := gin.New()

	router.GET("/all", suite.handler.GetAllOrders)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/all", nil)

	router.ServeHTTP(w, r)

	suite.Equal(http.StatusUnauthorized, w.Code)
	suite.Equal(`"user id not found"`, w.Body.String())
}

func (suite *OrderHandlerSuite) TestHandler_GetAllOrdersServiceFailure() {
	suite.service.On("GetAllOrders", 1).Return(nil, errors.New("error"))

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", 1)
		c.Next()
	})

	router.GET("/all", suite.handler.GetAllOrders)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/all", nil)

	router.ServeHTTP(w, r)

	suite.Equal(http.StatusInternalServerError, w.Code)
	suite.Equal(`"Orders not found"`, w.Body.String())
}

// =====================================================================================================================

func (suite *OrderHandlerSuite) TestHandler_CancelOrderSuccess() {
	suite.service.On("CancelOrder", 1, 1).Return(&model.Order{
		ID:     1,
		UserID: 1,
		Lines: []model.OrderLine{
			{
				ProductID: 1,
				Quantity:  1,
			},
			{
				ProductID: 2,
				Quantity:  2,
			},
		},
		Status:     model.StatusOrderCreated,
		TotalPrice: 123,
	}, nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", 1)
		c.AddParam("id", strconv.Itoa(1))
		c.Next()
	})
	router.POST("/cancel/:id", suite.handler.CancelOrder)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/cancel/1", nil)

	router.ServeHTTP(w, r)

	var res interface{}
	var orderRes model.Order
	_ = json.Unmarshal(w.Body.Bytes(), &res)
	lib.Copy(&orderRes, &res)

	suite.Equal(http.StatusOK, w.Code)
	suite.Equal(model.StatusOrderCreated, orderRes.Status)
	suite.Equal(float64(123), orderRes.TotalPrice)
	suite.Equal(2, len(orderRes.Lines))
}

func (suite *OrderHandlerSuite) TestHandler_CancelOrderUnauthorized() {
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.AddParam("id", strconv.Itoa(1))
		c.Next()
	})
	router.POST("/cancel/:id", suite.handler.CancelOrder)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/cancel/1", nil)

	router.ServeHTTP(w, r)

	suite.Equal(http.StatusUnauthorized, w.Code)
	suite.Equal(`"user id not found"`, w.Body.String())
}

func (suite *OrderHandlerSuite) TestHandler_CancelOrderServiceFailure() {
	suite.service.On("CancelOrder", 1, 1).Return(nil, errors.New("error"))

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", 1)
		c.AddParam("id", strconv.Itoa(1))
		c.Next()
	})
	router.POST("/cancel/:id", suite.handler.CancelOrder)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/cancel/1", nil)

	router.ServeHTTP(w, r)

	suite.Equal(http.StatusInternalServerError, w.Code)
	suite.Equal(`"Failed to cancel order"`, w.Body.String())
}
