package repository

import (
	"context"
	"database/sql"
	"github.com/aaanger/ecommerce/internal/order/model"
	"go.uber.org/zap"
	"time"
)

//go:generate mockery --name=IOrderRepository

type IOrderRepository interface {
	GetOrderByID(userID, orderID int) (*model.Order, error)
	GetAllOrders(userID int) ([]model.Order, error)
	UpdateOrder(userID, orderID int, status string) error

	BeginTx(ctx context.Context) (*OrderTxRepository, error)
}

type IOrderTxRepository interface {
	CreateOrder(userID int, userEmail string, lines []model.OrderLine) (*model.Order, error)
	Commit() error
	Rollback() error
}

type OrderRepository struct {
	db  *sql.DB
	log *zap.Logger
}

type OrderTxRepository struct {
	tx  *sql.Tx
	log *zap.Logger
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

func (r *OrderRepository) BeginTx(ctx context.Context) (*OrderTxRepository, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &OrderTxRepository{tx: tx}, nil
}

func (r *OrderTxRepository) CreateOrder(userID int, userEmail string, lines []model.OrderLine) (*model.Order, error) {
	log := r.log.With(
		zap.String("service", "order"),
		zap.String("layer", "repository"),
		zap.String("method", "CreateOrder"),
		zap.Int("userID", userID))

	var totalPrice float64

	for _, line := range lines {
		totalPrice += line.Price
	}

	order := model.Order{
		UserID:     userID,
		UserEmail:  userEmail,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		Lines:      lines,
		Status:     model.StatusOrderCreated,
		TotalPrice: totalPrice,
	}

	log.Debug("Executing INSERT query on orders")
	row := r.tx.QueryRow(`INSERT INTO orders (user_id, user_email, created_at, updated_at, status, total_price) VALUES($1, $2, $3, $4, $5, $6) RETURNING id;`,
		order.UserID, order.UserEmail, order.CreatedAt, order.UpdatedAt, order.Status, order.TotalPrice)

	err := row.Scan(&order.ID)
	if err != nil {
		log.Error("Failed to create order", zap.Error(err))
		return nil, err
	}

	log.Debug("Executing INSERT query on orderline")
	for _, line := range lines {
		_, err = r.tx.Exec(`INSERT INTO orderline (order_id, product_id, quantity, price) VALUES($1, $2, $3, $4);`,
			order.ID, line.ProductID, line.Quantity, line.Price)
		if err != nil {
			log.Error("Failed to create orderline", zap.Error(err))
			return nil, err
		}
	}

	log.Info("Order successfully created", zap.Int("orderID", order.ID))
	return &order, nil
}

func (r *OrderTxRepository) Commit() error {
	return r.tx.Commit()
}

func (r *OrderTxRepository) Rollback() error {
	return r.tx.Rollback()
}

func (r *OrderRepository) GetOrderByID(userID, orderID int) (*model.Order, error) {
	var order model.Order

	row := r.db.QueryRow(`SELECT id, created_at, updated_at, status, total_price FROM orders WHERE id=$1 AND user_id=$2;`, orderID, userID)
	err := row.Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt, &order.Status, &order.TotalPrice)
	if err != nil {
		return nil, err
	}

	var lines []model.OrderLine

	rows, err := r.db.Query(`SELECT product_id, quantity, price FROM orderline ol INNER JOIN orders o ON ol.order_id=o.id WHERE o.id=$1;`, orderID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var line model.OrderLine

		err = rows.Scan(&line.ProductID, &line.Quantity, &line.Price)
		if err != nil {
			return nil, err
		}

		lines = append(lines, line)
	}

	order.Lines = lines

	return &order, nil
}

func (r *OrderRepository) GetAllOrders(userID int) ([]model.Order, error) {
	var orders []model.Order

	rows, err := r.db.Query(`SELECT id, created_at, updated_at, status, total_price FROM orders WHERE user_id=$1;`, userID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var order model.Order

		err = rows.Scan(&order.ID, &order.UpdatedAt, &order.UpdatedAt, &order.Status, &order.TotalPrice)
		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	return orders, nil
}

func (r *OrderRepository) UpdateOrder(userID, orderID int, status string) error {
	_, err := r.db.Exec(`UPDATE orders SET updated_at = current_timestamp, status=$1 WHERE id=$2 AND user_id=$3;`, status, orderID, userID)
	if err != nil {
		return err
	}
	return nil
}
