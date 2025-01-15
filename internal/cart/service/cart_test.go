package service

import (
	"errors"
	"github.com/aaanger/ecommerce/internal/cart/model"
	"github.com/aaanger/ecommerce/internal/cart/repository/mocks"
	productModel "github.com/aaanger/ecommerce/internal/product/model"
	productMocks "github.com/aaanger/ecommerce/internal/product/repository/mocks"
	"github.com/stretchr/testify/suite"
	"testing"
)

type CartServiceSuite struct {
	suite.Suite
	repo        *mocks.ICartRepository
	productRepo *productMocks.IProductRepository
	service     *CartService
}

func (suite *CartServiceSuite) SetupTest() {
	suite.repo = mocks.NewICartRepository(suite.T())
	suite.productRepo = productMocks.NewIProductRepository(suite.T())
	suite.service = NewCartService(suite.repo, suite.productRepo)
}

func TestCartServiceSuite(t *testing.T) {
	suite.Run(t, new(CartServiceSuite))
}

// ====================================================================================================================

func (suite *CartServiceSuite) TestService_GetCartByIDSuccess() {
	suite.repo.On("GetCartByUserID", 1).Return(&model.Cart{
		ID:     1,
		UserID: 1,
	}, nil)

	cart, err := suite.service.GetCartByUserID(1)

	suite.NotNil(cart)
	suite.Nil(err)
}

func (suite *CartServiceSuite) TestService_GetCartByIDFailure() {
	suite.repo.On("GetCartByUserID", 1).Return(nil, errors.New("error"))

	cart, err := suite.service.GetCartByUserID(1)

	suite.Nil(cart)
	suite.NotNil(err)
}

// ====================================================================================================================

func (suite *CartServiceSuite) TestService_AddProductSuccess() {
	suite.productRepo.On("GetProductByID", 1).Return(&productModel.Product{
		ID:          1,
		Name:        "test",
		Description: "test",
		Price:       5,
		Amount:      5,
		InStock:     true,
	}, nil)

	suite.repo.On("GetCartByUserID", 1).Return(&model.Cart{
		ID:     1,
		UserID: 1,
	}, nil)

	suite.repo.On("AddProduct", 1, 1, 1).Return(nil)

	cart, err := suite.service.AddProduct(1, 1, 1)

	suite.NotNil(cart)
	suite.Nil(err)
}

func (suite *CartServiceSuite) TestService_AddProductGetCartFailure() {
	suite.productRepo.On("GetProductByID", 1).Return(&productModel.Product{
		ID:          1,
		Name:        "test",
		Description: "test",
		Price:       5,
		Amount:      5,
		InStock:     true,
	}, nil)

	suite.repo.On("GetCartByUserID", 1).Return(nil, errors.New("error"))

	suite.repo.On("CreateCart", 1).Return(1, nil)

	suite.repo.On("AddProduct", 1, 1, 1).Return(nil)

	cart, err := suite.service.AddProduct(1, 1, 1)

	suite.NotNil(cart)
	suite.Nil(err)
}

func (suite *CartServiceSuite) TestService_AddProductCreateCartFailure() {
	suite.productRepo.On("GetProductByID", 1).Return(&productModel.Product{
		ID:          1,
		Name:        "test",
		Description: "test",
		Price:       5,
		Amount:      5,
		InStock:     true,
	}, nil)

	suite.repo.On("GetCartByUserID", 1).Return(nil, errors.New("error"))

	suite.repo.On("CreateCart", 1).Return(0, errors.New("error"))

	cart, err := suite.service.AddProduct(1, 1, 1)

	suite.Nil(cart)
	suite.NotNil(err)
}

func (suite *CartServiceSuite) TestService_AddProductGetProductFailure() {
	suite.productRepo.On("GetProductByID", 1).Return(nil, errors.New("error"))

	cart, err := suite.service.AddProduct(1, 1, 1)

	suite.Nil(cart)
	suite.NotNil(err)
}

func (suite *CartServiceSuite) TestService_AddProductFailure() {
	suite.productRepo.On("GetProductByID", 1).Return(&productModel.Product{
		ID:          1,
		Name:        "test",
		Description: "test",
		Price:       5,
		Amount:      5,
		InStock:     true,
	}, nil)

	suite.repo.On("GetCartByUserID", 1).Return(&model.Cart{
		ID:     1,
		UserID: 1,
	}, nil)

	suite.repo.On("AddProduct", 1, 1, 1).Return(errors.New("error"))

	cart, err := suite.service.AddProduct(1, 1, 1)

	suite.Nil(cart)
	suite.NotNil(err)
}

// ====================================================================================================================

func (suite *CartServiceSuite) TestService_DeleteProductSuccess() {
	suite.repo.On("GetCartByUserID", 1).Return(&model.Cart{
		ID:     1,
		UserID: 1,
	}, nil)

	suite.repo.On("DeleteProduct", 1, 1).Return(nil)

	cart, err := suite.service.DeleteProduct(1, 1)

	suite.NotNil(cart)
	suite.Nil(err)
}

func (suite *CartServiceSuite) TestService_DeleteProductGetCartFailure() {
	suite.repo.On("GetCartByUserID", 1).Return(nil, errors.New("error"))

	cart, err := suite.service.DeleteProduct(1, 1)

	suite.Nil(cart)
	suite.NotNil(err)
}

func (suite *CartServiceSuite) TestService_DeleteProductFailure() {
	suite.repo.On("GetCartByUserID", 1).Return(&model.Cart{
		ID:     1,
		UserID: 1,
	}, nil)

	suite.repo.On("DeleteProduct", 1, 1).Return(errors.New("error"))

	cart, err := suite.service.DeleteProduct(1, 1)

	suite.Nil(cart)
	suite.NotNil(err)
}
