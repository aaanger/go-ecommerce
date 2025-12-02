package service

import (
	"errors"
	"github.com/aaanger/ecommerce/internal/product/model"
	"github.com/aaanger/ecommerce/internal/product/repository/mocks"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ProductServiceSuite struct {
	suite.Suite
	repo    *mocks.IProductRepository
	service *ProductService
}

func (suite *ProductServiceSuite) SetupTest() {
	suite.repo = mocks.NewIProductRepository(suite.T())
	suite.service = NewProductService(suite.repo)
}

func TestProductServiceSuite(t *testing.T) {
	suite.Run(t, new(ProductServiceSuite))
}

func (suite *ProductServiceSuite) TestService_CreateProductSuccess() {
	req := &model.ProductReq{
		Name:        "test",
		Description: "test",
		Price:       5,
		Amount:      5,
		InStock:     true,
	}

	res := &model.Product{
		ID:          1,
		Name:        "test",
		Description: "test",
		Price:       5,
		Amount:      5,
		InStock:     true,
	}

	suite.repo.On("CreateProduct", req).Return(res, nil)

	product, err := suite.service.CreateProduct(req)

	suite.NotNil(product)
	suite.Equal("test", res.Name)
	suite.Equal("test", res.Description)
	suite.Nil(err)
}

func (suite *ProductServiceSuite) TestService_CreateProductRepoFailure() {
	req := &model.ProductReq{
		Name:        "test",
		Description: "test",
		Price:       5,
		Amount:      5,
		InStock:     true,
	}

	suite.repo.On("CreateProduct", req).Return(nil, errors.New("error"))

	product, err := suite.service.CreateProduct(req)
	suite.Nil(product)
	suite.NotNil(err)
}

// ====================================================================================================================

func (suite *ProductServiceSuite) TestService_GetAllProductsSuccess() {
	suite.repo.On("GetAllProducts").Return([]model.Product{
		{
			ID:   1,
			Name: "test",
		},
	}, nil)

	products, err := suite.service.GetAllProducts()
	suite.NotNil(products)
	suite.Equal(1, products[0].ID)
	suite.Equal("test", products[0].Name)
	suite.Nil(err)
}

func (suite *ProductServiceSuite) TestService_GetAllProductsRepoFailure() {
	suite.repo.On("GetAllProducts").Return(nil, errors.New("error"))

	products, err := suite.service.GetAllProducts()
	suite.Nil(products)
	suite.NotNil(err)
}

// ====================================================================================================================

func (suite *ProductServiceSuite) TestService_GetProductByIDSuccess() {
	suite.repo.On("GetProductByID", 1).Return(&model.Product{
		ID:   1,
		Name: "test",
	}, nil)

	product, err := suite.service.GetProductByID(1)
	suite.NotNil(product)
	suite.Equal(1, product.ID)
	suite.Equal("test", product.Name)
	suite.Nil(err)
}

func (suite *ProductServiceSuite) TestService_GetProductByIDRepoFailure() {
	suite.repo.On("GetProductByID", 1).Return(nil, errors.New("error"))

	product, err := suite.service.GetProductByID(1)
	suite.Nil(product)
	suite.NotNil(err)
}

// ====================================================================================================================

func stringPtr(s string) *string {
	return &s
}

func (suite *ProductServiceSuite) TestService_UpdateProductSuccess() {
	req := model.UpdateProduct{
		Name: stringPtr("test"),
	}

	suite.repo.On("UpdateProduct", 1, req).Return(nil)

	err := suite.service.UpdateProduct(1, req)
	suite.Nil(err)
}

func (suite *ProductServiceSuite) TestService_UpdateProductRepoFailure() {
	req := model.UpdateProduct{
		Name: stringPtr("test"),
	}

	suite.repo.On("UpdateProduct", 1, req).Return(errors.New("error"))

	err := suite.service.UpdateProduct(1, req)
	suite.NotNil(err)
}

// ====================================================================================================================

func (suite *ProductServiceSuite) TestService_DeleteProductSuccess() {
	suite.repo.On("DeleteProduct", 1).Return(nil)

	err := suite.service.DeleteProduct(1)

	suite.Nil(err)
}

func (suite *ProductServiceSuite) TestService_DeleteProductRepoFailure() {
	suite.repo.On("DeleteProduct", 1).Return(errors.New("error"))

	err := suite.service.DeleteProduct(1)

	suite.NotNil(err)
}
