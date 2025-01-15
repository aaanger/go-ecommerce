package repository

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type CartRepositorySuite struct {
	suite.Suite
	db   *sql.DB
	mock sqlmock.Sqlmock
	repo *CartRepository
}

func (suite *CartRepositorySuite) SetupTest() {
	var err error
	suite.db, suite.mock, err = sqlmock.New()
	assert.NoError(suite.T(), err)
	suite.repo = NewCartRepository(suite.db)
}

func TestCartRepositorySuite(t *testing.T) {
	suite.Run(t, new(CartRepositorySuite))
}

// ====================================================================================================================

func (suite *CartRepositorySuite) TestRepository_CreateCartSuccess() {
	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
	suite.mock.ExpectQuery("INSERT INTO carts").WithArgs(1).WillReturnRows(rows)

	cartID, err := suite.repo.CreateCart(1)

	suite.NotNil(cartID)
	suite.Nil(err)
}

func (suite *CartRepositorySuite) TestRepository_CreateCartFailure() {
	rows := sqlmock.NewRows([]string{"id"})
	suite.mock.ExpectQuery("INSERT INTO carts").WithArgs(1).WillReturnRows(rows)

	cartID, err := suite.repo.CreateCart(1)

	suite.Equal(0, cartID)
	suite.NotNil(err)
}

// ====================================================================================================================

func (suite *CartRepositorySuite) TestRepository_GetCartByUserIDSuccess() {
	rows := sqlmock.NewRows([]string{"id", "user_id", "created_at", "updated_at"}).AddRow(1, 1, time.Now(), time.Now())
	suite.mock.ExpectQuery("SELECT id, user_id, created_at, updated_at FROM carts WHERE user_id=\\$1").WithArgs(1).WillReturnRows(rows)

	lineRows := sqlmock.NewRows([]string{"product_id", "quantity"}).AddRow(1, 1)
	suite.mock.ExpectQuery("SELECT product_id, quantity FROM cartline l INNER JOIN carts c ON c.id=l.cart_id WHERE c.id=\\$1").
		WithArgs(1).WillReturnRows(lineRows)

	cart, err := suite.repo.GetCartByUserID(1)

	suite.NotNil(cart)
	suite.Nil(err)
}

func (suite *CartRepositorySuite) TestRepository_GetCartByIDFailure() {
	rows := sqlmock.NewRows([]string{"id", "user_id", "created_at", "updated_at"})
	suite.mock.ExpectQuery("SELECT id, user_id, created_at, updated_at FROM carts WHERE user_id=\\$1").WithArgs(1).WillReturnRows(rows)

	cart, err := suite.repo.GetCartByUserID(1)

	suite.Nil(cart)
	suite.NotNil(err)
}

func (suite *CartRepositorySuite) TestRepository_GetCartByIDGetLinesFailure() {
	rows := sqlmock.NewRows([]string{"id", "user_id", "created_at", "updated_at"}).AddRow(1, 1, time.Now(), time.Now())
	suite.mock.ExpectQuery("SELECT id, user_id, created_at, updated_at FROM carts WHERE user_id=\\$1").WithArgs(1).WillReturnRows(rows)

	suite.mock.ExpectQuery("SELECT product_id, quantity FROM cartline l INNER JOIN carts c ON c.id=l.cart_id WHERE c.id=\\$1").
		WithArgs(1).WillReturnError(errors.New("error"))

	cart, err := suite.repo.GetCartByUserID(1)

	suite.Nil(cart)
	suite.NotNil(err)
}

// ====================================================================================================================

func (suite *CartRepositorySuite) TestRepository_AddProductSuccess() {
	suite.mock.ExpectExec("INSERT INTO cartline").WithArgs(1, 1, 1).WillReturnResult(sqlmock.NewResult(1, 1))

	err := suite.repo.AddProduct(1, 1, 1)

	suite.Nil(err)
}

func (suite *CartRepositorySuite) TestRepository_AddProductFailure() {
	suite.mock.ExpectExec("INSERT INTO cartline").WithArgs(1, 1, 1).WillReturnError(errors.New("error"))

	err := suite.repo.AddProduct(1, 1, 1)

	suite.NotNil(err)
}

// ====================================================================================================================

func (suite *CartRepositorySuite) TestRepository_DeleteProductSuccess() {
	suite.mock.ExpectExec("DELETE FROM cartline WHERE cart_id=\\$1 AND product_id=\\$2").WithArgs(1, 1).WillReturnResult(sqlmock.NewResult(1, 1))

	err := suite.repo.DeleteProduct(1, 1)

	suite.Nil(err)
}

func (suite *CartRepositorySuite) TestRepository_DeleteProductFailure() {
	suite.mock.ExpectExec("DELETE FROM cartline WHERE cart_id=\\$1 AND product_id=\\$2").WithArgs(1, 1).WillReturnError(errors.New("error"))

	err := suite.repo.DeleteProduct(1, 1)

	suite.NotNil(err)
}
