package repository

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aaanger/ecommerce/internal/order/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type OrderRepositorySuite struct {
	suite.Suite
	db   *sql.DB
	mock sqlmock.Sqlmock
	repo *OrderRepository
}

func (suite *OrderRepositorySuite) SetupTest() {
	var err error
	suite.db, suite.mock, err = sqlmock.New()
	assert.NoError(suite.T(), err)
	suite.repo = NewOrderRepository(suite.db)
}

func TestOrderRepositorySuite(t *testing.T) {
	suite.Run(t, new(OrderRepositorySuite))
}

// ====================================================================================================================

func (suite *OrderRepositorySuite) TestRepository_CreateOrderSuccess() {
	reqLines := []model.OrderLine{
		{
			ProductID: 1,
			Quantity:  1,
			Price:     5,
		},
	}

	suite.mock.ExpectBegin()
	orderRows := sqlmock.NewRows([]string{"id"}).AddRow(1)
	suite.mock.ExpectQuery("INSERT INTO orders").
		WithArgs(1, time.Now(), time.Now(), model.StatusOrderCreated, float64(5)).WillReturnRows(orderRows)

	suite.mock.ExpectExec("INSERT INTO orderline").WithArgs(1, reqLines[0].ProductID, reqLines[0].Quantity, reqLines[0].Price).
		WillReturnResult(sqlmock.NewResult(1, 1))

	suite.mock.ExpectCommit()

	order, err := suite.repo.CreateOrder(1, reqLines)

	suite.NotNil(order)
	suite.Nil(err)
}

func (suite *OrderRepositorySuite) TestRepository_CreateOrderFailureOrder() {
	reqLines := []model.OrderLine{
		{
			ProductID: 1,
			Quantity:  1,
			Price:     5,
		},
	}

	suite.mock.ExpectBegin()
	suite.mock.ExpectQuery("INSERT INTO orders").WithArgs(1, time.Now(), time.Now(), model.StatusOrderCreated, float64(5)).
		WillReturnError(errors.New("error"))
	suite.mock.ExpectRollback()

	order, err := suite.repo.CreateOrder(1, reqLines)

	suite.Nil(order)
	suite.NotNil(err)
}

func (suite *OrderRepositorySuite) TestRepository_CreateOrderFailureOrderLines() {
	reqLines := []model.OrderLine{
		{
			ProductID: 1,
			Quantity:  1,
			Price:     5,
		},
	}

	suite.mock.ExpectBegin()
	orderRows := sqlmock.NewRows([]string{"id"}).AddRow(1)
	suite.mock.ExpectQuery("INSERT INTO orders").WithArgs(1, reqLines[0].ProductID, reqLines[0].Quantity, reqLines[0].Price).
		WillReturnRows(orderRows)

	suite.mock.ExpectExec("INSERT INTO orderline").WithArgs(1, reqLines[0].ProductID, reqLines[0].Quantity, reqLines[0].Price).
		WillReturnError(errors.New("error"))
	suite.mock.ExpectRollback()

	order, err := suite.repo.CreateOrder(1, reqLines)

	suite.Nil(order)
	suite.NotNil(err)
}

// ====================================================================================================================

func (suite *OrderRepositorySuite) TestRepository_GetOrderByIDSuccess() {
	orderRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "status", "total_price"}).AddRow(1, time.Now(), time.Now(), model.StatusOrderCreated, float64(5))
	suite.mock.ExpectQuery("SELECT id, created_at, updated_at, status, total_price FROM orders").WithArgs(1, 1).WillReturnRows(orderRows)

	lineRows := sqlmock.NewRows([]string{"product_id", "quantity", "price"}).AddRow(1, 1, float64(5))
	suite.mock.ExpectQuery("SELECT product_id, quantity, price FROM orderline ol INNER JOIN orders o ON ol.order_id=o.id WHERE o.id=\\$1").WithArgs(1).WillReturnRows(lineRows)

	order, err := suite.repo.GetOrderByID(1, 1)

	suite.NotNil(order)
	suite.Nil(err)
}

func (suite *OrderRepositorySuite) TestRepository_GetOrderByIDFailure() {
	orderRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "status", "total_price"})
	suite.mock.ExpectQuery("SELECT id, created_at, updated_at, status, total_price FROM orders").WithArgs(1, 1).WillReturnRows(orderRows)

	lineRows := sqlmock.NewRows([]string{"product_id", "quantity", "price"}).AddRow(1, 1, float64(5))
	suite.mock.ExpectQuery("SELECT product_id, quantity, price FROM orderline ol INNER JOIN orders o ON ol.order_id=o.id WHERE o.id=\\$1").WithArgs(1).WillReturnRows(lineRows)

	order, err := suite.repo.GetOrderByID(1, 1)

	suite.Nil(order)
	suite.NotNil(err)
}

// ====================================================================================================================

func (suite *OrderRepositorySuite) TestRepository_GetAllOrdersSuccess() {
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "status", "total_price"}).AddRow(1, time.Now(), time.Now(), model.StatusOrderCreated, float64(5))
	suite.mock.ExpectQuery("SELECT id, created_at, updated_at, status, total_price FROM orders").WithArgs(1).WillReturnRows(rows)

	orders, err := suite.repo.GetAllOrders(1)

	suite.NotNil(orders)
	suite.Nil(err)
}

func (suite *OrderRepositorySuite) TestRepository_GetAllOrdersFailure() {

	suite.mock.ExpectQuery("SELECT id, created_at, updated_at, status, total_price FROM orders").WithArgs(1).WillReturnError(errors.New("error"))

	orders, err := suite.repo.GetAllOrders(1)

	suite.Nil(orders)
	suite.NotNil(err)
}

// ====================================================================================================================

func (suite *OrderRepositorySuite) TestRepository_UpdateOrderSuccess() {
	suite.mock.ExpectExec("UPDATE orders SET updated_at = current_timestamp, status=\\$1 WHERE id=\\$2 AND user_id=\\$3").WithArgs(model.StatusOrderDelivered, 1, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := suite.repo.UpdateOrder(1, 1, model.StatusOrderDelivered)

	suite.Nil(err)
}

func (suite *OrderRepositorySuite) TestRepository_UpdateOrderFailure() {
	suite.mock.ExpectExec("UPDATE orders SET updated_at = current_timestamp, status=\\$1 WHERE id=\\$2 AND user_id=\\$3").WithArgs(model.StatusOrderDelivered, 1, 1)

	err := suite.repo.UpdateOrder(1, 1, model.StatusOrderDelivered)

	suite.NotNil(err)
}
