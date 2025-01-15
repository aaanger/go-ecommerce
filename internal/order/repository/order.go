package repository

import (
	"database/sql"
	"github.com/aaanger/ecommerce/internal/order/model"
	"time"
)

//go:generate mockery --name=IOrderRepository

type IOrderRepository interface {
	CreateOrder(userID int, lines []model.OrderLine) (*model.Order, error)
	GetOrderByID(userID, orderID int) (*model.Order, error)
	GetAllOrders(userID int) ([]model.Order, error)
	UpdateOrder(userID, orderID int, status string) error
}

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

func (r *OrderRepository) CreateOrder(userID int, lines []model.OrderLine) (*model.Order, error) {
	var totalPrice float64

	for _, line := range lines {
		totalPrice += line.Price
	}

	order := model.Order{
		UserID:     userID,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		Lines:      lines,
		Status:     model.StatusOrderCreated,
		TotalPrice: totalPrice,
	}

	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}

	row := tx.QueryRow(`INSERT INTO orders (user_id, created_at, updated_at, status, total_price) VALUES($1, $2, $3, $4, $5) RETURNING id;`,
		order.UserID, order.CreatedAt, order.UpdatedAt, order.Status, order.TotalPrice)

	err = row.Scan(&order.ID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	for _, line := range lines {
		_, err = tx.Exec(`INSERT INTO orderline (order_id, product_id, quantity, price) VALUES($1, $2, $3, $4);`,
			order.ID, line.ProductID, line.Quantity, line.Price)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	return &order, tx.Commit()
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
