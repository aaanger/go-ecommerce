package service

import (
	"errors"
	"github.com/aaanger/ecommerce/internal/order/model"
	"github.com/aaanger/ecommerce/internal/order/repository/mocks"
	productModel "github.com/aaanger/ecommerce/internal/product/model"
	productMocks "github.com/aaanger/ecommerce/internal/product/repository/mocks"
	"github.com/stretchr/testify/suite"
	"testing"
)

type OrderServiceSuite struct {
	suite.Suite
	repo        *mocks.IOrderRepository
	productRepo *productMocks.IProductRepository
	service     *OrderService
}

func (suite *OrderServiceSuite) SetupTest() {
	suite.repo = mocks.NewIOrderRepository(suite.T())
	suite.productRepo = productMocks.NewIProductRepository(suite.T())
	suite.service = NewOrderService(suite.repo, suite.productRepo)
}

func TestOrderServiceSuite(t *testing.T) {
	suite.Run(t, new(OrderServiceSuite))
}

// ====================================================================================================================

func (suite *OrderServiceSuite) TestService_CreateOrderSuccess() {
	req := &model.CreateOrderReq{
		Lines: []model.OrderLineReq{
			{
				ProductID: 1,
				Quantity:  1,
			},
		},
	}

	suite.productRepo.On("GetProductByID", 1).Return(&productModel.Product{
		ID:          1,
		Name:        "test",
		Description: "test",
		Price:       5,
	}, nil)

	suite.repo.On("CreateOrder", 1, []model.OrderLine{
		{
			ProductID: 1,
			Quantity:  1,
			Price:     5,
		},
	}).Return(&model.Order{
		ID:     1,
		UserID: 1,
		Lines: []model.OrderLine{
			{
				ProductID: 1,
				Quantity:  1,
			},
		},
	}, nil)

	order, err := suite.service.CreateOrder(1, req)

	suite.NotNil(order)
	suite.Nil(err)

}

func (suite *OrderServiceSuite) TestService_CreateOrderFailure() {
	req := &model.CreateOrderReq{
		Lines: []model.OrderLineReq{
			{
				ProductID: 1,
				Quantity:  1,
			},
		},
	}

	suite.productRepo.On("GetProductByID", 1).Return(&productModel.Product{
		ID:          1,
		Name:        "test",
		Description: "test",
		Price:       5,
	}, nil)

	suite.repo.On("CreateOrder", 1, []model.OrderLine{
		{
			ProductID: 1,
			Quantity:  1,
			Price:     5,
		}}).Return(nil, errors.New("error"))

	order, err := suite.service.CreateOrder(1, req)

	suite.Nil(order)
	suite.NotNil(err)
}

func (suite *OrderServiceSuite) TestService_CreateOrderGetProductFailure() {
	req := &model.CreateOrderReq{
		Lines: []model.OrderLineReq{
			{
				ProductID: 1,
				Quantity:  1,
			},
		},
	}

	suite.productRepo.On("GetProductByID", 1).Return(nil, errors.New("error"))

	order, err := suite.service.CreateOrder(1, req)

	suite.Nil(order)
	suite.NotNil(err)
}

// ====================================================================================================================

func (suite *OrderServiceSuite) TestService_GetOrderByIDSuccess() {
	suite.repo.On("GetOrderByID", 1, 1).Return(&model.Order{
		ID:     1,
		UserID: 1,
		Lines: []model.OrderLine{
			{
				ProductID: 1,
			},
		},
	}, nil)

	suite.productRepo.On("GetProductByID", 1).Return(&productModel.Product{
		ID:   1,
		Name: "test",
	}, nil)

	order, err := suite.service.GetOrderByID(1, 1)

	suite.NotNil(order)
	suite.Nil(err)
}

func (suite *OrderServiceSuite) TestService_GetOrderByIDFailure() {
	suite.repo.On("GetOrderByID", 1, 1).Return(nil, errors.New("error"))

	order, err := suite.service.GetOrderByID(1, 1)

	suite.Nil(order)
	suite.NotNil(err)
}

func (suite *OrderServiceSuite) TestService_GetOrderByIDGetProductFailure() {
	suite.repo.On("GetOrderByID", 1, 1).Return(&model.Order{
		ID:     1,
		UserID: 1,
		Lines: []model.OrderLine{
			{
				ProductID: 1,
			},
		},
	}, nil)

	suite.productRepo.On("GetProductByID", 1).Return(nil, errors.New("error"))

	order, err := suite.service.GetOrderByID(1, 1)

	suite.Nil(order)
	suite.NotNil(err)
}

// ====================================================================================================================

func (suite *OrderServiceSuite) TestService_GetAllOrdersSuccess() {
	suite.repo.On("GetAllOrders", 1).Return([]model.Order{
		{
			ID:     1,
			UserID: 1,
			Lines: []model.OrderLine{
				{
					ProductID: 1,
				},
			},
		},
	}, nil)

	orders, err := suite.service.GetAllOrders(1)

	suite.NotNil(orders)
	suite.Nil(err)
}

func (suite *OrderServiceSuite) TestService_GetAllOrdersFailure() {
	suite.repo.On("GetAllOrders", 1).Return(nil, errors.New("error"))

	orders, err := suite.service.GetAllOrders(1)

	suite.Nil(orders)
	suite.NotNil(err)
}

// ====================================================================================================================

func (suite *OrderServiceSuite) TestService_CancelOrderSuccess() {
	suite.repo.On("GetOrderByID", 1, 1).Return(&model.Order{
		ID:     1,
		UserID: 1,
		Lines: []model.OrderLine{
			{
				ProductID: 1,
			},
		},
	}, nil)

	suite.repo.On("UpdateOrder", 1, 1, model.StatusOrderCanceled).Return(nil)

	order, err := suite.service.CancelOrder(1, 1)

	suite.NotNil(order)
	suite.Nil(err)
}

func (suite *OrderServiceSuite) TestService_CancelOrderFailure() {
	suite.repo.On("GetOrderByID", 1, 1).Return(&model.Order{
		ID:     1,
		UserID: 1,
		Lines: []model.OrderLine{
			{
				ProductID: 1,
			},
		},
	}, nil)

	suite.repo.On("UpdateOrder", 1, 1, model.StatusOrderCanceled).Return(errors.New("error"))

	order, err := suite.service.CancelOrder(1, 1)

	suite.Nil(order)
	suite.NotNil(err)
}
