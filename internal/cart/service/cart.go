package service

import (
	"errors"
	"github.com/aaanger/ecommerce/internal/cart/model"
	"github.com/aaanger/ecommerce/internal/cart/repository"
	productRepository "github.com/aaanger/ecommerce/internal/product/repository"
)

//go:generate mockery --name=ICartService

type ICartService interface {
	GetCartByUserID(userID int) (*model.Cart, error)
	AddProduct(userID, productID, quantity int) (*model.Cart, error)
	DeleteProduct(userID, productID int) (*model.Cart, error)
}

type CartService struct {
	repo        repository.ICartRepository
	productRepo productRepository.IProductRepository
}

func NewCartService(repo repository.ICartRepository, productRepo productRepository.IProductRepository) *CartService {
	return &CartService{
		repo:        repo,
		productRepo: productRepo,
	}
}

func (s *CartService) GetCartByUserID(userID int) (*model.Cart, error) {
	return s.repo.GetCartByUserID(userID)
}

func (s *CartService) AddProduct(userID, productID, quantity int) (*model.Cart, error) {
	product, err := s.productRepo.GetProductByID(productID)
	if err != nil {
		return nil, err
	}
	if product.InStock == false {
		return nil, errors.New("product is not in stock")
	}

	cart, err := s.repo.GetCartByUserID(userID)
	if err != nil {
		cartID, err := s.repo.CreateCart(userID)
		if err != nil {
			return nil, err
		}
		cart = &model.Cart{
			ID:     cartID,
			UserID: userID,
		}
	}

	err = s.repo.AddProduct(cart.ID, productID, quantity)
	if err != nil {
		return nil, err
	}

	cart.Lines = append(cart.Lines, model.CartLine{
		ProductID: productID,
		Product:   product,
		Quantity:  quantity,
	})

	var totalPrice float64

	for _, line := range cart.Lines {
		totalPrice += line.Product.Price
	}

	cart.TotalPrice = totalPrice

	return cart, nil
}

func (s *CartService) DeleteProduct(userID, productID int) (*model.Cart, error) {
	cart, err := s.repo.GetCartByUserID(userID)
	if err != nil {
		return nil, err
	}

	err = s.repo.DeleteProduct(cart.ID, productID)
	if err != nil {
		return nil, err
	}

	for i, line := range cart.Lines {
		if line.ProductID == productID {
			cart.Lines = append(cart.Lines[:i], cart.Lines[i+1:]...)
			break
		}
	}

	return cart, nil
}
