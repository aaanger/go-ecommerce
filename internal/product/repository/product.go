package repository

import (
	"database/sql"
	"fmt"
	"github.com/aaanger/ecommerce/internal/product/model"
	"strings"
)

//go:generate mockery --name=IProductRepository

type IProductRepository interface {
	CreateProduct(req *model.ProductReq) (*model.Product, error)
	GetAllProducts() ([]model.Product, error)
	GetProductByID(id int) (*model.Product, error)
	UpdateProduct(id int, input model.UpdateProduct) error
	DeleteProduct(id int) error
}

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

func (r *ProductRepository) CreateProduct(req *model.ProductReq) (*model.Product, error) {
	product := model.Product{
		Name:        req.Name,
		Description: req.Description,
		Amount:      req.Amount,
		Price:       req.Price,
		InStock:     req.InStock,
	}

	row := r.db.QueryRow(`INSERT INTO products (name, description, amount, price, in_stock) VALUES($1, $2, $3, $4, $5) RETURNING id;`, req.Name, req.Description, req.Amount, req.Price, req.InStock)
	err := row.Scan(&product.ID)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *ProductRepository) GetAllProducts() ([]model.Product, error) {
	var products []model.Product

	rows, err := r.db.Query(`SELECT id, name, description, price, amount, in_stock FROM products;`)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var product model.Product

		err = rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Amount, &product.InStock)
		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	return products, rows.Err()
}

func (r *ProductRepository) GetProductByID(id int) (*model.Product, error) {
	product := model.Product{
		ID: id,
	}

	row := r.db.QueryRow(`SELECT name, description, price, amount, in_stock FROM products WHERE id=$1;`, id)
	err := row.Scan(&product.Name, &product.Description, &product.Amount, &product.Price, &product.InStock)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *ProductRepository) UpdateProduct(id int, input model.UpdateProduct) error {
	keys := make([]string, 0)
	values := make([]interface{}, 0)
	arg := 1

	if input.Name != nil {
		keys = append(keys, fmt.Sprintf("name=$%d", arg))
		values = append(values, *input.Name)
		arg++
	}
	if input.Description != nil {
		keys = append(keys, fmt.Sprintf("description=$%d", arg))
		values = append(values, *input.Description)
		arg++
	}
	if input.Price != nil {
		keys = append(keys, fmt.Sprintf("price=$%d", arg))
		values = append(values, *input.Price)
		arg++
	}
	if input.Amount != nil {
		keys = append(keys, fmt.Sprintf("amount=$%d", arg))
		values = append(values, *input.Amount)
		arg++
	}
	if input.InStock != nil {
		keys = append(keys, fmt.Sprintf("in_stock=$%d", arg))
		values = append(values, *input.InStock)
		arg++
	}

	joinKeys := strings.Join(keys, ", ")

	query := fmt.Sprintf(`UPDATE products SET %s WHERE id=$%d;`, joinKeys, arg)

	values = append(values, id)

	_, err := r.db.Exec(query, values...)
	if err != nil {
		return err
	}

	return nil
}

func (r *ProductRepository) DeleteProduct(id int) error {
	_, err := r.db.Exec(`DELETE FROM products WHERE id=$1;`, id)
	if err != nil {
		return err
	}
	return nil
}
