package repository

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aaanger/ecommerce/internal/product/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ProductRepositorySuite struct {
	suite.Suite
	db   *sql.DB
	mock sqlmock.Sqlmock
	repo *ProductRepository
}

func (suite *ProductRepositorySuite) SetupTest() {
	var err error
	suite.db, suite.mock, err = sqlmock.New()
	assert.NoError(suite.T(), err)
	suite.repo = NewProductRepository(suite.db)
}

func TestProductRepositorySuite(t *testing.T) {
	suite.Run(t, new(ProductRepositorySuite))
}

// ====================================================================================================================

func (suite *ProductRepositorySuite) TestRepository_CreateProductSuccess() {
	req := &model.ProductReq{
		Name:        "test",
		Description: "test",
		Price:       5,
		Amount:      5,
		InStock:     true,
	}

	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
	suite.mock.ExpectQuery("INSERT INTO products").WithArgs(req.Name, req.Description, req.Amount, req.Price, req.InStock).WillReturnRows(rows)

	product, err := suite.repo.CreateProduct(req)
	suite.NotNil(product)
	suite.Equal(1, product.ID)
	suite.Equal("test", product.Name)
	suite.Nil(err)
}

func (suite *ProductRepositorySuite) TestRepository_CreateProductFailure() {
	req := &model.ProductReq{}

	rows := sqlmock.NewRows([]string{"id"})
	suite.mock.ExpectQuery("INSERT INTO products").WithArgs().WillReturnRows(rows)

	product, err := suite.repo.CreateProduct(req)
	suite.Nil(product)
	suite.NotNil(err)
}

// ====================================================================================================================

func (suite *ProductRepositorySuite) TestRepository_GetAllProductsSuccess() {
	rows := sqlmock.NewRows([]string{"id", "name", "description", "price", "amount", "in_stock"}).AddRow(1, "test", "test", float64(5), 5, true).AddRow(2, "test2", "test2", float64(10), 10, true)
	suite.mock.ExpectQuery("SELECT id, name, description, price, amount, in_stock FROM products").WillReturnRows(rows)

	products, err := suite.repo.GetAllProducts()

	expected := []model.Product{
		{
			ID:          1,
			Name:        "test",
			Description: "test",
			Price:       5,
			Amount:      5,
			InStock:     true,
		},
		{
			ID:          2,
			Name:        "test2",
			Description: "test2",
			Price:       10,
			Amount:      10,
			InStock:     true,
		},
	}

	suite.NotNil(products)
	suite.Nil(err)
	suite.Equal(expected, products)
}

func (suite *ProductRepositorySuite) TestRepository_GetAllProductFailure() {
	suite.mock.ExpectQuery("SELECT id, name, description, price, amount, in_stock FROM products")

	products, err := suite.repo.GetAllProducts()

	suite.Nil(products)
	suite.NotNil(err)
}

// ====================================================================================================================

func (suite *ProductRepositorySuite) TestRepository_GetProductByIDSuccess() {
	rows := sqlmock.NewRows([]string{"name", "description", "price", "amount", "in_stock"}).AddRow("test", "test", float64(5), 5, true)

	suite.mock.ExpectQuery("SELECT name, description, price, amount, in_stock FROM products").WithArgs(1).WillReturnRows(rows)

	product, err := suite.repo.GetProductByID(1)

	expected := &model.Product{
		ID:          1,
		Name:        "test",
		Description: "test",
		Price:       5,
		Amount:      5,
		InStock:     true,
	}

	suite.NotNil(product)
	suite.Equal(expected, product)
	suite.Nil(err)
}

func (suite *ProductRepositorySuite) TestRepository_GetProductByIDFailure() {
	rows := sqlmock.NewRows([]string{"name", "description", "price", "amount", "in_stock"})

	suite.mock.ExpectQuery("SELECT name, description, price, amount, in_stock FROM products").WithArgs(1).WillReturnRows(rows)

	product, err := suite.repo.GetProductByID(1)

	suite.Nil(product)
	suite.NotNil(err)
}

// ====================================================================================================================

func strPtr(s string) *string {
	return &s
}

func (suite *ProductRepositorySuite) TestRepository_UpdateProductSuccess() {
	suite.mock.ExpectExec("UPDATE products").WithArgs("test", 1).WillReturnResult(sqlmock.NewResult(1, 1))

	req := model.UpdateProduct{
		Name: strPtr("test"),
	}
	err := suite.repo.UpdateProduct(1, req)

	suite.Nil(err)
}

func (suite *ProductRepositorySuite) TestRepository_UpdateProductFailure() {
	suite.mock.ExpectExec("UPDATE products SET name=$1 WHERE id=$2").WithArgs()

	req := model.UpdateProduct{}
	err := suite.repo.UpdateProduct(1, req)

	suite.NotNil(err)
}

// ====================================================================================================================

func (suite *ProductRepositorySuite) TestRepository_DeleteProductSuccess() {
	suite.mock.ExpectExec("DELETE FROM products").WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))

	err := suite.repo.DeleteProduct(1)

	suite.Nil(err)
}

func (suite *ProductRepositorySuite) TestRepository_DeleteProductFailure() {
	suite.mock.ExpectExec("DELETE FROM products").WithArgs()

	err := suite.repo.DeleteProduct(1)

	suite.NotNil(err)
}
