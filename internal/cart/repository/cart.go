package repository

import (
	"database/sql"
	"errors"
	"github.com/aaanger/ecommerce/internal/cart/model"
)

//go:generate mockery --name=ICartRepository

type ICartRepository interface {
	CreateCart(userID int) (int, error)
	GetCartByUserID(userID int) (*model.Cart, error)
	AddProduct(cartID, productID, quantity int) error
	DeleteProduct(cartID, productID int) error
}

type CartRepository struct {
	db *sql.DB
}

func NewCartRepository(db *sql.DB) *CartRepository {
	return &CartRepository{
		db: db,
	}
}

func (r *CartRepository) CreateCart(userID int) (int, error) {
	var id int

	row := r.db.QueryRow(`INSERT INTO carts (user_id, created_at) VALUES($1, current_timestamp) RETURNING id;`, userID)
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *CartRepository) GetCartByUserID(userID int) (*model.Cart, error) {
	var cart model.Cart

	row := r.db.QueryRow(`SELECT id, user_id, created_at, updated_at FROM carts WHERE user_id=$1;`, userID)
	err := row.Scan(&cart.ID, &cart.UserID, &cart.CreatedAt, &cart.UpdatedAt)
	if err != nil {
		return nil, err
	}

	var lines []model.CartLine

	rows, err := r.db.Query(`SELECT product_id, quantity FROM cartline l INNER JOIN carts c ON c.id=l.cart_id WHERE c.id=$1;`, cart.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &cart, nil
		}
		return nil, err
	}

	for rows.Next() {
		var line model.CartLine

		err = rows.Scan(&line.ProductID, &line.Quantity)
		if err != nil {
			return nil, err
		}

		lines = append(lines, line)
	}

	cart.Lines = lines

	return &cart, nil
}

func (r *CartRepository) AddProduct(cartID, productID, quantity int) error {
	_, err := r.db.Exec(`INSERT INTO cartline (cart_id, product_id, quantity) VALUES($1, $2, $3);`, cartID, productID, quantity)
	if err != nil {
		return err
	}

	return nil
}

func (r *CartRepository) DeleteProduct(cartID, productID int) error {
	_, err := r.db.Exec(`DELETE FROM cartline WHERE cart_id=$1 AND product_id=$2;`, cartID, productID)
	if err != nil {
		return err
	}

	return nil
}
